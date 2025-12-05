package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userHandler "github.com/swanhubx/swanlab-helper/argo/internal/handler/user"
	userRepo "github.com/swanhubx/swanlab-helper/argo/internal/repo/user"
)

// Group 代表实现一个路由组的相关方法
type Group interface {
	// AddRoutes 用于增加路由，传递一个上层路由组和资源标识
	AddRoutes(api *gin.RouterGroup)
}

type UserRouter struct {
	db *gorm.DB
}

func NewUserRouter(db *gorm.DB) Group {
	return &UserRouter{db: db}
}

// AddRoutes 增加用户相关的接口
// - PUT 	/user/profile	更新用户配置.
func (r *UserRouter) AddRoutes(api *gin.RouterGroup) {
	userProfileRepo := userRepo.NewProfileRepo(r.db)

	uh := userHandler.NewProfileHandler(userProfileRepo)

	userApis := api.Group("/user")
	{
		userApis.PUT("/profile", uh.UpdateUserProfile)
	}
}
