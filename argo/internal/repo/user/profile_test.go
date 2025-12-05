package user

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/swanhubx/swanlab-helper/argo/internal/pkg/mock"
	"github.com/swanhubx/swanlab-helper/argo/internal/pkg/model"
	"gorm.io/gorm"
	"testing"
)

type UserProfileRepoSuite struct {
	suite.Suite

	db   *gorm.DB
	repo ProfileRepo
}

func TestUserProfileRepoSuite(t *testing.T) {
	suite.Run(t, new(UserProfileRepoSuite))
}

func (s *UserProfileRepoSuite) SetupSuite() {
	s.db = mock.NewMySQL()
	s.repo = NewProfileRepo(s.db)
}

func (s *UserProfileRepoSuite) SetupTest() {
	err := s.db.Migrator().DropTable(model.Profile{})
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.db.AutoMigrate(model.Profile{})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *UserProfileRepoSuite) TearDownSuite() {
	conn, _ := s.db.DB()
	_ = conn.Close()
}

// 模拟创建一个用户.

// MockUserProfileRecord 模拟一条用户资料记录.
func MockProfileRecord(db *gorm.DB, t *testing.T) *model.Profile {
	profile := model.Profile{
		Name:   mock.ID(12),
		Avatar: "https://test.swanlab.cn/avatar.png",
		UserID: 123456,
	}
	err := db.Create(&profile).Error
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return &profile
}

// -- FindByUID and Create

// TestFindByUID 查找用户资料.
func (s *UserProfileRepoSuite) TestFindByUID() {
	record := MockProfileRecord(s.db, s.T())
	profile, err := s.repo.FindByUID(context.Background(), record.UserID)
	if err != nil {
		assert.FailNow(s.T(), err.Error())
	}
	s.Equal(record.ID, profile.ID)
}

// TestFindByUID_NotFound 用户资料不存在.
func (s *UserProfileRepoSuite) TestFindByUID_NotFound() {
	profile, err := s.repo.FindByUID(context.Background(), 0)
	s.Error(err)
	s.Nil(profile)
}

// -- Function: Update

// TestUpdate 更新用户资料.
func (s *UserProfileRepoSuite) TestUpdate() {
	record := MockProfileRecord(s.db, s.T())
	record.Name = "Updated Name"
	record.Avatar = "https://test.swanlab.cn/avatar.png.updated"
	err := s.repo.Update(context.Background(), record)
	s.NoError(err)

	// 再次查询，验证更新结果
	updatedProfile, err := s.repo.FindByUID(context.Background(), record.UserID)
	s.NoError(err)
	s.Equal("Updated Name", updatedProfile.Name)
	s.Equal("https://test.swanlab.cn/avatar.png.updated", updatedProfile.Avatar)
}
