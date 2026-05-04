package main

import (
	"fmt"
	"log"
	"platform/go-pkg/database"
	"platform/go-pkg/jwt"
	"platform/go-pkg/logger"
	"platform/go-pkg/middleware"
	"platform/go-pkg/validator"
	"strconv"
	"todo/config"
	"todo/controllers"
	"todo/handlers"
	"todo/repositories"
	"todo/routers"

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

	todoRepo := repositories.NewTodoRepository(db)
	tagRepo := repositories.NewTagRepository(db)

	jwt, err := jwt.NewService(jwt.Config{
		Secret: cfg.JWTSecret,
		Issuer: cfg.AppName,
	})
	if err != nil {
		log.Fatalf("failed jwt service: %v", err)
	}

	todoController := controllers.NewTodoController(todoRepo, tagRepo, jwt)

	v := validator.New()

	todoHandler := handlers.NewTodoHandler(todoController, v)

	// setup gin
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.UseRawPath = true

	r.Use(middleware.CORS())

	api := r.Group("/api")

	routers.RegisterTodoRoutes(api, todoHandler, middleware.Auth(middleware.AuthConfig{
		SecretKey: cfg.JWTSecret,
		Logger: logger.New("todo"),
	}))

	log.Printf("todo service running on port %s", cfg.AppPort)
	if err := r.Run(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}