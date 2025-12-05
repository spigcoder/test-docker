package errs

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	err := Wrap(Unknown, "I don't know what happened").WithMeta(Metadata{
		"foo": "bar",
	})
	assert.Equal(t, Unknown, err.Code)
	assert.Equal(t, "bar", err.Meta["foo"])
	assert.Equal(t, nil, err.Err)
}

func TestErrCode(t *testing.T) {
	assert.Equal(t, 200, OK.StatusCode())
	assert.Equal(t, "ok", OK.String())
}

// 测试 WithErr 方法
func TestError_WithErr(t *testing.T) {
	err := Error{Code: OK, Message: "ok"}
	assert.Equal(t, nil, err.Err)
	// 使用 WithErr 方法设置 Err 属性
	assert.Equal(t, "mock error", err.WithErr(errors.New("mock error")).Err.Error())
}

// 测试 WithMeta 方法
func TestError_WithMeta(t *testing.T) {
	err := Error{Code: OK, Message: "ok"}
	assert.Equal(t, nil, err.Err)
	// 使用 WithMeta 方法设置 Meta 属性
	assert.Equal(t, "bar", err.WithMeta(Metadata{"foo": "bar"}).Meta["foo"])
}

// Body 方法响应状态码和请求响应体
func TestError_Body(t *testing.T) {
	status, body := Wrap(Unknown, "invalid message").
		WithErr(errors.New("invalid error")).Body()
	assert.Equal(t, Unknown.StatusCode(), status)
	assert.Equal(t, gin.H{
		"status":      Unknown.StatusCode(),
		"code":        Unknown.String(),
		"message":     "invalid message",
		"description": "500 - unknown: invalid message, invalid error",
	}, body)
}

// 返回错误的文本信息
func TestError_Error(t *testing.T) {
	// 没有错误详情
	err := Wrap(OK, "I'm fine").Error()
	assert.Equal(t, "200 - ok: I'm fine", err)
	// 有错误详情
	err = Wrap(OK, "I'm fine").WithErr(errors.New("mock error")).Error()
	assert.Equal(t, "200 - ok: I'm fine, mock error", err)
}
