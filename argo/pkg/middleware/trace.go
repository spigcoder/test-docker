package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	TraceIDHeader        = "X-Trace-ID"
	TraceIDKey    ctxKey = "trace_id"
)

// TraceMiddleware 中间件
func TraceMiddleware(c *gin.Context) {
	traceID := c.GetHeader(TraceIDHeader)
	// 设置 Response 中的 X-Trace-ID
	if traceID != "" {
		c.Header(TraceIDHeader, traceID)
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)
	}
	c.Next()
}
