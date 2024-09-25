package mysql

import "errors"

// 这里是定义的一些常量，需要的时候调用即可，最好不要在程序中出现莫名的字符串
var (
	ErrorUserExist     = errors.New("the user is already exist")
	ErrorUserNotExist  = errors.New("the user is not exist")
	ErrorWrongPassword = errors.New("wrong password")
	ErrorInvalidID     = errors.New("wrong ID")
)

var (
	WarnNoRows = "no rows in community db"
)
