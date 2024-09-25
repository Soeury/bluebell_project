package controllers

type ResCode int64

// 自定义的一些错误码
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
)

// 将错误码与错误信息对应
var CodeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "wrong query params",
	CodeUserExist:       "user exist",
	CodeUserNotExist:    "user not exist",
	CodeInvalidPassword: "wrong password",
	CodeServerBusy:      "busy server",
	CodeNeedLogin:       "need login",
	CodeInvalidToken:    "wrong token",
}

// Msg 方法实现通过错误码取出错误信息
func (code ResCode) Msg() string {
	msg, ok := CodeMsgMap[code]
	if !ok {
		return CodeMsgMap[CodeServerBusy] // 除自定义错误码之外的错误信息都是 [busy service]
	}
	return msg
}
