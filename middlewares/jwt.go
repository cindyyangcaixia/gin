package middlewares

import (
	stderrors "errors"
	"net/http"
	"scalper/errors"
	"scalper/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func JWTAuth(secret string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.Error("Missing Authorization header")
			c.Set("app_error", errors.NewAppError(errors.ErrCodeInvalidToken, stderrors.New("Missing Authorization header"), ""))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Error("Invalid Authorization header")
			c.Set("app_error", errors.NewAppError(errors.ErrCodeInvalidToken, stderrors.New("Invalid Authorization header"), ""))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &services.Claims{}

		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error("unexpected signing method")
				return nil, stderrors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			logger.Error("Invalid token", zap.Error(err))
			c.Set("app_error", errors.NewAppError(errors.ErrCodeInvalidToken, stderrors.New("Invalid token"), ""))
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		c.Set("serial_number", claims.SerialNumber)
		c.Set("phone_number", claims.PhoneNumber)
		c.Set("am_id", claims.AmID)

		c.Next()
	}
}
