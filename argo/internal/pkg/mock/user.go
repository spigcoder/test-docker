package mock

import (
	"context"

	"github.com/swanhubx/swanlab-helper/argo/internal/pkg/model"
)

type UserProfileRepo struct {
	CreateFunc    func(profile *model.Profile) error
	FindByUIDFunc func(uid uint) (*model.Profile, error)
	UpdateFunc    func(profile *model.Profile) error
}

func (m *UserProfileRepo) Create(_ context.Context, profile *model.Profile) error {
	return m.CreateFunc(profile)
}

func (m *UserProfileRepo) FindByUID(_ context.Context, uid uint) (*model.Profile, error) {
	return m.FindByUIDFunc(uid)
}

func (m *UserProfileRepo) Update(_ context.Context, profile *model.Profile) error {
	return m.UpdateFunc(profile)
}
