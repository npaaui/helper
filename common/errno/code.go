package errno

import (
	"fmt"
	"sync"
)

/**
错误码
*/
type CodeErr struct {
	HttpStatus int
	Code       int
	Msg        string
	Err        error
}

func (c CodeErr) GetHttpStatus() int {
	return c.HttpStatus
}

func (c CodeErr) GetCode() int {
	return c.Code
}

func (c CodeErr) GetMsg() string {
	return c.Msg
}

func (c CodeErr) Error() string {
	return c.Msg
}

func (c CodeErr) NotNil() bool {
	return c.Code != 0
}

const (
	Success     = 0      // 成功
	ErrUnknown  = -1     // 未知错误
	ErrNetwork  = 100000 // 网络异常
	ErrSystem   = 100001 // 系统异常
	ErrService  = 100002 // 服务异常
	ErrValid    = 100003 // 校验异常
	ErrDatabase = 100004 // 数据库异常
	ErrConfig   = 100005 // 配置异常
)

var codes = map[int]CodeErr{} // 初始化时存储所有注册的错误码
var codesMux sync.Mutex

func init() {
	Register(200, Success, "成功")
	Register(500, ErrUnknown, "未知错误")
	Register(500, ErrNetwork, "网络异常，请稍后重试")
	Register(500, ErrSystem, "系统异常")
	Register(500, ErrService, "服务异常")
	Register(500, ErrValid, "校验异常")
	Register(500, ErrDatabase, "数据库异常")
	Register(500, ErrConfig, "配置异常")
}

// 注册错误码
func Register(httpStatus, code int, msg string) {
	codesMux.Lock()
	defer codesMux.Unlock()
	codes[code] = CodeErr{httpStatus, code, msg, nil}
}

func NewCodeErr(code int) CodeErr {
	if codeErr, ok := codes[code]; ok {
		return codeErr
	}
	return codes[ErrUnknown]
}

func NewCodeErrAndMsg(code int, msg string) CodeErr {
	codeErr := NewCodeErr(code)
	if msg != "" {
		codeErr.Msg = msg
	}
	return codeErr
}

func NewCodeErrAndErr(code int, format string, args ...interface{}) CodeErr {
	codeErr := NewCodeErr(code)
	if format != "" {
		codeErr.Err = fmt.Errorf(format, args)
		codeErr.Msg = codeErr.Err.Error()
	}
	return codeErr
}
