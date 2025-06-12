package middlewares

import (
	stderrors "errors"
	"net/http"
	"scalper/errors"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validator(param interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// body, _ := c.GetRawData() body只能被获取一次
		// fmt.Println("Raw body:", string(body))
		var bindErr error
		switch c.Request.Method {
		case "POST", "PUT", "PATCH":
			bindErr = c.ShouldBindJSON(param)
		case "GET":
			bindErr = c.ShouldBind(param)
		default:
			c.Set("app_error", errors.NewAppError(errors.ErrCodeMethodNotAllowed, stderrors.New("Unsupported method"), ""))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if bindErr != nil {
			c.Set("app_error", errors.NewAppError(errors.ErrCodeInvalidParam, bindErr, ""))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		validator := validator.New()
		validator.RegisterValidation("phone", validatePhone)
		if err := validator.Struct(param); err != nil {
			c.Set("app_error", errors.NewAppError(errors.ErrCodeValidation, err, err.Error()))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Set("validated", param)
		c.Next()
	}
}

func validatePhone(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return regexp.MustCompile(`^1[0-9]{10}$`).MatchString(phone)
}
