package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 为全局错误处理中间件，会在控制台打印错误信息
func ErrorHandler(ctx *gin.Context) {
	ctx.Next()
	err := ctx.Errors.Last()
	if err != nil {
		slog.ErrorContext(ctx.Request.Context(), err.Error(),
			slog.String("event", "api_error"),
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.RequestURI),
			slog.Int("status", ctx.Writer.Status()),
		)
	}
}
