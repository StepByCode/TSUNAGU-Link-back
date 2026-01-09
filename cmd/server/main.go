package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/StepByCode/TSUNAGU-Link-back/internal/config"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/handler"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/repository"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// デバッグ情報を出力（パスワードは除く）
	log.Printf("Connecting to database: host=%s port=%d user=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.DBName, cfg.Database.SSLMode)

	db, err := sql.Open("postgres", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// データベース接続をリトライ（最大30秒）
	maxRetries := 10
	retryInterval := 3 * time.Second
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err != nil {
			if i == maxRetries-1 {
				log.Fatalf("Failed to ping database after %d attempts: %v", maxRetries, err)
			}
			log.Printf("Database not ready, retrying in %v... (attempt %d/%d)", retryInterval, i+1, maxRetries)
			time.Sleep(retryInterval)
			continue
		}
		break
	}

	log.Println("Successfully connected to database")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	userHandler.RegisterRoutes(e)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
