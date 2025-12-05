package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/swanhubx/swanlab-helper/argo/pkg/middleware"
)

// Env 环境结构体，传递给app的入口参数.
type Env struct {
	DB *gorm.DB
}

func NewRouter(env Env) *gin.Engine {
	r := gin.Default()
	r.ContextWithFallback = true
	r.Use(middleware.ErrorHandler)
	r.Use(middleware.TraceMiddleware)

	api := r.Group("/" + viper.GetString("api_prefix"))
	NewUserRouter(env.DB).AddRoutes(api)
	return r
}
