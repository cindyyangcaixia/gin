package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"scalper/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return len(b), nil
	// return w.ResponseWriter.Write(b)
}

func writeJsonResponse(ctx *gin.Context, logger *zap.Logger, status int, code int, message string, data interface{}, writer gin.ResponseWriter) {
	formatted := gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}

	// logger.Info(
	// 	"Preparing to write JSON response",
	// 	zap.Int("status", status),
	// 	zap.Any("formatted", formatted),
	// 	zap.Any("headers", writer.Header()),
	// 	zap.Bool("written", writer.Written()),
	// )
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(formatted); err != nil {
		logger.Error("Failed to encode response", zap.Error(err), zap.Int("code", code))
		return
	}

	// if flusher, ok := writer.(http.Flusher); ok {
	// 	flusher.Flush()
	// }

	logger.Info("response formatted", zap.Int("status", status), zap.Int("code", code), zap.String("message", message))

	logger.Info(
		"Response",
		zap.String("level", "info"),
		zap.Int("status", status),
		zap.Int("code", code),
		zap.String("message", message),
		zap.Any("headers", writer.Header()),
	)
}

func ResponseFormatter(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		writer := &responseWriter{body: bytes.NewBuffer(nil), ResponseWriter: ctx.Writer}
		ctx.Writer = writer

		ctx.Next()

		if appErr, exists := ctx.Get("app_error"); exists {
			if err, ok := appErr.(*errors.AppError); ok {
				ctx.Error(err.Err)
				writeJsonResponse(ctx, logger, ctx.Writer.Status(), err.Code, err.Message, nil, writer.ResponseWriter)
				return
			}

			writeJsonResponse(ctx, logger, http.StatusInternalServerError, errors.ErrCodeInternalServer,
				errors.ErrorMessages[errors.ErrCodeInternalServer], nil, writer.ResponseWriter)
			return
		}

		// todo
		if ctx.Writer.Status() >= http.StatusBadRequest {
			return
		}

		var data interface{}
		if err := json.Unmarshal(writer.body.Bytes(), &data); err != nil {
			writeJsonResponse(ctx, logger, http.StatusInternalServerError, errors.ErrCodeInternalServer,
				errors.ErrorMessages[errors.ErrCodeInternalServer], nil, writer.ResponseWriter)
			return
		}

		writeJsonResponse(ctx, logger, ctx.Writer.Status(), 0, "Success", data, writer.ResponseWriter)
	}
}
