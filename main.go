package main

import (
	"context"
	"log"
	"scalper/config"
	"scalper/middlewares"
	"scalper/routers"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := config.InitLogger(cfg.LoggerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("starting application")

	// Initialize MongoDB
	client, repos, err := config.InitMongoDB(cfg.MongoConfig, logger)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB: %v", zap.Error(err))
	}
	defer client.Disconnect(context.Background())

	// Initialize UserService
	services := config.InitServices(repos, logger)

	// Initialize Gin
	r := gin.New()
	r.Use(middlewares.RequestID())
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		Context: func(c *gin.Context) []zapcore.Field {
			if id, exists := c.Get("request_id"); exists {
				return []zapcore.Field{
					zap.String("request_id", id.(string)),
				}
			}
			return nil
		},
	}))

	// r.Use(middlewares.Logger(logger))
	r.Use(middlewares.ResponseFormatter(logger))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	routers.SetupRoutes(r, &routers.Services{
		UserPhoneService: services.UserPhoneService,
	}, logger, viper.GetString("JwtSecret"))

	// Run server
	if err := r.Run(":3000"); err != nil {
		logger.Fatal("Failed to start serve %v", zap.Error(err))
	}
}
