package main

import (
	"flag"
	"log"
	"os"
	"time"

	"game-platform/internal/config"
	"game-platform/internal/handler"
	"game-platform/internal/middleware"
	"game-platform/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// --- Auth routes (no JWT required) ---
	authHandler := handler.NewAuthHandler(nil, nil)
	authHandler.SetJWTSecret(cfg.JWT.Secret)

	// Embed auth: uses X-Erlink-Embed-Secret header, NOT JWT
	embedHandler := handler.NewEmbedHandler(
		cfg.Embed.HandoffSecret,
		cfg.JWT.Secret,
		cfg.JWT.ExpiryHours,
	)
	r.POST("/api/v1/auth/embed", embedHandler.EmbedAuth)

	// Verify token: requires JWT Bearer token
	r.POST("/api/v1/auth/verify", authHandler.VerifyToken)

	// Health endpoint (public, no JWT required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// --- Protected routes (JWT required) ---
	protected := r.Group("/api/v1")
	protected.Use(middleware.Authenticate(cfg.JWT.Secret))

	// --- WebSocket game rooms ---
	// RoomManager with JWT secret for query param auth and matchRepo for persistence.
	// WS clients can authenticate either via:
	//   1. Bearer token (through the middleware on the protected group), or
	//   2. Query parameter ?token=JWT (directly on the WS endpoint)
	roomManager := websocket.NewRoomManagerWithDeps(cfg.JWT.Secret, nil)
	protected.GET("/ws/game/:match_id", roomManager.HandleWebSocket)

	// --- Game invite routes ---
	inviteHandler := handler.NewInviteHandler(nil, nil)
	protected.POST("/api/v1/games/invite", inviteHandler.CreateInvite)
	protected.POST("/api/v1/games/invite/:id/accept", inviteHandler.AcceptInvite)
	protected.POST("/api/v1/games/invite/:id/decline", inviteHandler.DeclineInvite)
	protected.GET("/api/v1/games/invite/pending", inviteHandler.GetPendingInvites)

	_ = os.Stdout // avoid unused import issues if needed later

	addr := cfg.Addr()
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
