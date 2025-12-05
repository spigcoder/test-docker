package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/swanhubx/swanlab-helper/argo/pkg/middleware"
)

// Init 初始化全局 JSON Logger
// levelStr: 日志级别 (debug, info, warn, error)
func Init(levelStr string) {
	// 1. 解析日志级别 (默认为 Info)
	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 2. 配置 Handler 选项
	opts := &slog.HandlerOptions{
		Level:     level, // 设置最低日志级别
		AddSource: true,  // 【可选】在日志中显示文件名和行号 (source: "main.go:20")
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	finalHandler := NewTraceHandler(handler)

	logger := slog.New(finalHandler)
	slog.SetDefault(logger)
}

type TraceHandler struct {
	slog.Handler
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID, ok := ctx.Value(middleware.TraceIDKey).(string); ok {
		r.AddAttrs(slog.String("trace_id", traceID))
	}
	return h.Handler.Handle(ctx, r)
}

// NewTraceHandler 创建带 Trace 能力的 Handler
func NewTraceHandler(h slog.Handler) *TraceHandler {
	return &TraceHandler{Handler: h}
}
