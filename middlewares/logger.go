package middlewares

import (
	"fmt"
	"scalper/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// start := time.Now()
		// path := ctx.Request.URL.Path
		// method := ctx.Request.Method
		// clientIp := ctx.ClientIP()

		ctx.Next()

		// latency := time.Since(start)
		// status := ctx.Writer.Status()

		requestID := ""
		if id, exists := ctx.Get("request_id"); exists {
			requestID = id.(string)
		}

		if appErr, exists := ctx.Get("app_error"); exists {
			if err, ok := appErr.(*errors.AppError); ok {
				// if len(ctx.Errors) > 0 {
				// 	for _, err := range ctx.Errors {
				errWithStack := fmt.Sprintf("%+v", err.Err)
				logger.Error("API 请求错误",
					zap.String("request_id", requestID),
					zap.String("error", errWithStack),
				)
			}
		}

		// logger.Info("Http request",
		// 	zap.String("method", method),
		// 	zap.String("path", path),
		// 	zap.Int("status", status),
		// 	zap.Duration("latency", latency),
		// 	zap.String("client_ip", clientIp),
		// 	zap.Any("errors", errWithStack),
		// )
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
