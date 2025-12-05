package errs

// 一个错误码具有三个属性：错误编码、错误码名称、HTTP 状态码。
// 错误码不能被程序任意修改，所以应该定义为常量而非变量，但由于其还具有其它额外的属性，错误码名称和对应的 HTTP 状态码，
// 所以将其对应为一个 int 类型，代表其对应的索引，然后将其它属性定义在数组中。
//
// 错误码的定义方法参考自：https://github.com/encoredev/encore.dev/blob/v1.46.1/beta/errs/codes.go

// ErrCode 为错误编码
type ErrCode int

// 错误码枚举值，此为定义错误码对应的编号
// Note：每添加一个错误码都需要在 codeNames 和 codeStatus 中定义对应的值
const (
	// OK 代表操作成功，对应 200 OK
	OK ErrCode = 0

	// Unknown 代表服务未知的错误，对应 500 Internal Server Error
	Unknown ErrCode = 1

	// BadRequest 代表传入的参数无效，对应 400 Bad Request
	BadRequest ErrCode = 2

	// DatabaseError 代表数据库操作错误，使用 500 Internal Server Error
	DatabaseError ErrCode = 3

	// NotFound 代表为资源不存在，对应 404 Not Found
	NotFound ErrCode = 4

	// Unauthorized 代表缺乏身份验证凭证，不允许接下来的操作，对应 401 Unauthorized
	Unauthorized ErrCode = 5

	// Forbidden 代表拒绝访问，对应 403 Forbidden
	Forbidden ErrCode = 6

	// TooManyRequests 代表请求过多，对应 429 Too Many Requests
	TooManyRequests ErrCode = 7

	// Conflict 代表资源冲突，对应 409 Conflict
	Conflict ErrCode = 8

	// WebsocketError 代表 websocket 错误
	WebsocketError ErrCode = 9

	// QueueError 代表 actor 消息发送失败
	QueueError ErrCode = 10
)

// String 返回错误编码实际对应的错误码
func (c ErrCode) String() string {
	return codeNames[c]
}

// StatusCode 用于返回错误码对应的 HTTP 状态码，对于 API 响应非常友好，
// 状态码代表的是一类错误，一个状态码下可能有多个错误码
func (c ErrCode) StatusCode() int {
	return codeStatus[c]
}

// [...] 表示由编译器根据初始化值的数量来确定数组长度
// codeNames 为错误码的实际名称，用于返回给应用
var codeNames = [...]string{
	OK:              "ok",
	Unknown:         "unknown",
	BadRequest:      "bad_request",
	DatabaseError:   "database_error",
	NotFound:        "not_found",
	Unauthorized:    "unauthorized",
	Forbidden:       "forbidden",
	TooManyRequests: "too_many_requests",
	Conflict:        "conflict",
	WebsocketError:  "websocket_error",
	QueueError:      "queue_error",
}

// codeStatus 为错误码对应的 HTTP 状态码
// 响应状态码参考：https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status
var codeStatus = [...]int{
	OK:              200,
	Unknown:         500,
	BadRequest:      400,
	DatabaseError:   500,
	NotFound:        404,
	Unauthorized:    401,
	Forbidden:       403,
	TooManyRequests: 429,
	Conflict:        409,
	WebsocketError:  500,
	QueueError:      500,
}
