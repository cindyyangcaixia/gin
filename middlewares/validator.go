package middlewares

import (
	"net/http"
	"scalper/errors"

	stderrors "github.com/pkg/errors"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Validator(param interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// body, _ := c.GetRawData() body只能被获取一次
		// fmt.Println("Raw body:", string(body))
		methods := map[string]bool{
			"POST":  true,
			"PUT":   true,
			"PATCH": true,
		}
		// Modify it when file upload is needed.
		if methods[c.Request.Method] && c.ContentType() != binding.MIMEJSON {
			c.Set("app_error", errors.NewAppError(errors.ErrCodeInvalidParam, http.StatusBadRequest, stderrors.New("only support json "), "only support json"))
			c.Abort()
			return
		}

		if err := c.ShouldBindUri(param); err != nil {
			c.Set("app_error", errors.NewAppError(errors.ErrCodeInvalidParam, http.StatusBadRequest, stderrors.New(err.Error()), ""))
			c.Abort()
			return
		}

		if err := c.ShouldBind(param); err != nil {
			c.Set("app_error", errors.NewAppError(errors.ErrCodeInvalidParam, http.StatusBadRequest, stderrors.New(err.Error()), ""))
			c.Abort()
			return
		}

		validator := validator.New()
		validator.RegisterValidation("phone", validatePhone)
		if err := validator.Struct(param); err != nil {
			c.Set("app_error", errors.NewAppError(errors.ErrCodeValidation, http.StatusBadRequest, stderrors.New(err.Error()), err.Error()))
			c.Abort()
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
