package main

import (
	"auth/config"
	"auth/controllers"
	"auth/handlers"
	"auth/repositories"
	"auth/routers"
	"fmt"
	"log"
	"platform/go-pkg/database"
	"platform/go-pkg/jwt"
	"platform/go-pkg/logger"
	"platform/go-pkg/middleware"
	"platform/go-pkg/validator"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// load configuration
	cfg := config.Load()

	// connecting database
	maxConns, err := strconv.Atoi(cfg.DBMaxOpenConns)
	if err != nil {
		log.Fatalf("failed parse string to int: %v", err)
	}

	db, err := database.NewPostgres(database.PostgresConfig{
		Host: cfg.DBHost,
		Port: cfg.DBPort,
		User: cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
		SSLMode: cfg.DBSSLMode,
		MaxConns: int32(maxConns),
	})

	if err != nil {
		log.Fatalf("failed connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)

	jwt, err := jwt.NewService(jwt.Config{
		Secret: cfg.JWTSecret,
		Issuer: cfg.AppName,
	})
	authController := controllers.NewAuthController(userRepo, jwt)

	v := validator.New()

	authHandler := handlers.NewAuthHandler(authController, v)

	// setup gin
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.UseRawPath = true

	r.Use(middleware.CORS())

	api := r.Group("/api")

	routers.RegisterAuthRoutes(api, authHandler, middleware.Auth(middleware.AuthConfig{
		SecretKey: cfg.JWTSecret,
		Logger: logger.New("auth"),
	}))

	log.Printf("auth service running on port %s", cfg.AppPort)
	if err := r.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}