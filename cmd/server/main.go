package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"game-platform/internal/config"
	"game-platform/internal/handler"
	"game-platform/internal/middleware"
	"game-platform/internal/repository"
	"game-platform/internal/service"
	"game-platform/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration (reads config file + env var overrides)
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	dbPool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)

	// Initialize services with config-driven values
	portalAPI := service.NewPortalAPIWithTimeout(
		cfg.PortalAPI.URL,
		cfg.PortalAPI.APIKey,
		cfg.PortalAPITimeout(),
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userRepo, portalAPI)
	authHandler.SetJWTSecret(cfg.JWT.Secret)
	userHandler := handler.NewUserHandler(userRepo, nil)

	// Setup Gin
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimitWithRedis(rdb, cfg.RateLimit.MaxRequests, time.Duration(cfg.RateLimit.WindowSeconds)*time.Second))

	// Health check with timestamp
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// API v1
	api := r.Group("/api/v1")
	{
		// Auth
		auth := api.Group("/auth")
		auth.POST("/verify", authHandler.VerifyToken)

		// Users — protected by JWT authentication
		users := api.Group("/users")
		users.Use(middleware.Authenticate(cfg.JWT.Secret))
		{
			users.GET("/:sid/profile", userHandler.GetProfile)
		}
	}

	// WebSocket endpoint — protected by JWT authentication
	ws := r.Group("/ws")
	ws.Use(middleware.Authenticate(cfg.JWT.Secret))
	ws.GET("/game/:game_id", func(c *gin.Context) {
		wsHub.HandleWebSocket(c)
	})

	// Static files (frontend)
	r.Static("/static", "./web")
	r.StaticFile("/", "./web/index.html")

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on %s", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 5 seconds to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
