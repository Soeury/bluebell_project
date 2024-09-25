package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	CtxUserIDKey      = "user_id"
	ErrorUserNotLogin = errors.New("user not login")
)

// GetUserID 获取当前用户的 ID
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	// 注意: 这里c.Get() 取出来的值是一个 interface{} 任意类型的接口值，需要转换成 int64 类型才能使用
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err := ErrorUserNotLogin
		return 0, err
	}

	userID, ok = uid.(int64)
	if !ok {
		err := ErrorUserNotLogin
		return 0, err
	}
	return userID, nil
}

// GetPageAndSize 获取分页参数
func GetPageAndSize(c *gin.Context) (int64, int64) {

	// 获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1 // 设置默认值
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10 // 默认值
	}
	return page, size
}
