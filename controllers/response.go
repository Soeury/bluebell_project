package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 封装一些返回响应的函数
/*
    code :  错误码
	msg  :  错误信息
	data :  相关数据
*/

// 这些方法都要传入 *gin.Context 参数，所以都是在 controllers 层调用的，用来返回响应
type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `josn:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseErr(c *gin.Context, code ResCode) {

	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseErrWithMsg(c *gin.Context, code ResCode, msg interface{}) {

	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {

	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
