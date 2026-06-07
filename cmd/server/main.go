package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"game-platform/internal/config"
	"game-platform/internal/handler"
	"game-platform/internal/logger"
	"game-platform/internal/middleware"
	"game-platform/internal/repository"
	"game-platform/internal/service"
	"game-platform/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	// ── 1. Load config ──────────────────────────────────────────────────────
	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// ── 2. Init logger ──────────────────────────────────────────────────────
	logger.Init(cfg.Server.Mode)

	// ── 3. Connect PostgreSQL ──────────────────────────────────────────────
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		slog.Error("failed to parse database URL", "error", err)
		os.Exit(1)
	}
	if cfg.Database.PoolSize > 0 {
		poolCfg.MaxConns = int32(cfg.Database.PoolSize)
	}

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		slog.Error("failed to create connection pool", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to PostgreSQL")

	// ── 4. Connect Redis ────────────────────────────────────────────────────
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to ping Redis", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to Redis")

	// ── 5. Migrations (logged; Docker CMD handles actual migration) ──────────
	slog.Info("migrations: expected to be run via Docker CMD or `make migrate` before server start")

	// ── 6. Init repositories ────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	ratingRepo := repository.NewRatingRepository(db)
	tournamentRepo := repository.NewTournamentRepository(db)
	inviteRepo := repository.NewInviteRepository(db)

	// ── 7. Init services ───────────────────────────────────────────────────
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, userRepo)

	// ── 8. Init handlers ────────────────────────────────────────────────────
	embedHandler := handler.NewEmbedHandler(
		cfg.Embed.HandoffSecret,
		cfg.JWT.Secret,
		cfg.JWT.ExpiryHours,
	)
	embedHandler.SetUserRepo(userRepo)

	portalAPI := service.NewPortalAPIWithTimeout(
		cfg.PortalAPI.URL,
		cfg.PortalAPI.APIKey,
		cfg.PortalAPITimeout(),
	)

	authHandler := handler.NewAuthHandler(userRepo, portalAPI)
	authHandler.SetJWTSecret(cfg.JWT.Secret)

	userHandler := handler.NewUserHandler(userRepo, nil, matchRepo, tournamentRepo)

	gameHandler := handler.NewGameHandler(matchRepo, ratingSvc)

	inviteHandler := handler.NewInviteHandler(inviteRepo, matchRepo)
	inviteHandler.SetUserRepo(userRepo)

	ratingHandler := handler.NewRatingHandler(ratingRepo, ratingSvc, userRepo)

	tournamentHandler := handler.NewTournamentHandler(tournamentRepo, userRepo)

	// ── 9. Init WebSocket room manager ───────────────────────────────────────
	roomManager := websocket.NewRoomManagerWithDeps(cfg.JWT.Secret, matchRepo)
	roomManager.SetRatingService(ratingSvc)

	// ── 10. Build router ────────────────────────────────────────────────────
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Rate limiting (Redis-backed with in-memory fallback)
	rateLimitMW := middleware.RateLimitWithRedis(
		rdb,
		cfg.RateLimit.MaxRequests,
		time.Duration(cfg.RateLimit.WindowSeconds)*time.Second,
	)
	r.Use(rateLimitMW)

	// ── Public routes ────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
	r.POST("/api/v1/auth/embed", embedHandler.EmbedAuth)
	r.POST("/api/v1/auth/verify", authHandler.VerifyToken)

	// ── Protected /api/v1/* routes (JWT required) ────────────────────────────
	protected := r.Group("/api/v1")
	protected.Use(middleware.Authenticate(cfg.JWT.Secret))

	// Users
	protected.GET("/users/:sid/profile", userHandler.GetProfile)
	protected.GET("/users/:sid/stats", userHandler.GetStats)

	// Games
	protected.GET("/games/available", gameHandler.AvailableGames)
	protected.POST("/games/match/start", gameHandler.StartMatch)
	protected.POST("/games/match/complete", gameHandler.CompleteMatch)
	protected.GET("/games/match/history", gameHandler.GetMatchHistory)

	// Invites
	protected.POST("/games/invite", inviteHandler.CreateInvite)
	protected.POST("/games/invite/:id/accept", inviteHandler.AcceptInvite)
	protected.POST("/games/invite/:id/decline", inviteHandler.DeclineInvite)
	protected.GET("/games/invite/pending", inviteHandler.GetPendingInvites)

	// Ratings
	protected.GET("/ratings/:game_type", ratingHandler.GetRatings)
	protected.GET("/ratings/:game_type/leaderboard", gameHandler.GetLeaderboard)
	protected.GET("/ratings/:game_type/leaderboard/department", ratingHandler.GetLeaderboardByDepartment)
	protected.GET("/ratings/:game_type/me", gameHandler.GetUserRating)

	// Tournaments
	protected.GET("/tournaments", tournamentHandler.ListTournaments)
	protected.POST("/tournaments", tournamentHandler.CreateTournament)
	protected.GET("/tournaments/:id", tournamentHandler.GetTournament)
	protected.PUT("/tournaments/:id", tournamentHandler.UpdateTournament)
	protected.DELETE("/tournaments/:id", tournamentHandler.DeleteTournament)
	protected.POST("/tournaments/:id/join", tournamentHandler.JoinTournament)
	protected.POST("/tournaments/:id/leave", tournamentHandler.LeaveTournament)
	protected.GET("/tournaments/:id/players", tournamentHandler.GetTournamentPlayers)

	// ── WebSocket (NOT under /api/v1!) ────────────────────────────────────────
	r.GET("/ws/game/:match_id", roomManager.HandleWebSocket)

	// ── SPA static files ────────────────────────────────────────────────────
	spaDir := filepath.Join(".", "web", "dist")
	if info, err := os.Stat(spaDir); err == nil && info.IsDir() {
		// Serve Vite assets
		r.Static("/assets", filepath.Join(spaDir, "assets"))

		// SPA fallback for /games/* routes (Vite base is /games/)
		r.GET("/games/*filepath", func(c *gin.Context) {
			c.File(filepath.Join(spaDir, "index.html"))
		})

		// Root redirects to /games/
		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/games/")
		})
	} else {
		slog.Warn("web/dist directory not found, SPA static serving disabled")
	}

	// ── 11. Start HTTP server with graceful shutdown ────────────────────────
	addr := cfg.Addr()
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Channel to listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", addr, "mode", cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen failed", "error", err)
			os.Exit(1)
		}
	}()

	// Block until signal received
	sig := <-quit
	slog.Info("shutdown signal received", "signal", sig.String())

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced shutdown", "error", err)
	}

	slog.Info("shutdown complete")
}
