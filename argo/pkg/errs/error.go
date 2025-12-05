package errs

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// 定义错误结构体应该考虑三个角色：应用程序、终端用户和运维人员，不同角色的用户对于错误有不同的理解
// - 应用程序需要的是固定的错误编码，便于对错误做出对应的动作
// - 终端用户则需要人类可读的语言便于理解错误内容
// - 运维人员需要的是更细节的信息，包括一些入参和错误堆栈，便于排查错误
//
// 参考：https://blog.carlana.net/post/2020/working-with-errors-as/

// Metadata 用于记录结构化的键值对，用于将任意元数据附加到错误，
// 元数据只在内部传递，不对客户端公开，对于运维人员友好
type Metadata map[string]interface{}

// Error 为错误结构体
// 参考自：https://encore.dev/docs/go/primitives/api-errors
type Error struct {
	// Code 为错误码，便于应用读取
	Code ErrCode `json:"code"`

	// Message 为人类可读的响应信息
	Message string `json:"message"`

	// Err 代表详细的错误信息，包含一些错误堆栈
	Err error `json:"error"`

	// Meta 为键值对信息，包含入参和一些资源 ID 信息，运维人员友好，
	// 程序内部使用，不对外暴露
	Meta Metadata `json:"-"`
}

// WithErr 方法用于增加详细的错误信息，返回一个结构体指针
func (e *Error) WithErr(err error) *Error {
	e.Err = err
	return e
}

// 上述的链式调用方法可能存在一些问题，如下：
//
// 这里是直接修改原对象的方法
// ```
// func (e *Error) WithErr(err error) *Error {
//     e.Err = err    // 直接修改原对象
// 	   return e
// }
// ```
//
// 使用时可能产生的问题
// ```
// err1 := Wrap(errs.NotFound, "not found")
// err2 := err1.WithErr(someError)  // err1 和 err2 指向同一个对象
// ```
//
// 因此还有一种做法是直接创建一个新的对象，如下：
// ```
// func (e *Error) WithErr(err error) *Error {
//     return &Error{
//        Code:    e.Code,
//		  Message: e.Message,
//		  Err: err,
//		  Meta: e.Meta
//     }
// }
// ```

// WithMeta 用于携带结构化的键值对，便于定位错误
func (e *Error) WithMeta(fields Metadata) *Error {
	e.Meta = fields
	return e
}

// Body 方法返回 gin 中 ctx.JSON() 方法所需要的参数，
// 分别是一个响应状态码和一个 JSON 响应体。固定的响应体格式为：
//
//	{
//	  "status": 200,
//	  "code": "OK",
//	  "message": "created successfully"
//	}
//
// 其它信息不会响应给客户端，只用于系统内部排查和定位错误
func (e *Error) Body() (int, gin.H) {
	status := e.Code.StatusCode()
	return status, gin.H{
		"status":      status,
		"code":        e.Code.String(),
		"message":     e.Message,
		"description": e.Error(),
	}
}

// Error 用于返回一个错误文本信息，格式为 {status} - {code}: {message}, {detail}，
// 例如：400 - not_found: the user is not found
func (e *Error) Error() string {
	detail := ""
	// 包含错误详细信息
	if e.Err != nil {
		detail = ", " + e.Err.Error()
	}
	return fmt.Sprintf("%d - %s: %s%s", e.Code.StatusCode(), e.Code.String(), e.Message, detail)
}

// Response 用于接口响应，这里传入一个 gin.Context 上下文，
// 然后会将 Error 对象通过 ctx.Set() 传递给后面的中间件用于打印错误信息，设置属性为 errs，
// 使用 ctx.JSON() 方法设置响应状态码和响应体
func (e *Error) Response(ctx *gin.Context) {
	slog.ErrorContext(ctx, e.Error())
	// ctx.Set("errs", &e)
	err := ctx.Error(e)
	if err != nil {
		slog.Error("Can't Set error to gin.Context")
	}
	ctx.JSON(e.Body())
}

// Wrap 用于组装一个新的 Error 对象，需要传递错误码和错误信息，
// 如果需要传递额外的信息，例如 Err 和 Meta，则需要使用方法 WithErr 和 WithMeta,
// 使用方式例如 `errs.Wrap(errs.OK, "It's ok").WithErr(errors.New("hello"))`,
// 使用链式调用的方法。在应用中都需要使用 errs.Wrap() 来包装一个错误类型。
func Wrap(code ErrCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
