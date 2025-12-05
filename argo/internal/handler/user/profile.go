package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userApi "github.com/swanhubx/swanlab-helper/argo/api/v1/user"
	userRepo "github.com/swanhubx/swanlab-helper/argo/internal/repo/user"
	"github.com/swanhubx/swanlab-helper/argo/pkg/errs"
)

type ProfileHandler struct {
	repo userRepo.ProfileRepo
}

func NewProfileHandler(repo userRepo.ProfileRepo) ProfileHandler {
	return ProfileHandler{repo: repo}
}

// UpdateUserProfile 修改用户信息
// API: PUT /user/profile
// Response:
// - 200 OK				修改用户信息成功
// - 400 BadRequest		请求参数错误
// - 401 Unauthorized	未认证
// - 404 NotFound		用户信息不存在
// - 500 DatabaseError	数据库错误
// - 500 Unknown		未知错误.
func (h *ProfileHandler) UpdateUserProfile(ctx *gin.Context) {
	var req userApi.UpdateUserProfileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		errs.Wrap(errs.BadRequest, "Request format error").Response(ctx)
		return
	}

	if req.Name == "" && req.Avatar == "" {
		errs.Wrap(errs.BadRequest, "No field to update").Response(ctx)
		return
	}

	profile, err := h.repo.FindByUID(ctx, req.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		errs.Wrap(errs.DatabaseError, "Failed to find user profile").WithErr(err).Response(ctx)
		return
	}
	if profile == nil {
		errs.Wrap(errs.NotFound, "User profile not found").Response(ctx)
		return
	}

	if req.Name != "" {
		profile.Name = req.Name
	}
	if req.Avatar != "" {
		profile.Avatar = req.Avatar
	}
	if err := h.repo.Update(ctx, profile); err != nil {
		errs.Wrap(errs.DatabaseError, "Failed to update user profile").WithErr(err).Response(ctx)
		return
	}

	ctx.JSON(http.StatusOK, &userApi.UpdateUserProfileResponse{
		Name:   profile.Name,
		Avatar: profile.Avatar,
	})
}
