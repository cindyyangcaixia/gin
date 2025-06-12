package routers

import (
	"net/http"
	"scalper/middlewares"
	"scalper/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Services struct {
	UserPhoneService *services.UserPhoneService
}

func SetupRoutes(r *gin.Engine, services *Services, logger *zap.Logger, jwtSecret string) {
	api := r.Group("/api/v1")

	api.GET("version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"version": "1.0.0"})
	})

	UserPhoneLoginRouters(api, services.UserPhoneService, logger)
	protected := api.Group("")
	protected.Use(middlewares.JWTAuth(jwtSecret, logger))
	UserPhoneRouters(protected, services.UserPhoneService, logger)
}
