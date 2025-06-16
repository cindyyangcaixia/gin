package routers

import (
	"net/http"
	"scalper/middlewares"
	"scalper/models"
	"scalper/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type UserPhoneGetByPhone struct {
	PhoneNumber string `uri:"phone_number" validate:"required,phone"`
}

type UserPhoneQuery struct {
	PhoneNumber  string  `form:"phone_number,omitempty" validate:"omitempty,phone"`
	SerialNumber *string `form:"serial_number,omitempty" validate:"omitempty"`
	Page         int64   `form:"page" validate:"gte=1"`
	Limit        int64   `form:"limit" validate:"gte=1"`
}

func UserPhoneRouters(r *gin.RouterGroup, svc *services.UserPhoneService, logger *zap.Logger) {
	router := r.Group("user-phones")

	router.POST("", middlewares.Validator(&models.UserPhone{}), func(ctx *gin.Context) {
		validated, _ := ctx.Get("validated")
		userPhone := validated.(*models.UserPhone)
		res, err := svc.CreateUserPhone(ctx.Request.Context(), userPhone)
		if err != nil {
			ctx.Set("app_error", err)
			ctx.Abort()
			return
		}
		userPhone.ID = res.InsertedID.(primitive.ObjectID)
		ctx.JSON(http.StatusCreated, userPhone)
	})

	router.GET(":phone_number", middlewares.Validator(&UserPhoneGetByPhone{}), func(ctx *gin.Context) {
		validated, _ := ctx.Get("validated")
		data := validated.(*UserPhoneGetByPhone)
		phone, err := svc.GetUserPhone(ctx, data.PhoneNumber)
		if err != nil {
			ctx.Set("app_error", err)
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, phone)
	})

	router.GET("", middlewares.Validator(&UserPhoneQuery{}), func(ctx *gin.Context) {
		validated, _ := ctx.Get("validated")
		params := validated.(*UserPhoneQuery)
		items, total, err := svc.ListUserPhones(ctx.Request.Context(), params.PhoneNumber, params.SerialNumber, params.Page, params.Limit)
		if err != nil {
			ctx.Set("app_error", err)
			ctx.Abort()
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"items": items,
			"total": total,
			"page":  params.Page,
			"limit": params.Limit,
		})
	})
}
