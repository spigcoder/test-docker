package errs

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
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

// HTTP 响应
func TestError_Response(t *testing.T) {
	// 先模拟一个 gin.Context()
	gin.SetMode(gin.ReleaseMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	Wrap(OK, "I'm fine").WithMeta(Metadata{
		"foo": "bar",
	}).Response(ctx)
	// 查看 ctx 响应
	assert.Equal(t, OK.StatusCode(), r.Code)
	// 对请求体进行解码
	var data gin.H
	err := json.NewDecoder(r.Body).Decode(&data)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, gin.H{
		// 这里转成 float64 是因为 Go 的 JSON 包会将数字默认解析为 float64
		"status":      float64(OK.StatusCode()),
		"code":        OK.String(),
		"message":     "I'm fine",
		"description": "200 - ok: I'm fine",
	}, data)
	// 通过 ctx.Get 获取传递的错误对象
	e, ok := ctx.Get("errs")
	assert.True(t, ok)
	d, _ := e.(Error)
	assert.Equal(t, OK, d.Code)
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
