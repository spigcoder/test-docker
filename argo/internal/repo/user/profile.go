package user

import (
	"context"

	"gorm.io/gorm"

	"github.com/swanhubx/swanlab-helper/argo/internal/pkg/model"
)

type ProfileRepo interface {
	// Create 创建用户信息
	Create(ctx context.Context, profile *model.Profile) error
	// FindByUID 获取用户信息
	FindByUID(ctx context.Context, uid uint) (*model.Profile, error)
	// Update 更新用户信息
	Update(ctx context.Context, profile *model.Profile) error
}

type userProfileRepo struct {
	db *gorm.DB
}

func NewProfileRepo(db *gorm.DB) ProfileRepo {
	return &userProfileRepo{db: db}
}

// Create 创建用户信息.
func (r *userProfileRepo) Create(ctx context.Context, profile *model.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// FindByUID 查找用户信息.
func (r *userProfileRepo) FindByUID(ctx context.Context, uid uint) (*model.Profile, error) {
	var profile model.Profile
	// 需要用双引号包裹字段名, 防止 gorm 自动转为纯小写形式
	res := r.db.WithContext(ctx).Where(`userId = ?`, uid).First(&profile)
	if res.Error != nil {
		return nil, res.Error
	}

	return &profile, nil
}

// Update 更新用户信息.
func (r *userProfileRepo) Update(ctx context.Context, profile *model.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}
