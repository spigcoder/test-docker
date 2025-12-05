package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/swanhubx/swanlab-helper/argo/internal/pkg/mock"
	"github.com/swanhubx/swanlab-helper/argo/internal/pkg/model"
	"github.com/swanhubx/swanlab-helper/argo/pkg/errs"
)

type ProfileHandlerSuite struct {
	suite.Suite

	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestProfileHandlerSuite(t *testing.T) {
	suite.Run(t, new(ProfileHandlerSuite))
}

func (s *ProfileHandlerSuite) SetupTest() {
	gin.SetMode(gin.ReleaseMode)

	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
}

// -- Function: UpdateUserProfile

// 200 - OK.
func (s *ProfileHandlerSuite) TestUpdateUserProfile_OK() {
	mockRepo := &mock.UserProfileRepo{
		UpdateFunc: func(_ *model.Profile) error {
			return nil
		},
		FindByUIDFunc: func(uid uint) (*model.Profile, error) {
			return &model.Profile{
				BaseModel: model.BaseModel{
					ID: 1,
				},
				UserID: uid,
				Name:   "Test User",
				Avatar: "https://test.swanlab.cn/avatar.png",
			}, nil
		},
	}

	handler := NewProfileHandler(mockRepo)

	body, _ := mock.RequestBody(gin.H{
		"name":    "Updated User",
		"avatar":  "https://test.swanlab.cn/avatar2.png",
		"user_id": 1,
	})

	s.ctx.Request = httptest.NewRequest(http.MethodPut, "/user/profile", body)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	handler.UpdateUserProfile(s.ctx)

	s.Equal(errs.OK.StatusCode(), s.recorder.Code)
}

// 400 - BadRequest.
func (s *ProfileHandlerSuite) TestUpdateUserProfile_BadRequest() {
	handler := NewProfileHandler(&mock.UserProfileRepo{})
	badCases := []gin.H{
		{"name": "name", "avatar": 123, "user_id": 1},
		{"name": "name", "avatar": "invalid-url"},
		{"name": 123, "avatar": "https://test.swanlab.cn/avatar2.png"},
		{"name": "", "avatar": ""},
	}

	for _, bodyData := range badCases {
		body, _ := mock.RequestBody(bodyData)

		s.ctx.Request = httptest.NewRequest(http.MethodPut, "/user/profile", body)
		s.ctx.Request.Header.Set("Content-Type", "application/json")
		handler.UpdateUserProfile(s.ctx)

		s.Equal(errs.BadRequest.StatusCode(), s.recorder.Code)
	}
}

func (s *ProfileHandlerSuite) TestUpdateUserProfile_NotFound() {
	mockRepo := &mock.UserProfileRepo{
		FindByUIDFunc: func(_ uint) (*model.Profile, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	handler := NewProfileHandler(mockRepo)

	body, _ := mock.RequestBody(gin.H{
		"name":    "Updated User",
		"avatar":  "https://test.swanlab.cn/avatar2.png",
		"user_id": 1,
	})

	s.ctx.Request = httptest.NewRequest(http.MethodPut, "/user/profile", body)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	handler.UpdateUserProfile(s.ctx)

	s.Equal(errs.NotFound.StatusCode(), s.recorder.Code)
}

// 500 - DatabaseError.
func (s *ProfileHandlerSuite) TestUpdateUserProfile_DatabaseError() {
	mockRepo := &mock.UserProfileRepo{
		FindByUIDFunc: func(_ uint) (*model.Profile, error) {
			return nil, gorm.ErrInvalidData
		},
	}

	handler := NewProfileHandler(mockRepo)

	body, _ := mock.RequestBody(gin.H{
		"name":    "Updated User",
		"avatar":  "https://test.swanlab.cn/avatar2.png",
		"user_id": 1,
	})

	s.ctx.Request = httptest.NewRequest(http.MethodPut, "/user/profile", body)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	handler.UpdateUserProfile(s.ctx)

	s.Equal(errs.DatabaseError.StatusCode(), s.recorder.Code)
}
