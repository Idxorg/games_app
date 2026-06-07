package main

import (
	"context"
	"log"
	"os"

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
	// Инициализация БД
	dbPool, err := pgxpool.New(context.TODO(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	// Инициализация Redis
	redisAddr := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

	// Инициализация WebSocket хаб
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(dbPool)
	// TODO: Добавить остальные репозитории

	// Инициализация сервисов
	portalAPI := service.NewPortalAPI(
		os.Getenv("PORTAL_API_URL"),
		os.Getenv("PORTAL_API_KEY"),
	)
	// s3Client := service.NewS3Client(...) // TODO: Инициализировать S3 клиент

	// Инициализация handlers
	authHandler := handler.NewAuthHandler(userRepo, portalAPI)
	userHandler := handler.NewUserHandler(userRepo, nil) // ratingRepo пока nil
	// TODO: Инициализировать остальные handlers

	// Настройка Gin
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(100, 1))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "time": "2026-06-07T12:00:00Z"})
	})

	// API v1
	api := r.Group("/api/v1")
	{
		// Auth
		auth := api.Group("/auth")
		auth.POST("/verify", authHandler.VerifyToken)

		// Users
		users := api.Group("/users")
		users.Use(middleware.Authenticate(dbPool)) // JWT валидация
		{
			users.GET("/:sid/profile", userHandler.GetProfile)
			// TODO: Добавить остальные endpoints
		}

		// TODO: Добавить остальные API endpoints (tournaments, ratings, games)
	}

	// WebSocket endpoint
	r.GET("/ws/game/:game_id", func(c *gin.Context) {
		wsHub.HandleWebSocket(c)
	})

	// Static files (фронтенд)
	r.Static("/static", "./web")
	r.StaticFile("/", "./web/index.html")

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
