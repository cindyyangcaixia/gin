package routers

import (
	"net/http"
	"scalper/middlewares"
	"scalper/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserLoginRequest struct {
	PhoneNumber  string  `json:"phone_number" validate:"phone"`
	SerialNumber *string `json:"serial_number,omitempty" validate:"omitempty"`
	Password     string  `json:"password"`
}

func UserPhoneLoginRouters(r *gin.RouterGroup, svc *services.UserPhoneService, logger *zap.Logger) {
	router := r.Group("user-phones")

	router.POST("login", middlewares.Validator(&UserLoginRequest{}), func(ctx *gin.Context) {
		validated, _ := ctx.Get("validated")
		data := validated.(*UserLoginRequest)
		token, err := svc.Login(ctx, data.PhoneNumber, data.Password)
		if err != nil {
			ctx.Set("app_error", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})
}
