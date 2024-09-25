package controllers

import (
	"errors"
	"fmt"
	"project_bluebell/dao/mysql"
	"project_bluebell/logic"
	"project_bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GetPostDetailHandler
// @Summary 实现用户注册功能
// @Description 实现用户注册网站
// @Tags user_group
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamSignUp true "signup : username , password , re_password"
// @Success 200 {object} ResNil_
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	// SignUpHandler 实现注册功能
	// 1. 获取参数和参数校验 (创建数据，返回错误)    -1. query    -2. params √
	// 这里创建一个指针类型的对象
	p := new(models.ParamSignUp)

	// *在web项目中对请求的参数进行校验是非常常见的行为，可以考虑使用库validator来提高效率,
	// ShouldBindJson 只能简单的识别一下传入的数据类型是否正确，并且是否为 JSON 格式，无法处理键值为空的情况
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Sign up with invalid param", zap.Error(err))

		// 判断 err 是不是 validator.ValidationErrors 类型
		// 比如说数据的类型和结构体的类型对应不上，这时候err就是  validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseErr(c, CodeInvalidParam)
			return
		}
		ResponseErrWithMsg(c, CodeInvalidParam, RemoveTopStruct(errs.Translate(trans)))
		return
	}
	// 传入数据正确，打印输出
	fmt.Printf("user: %+v\n", *p)

	// 2. 业务处理   -   这里交给 logic 去实现  这里最好是使用结构体指针的方式传入数据
	// 最底层的处理将错误返回到上一层，上一层继续返回到再上一层，这时，最上层拿到的错误就是最底层传上来的
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseErr(c, CodeUserExist)
			return
		}
		ResponseErr(c, CodeServerBusy)
		return
	}

	// 3. 返回响应 (使用上面的数据，打印信息)
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler
// @Summary 实现用户登录
// @Description 实现用户进行网站登录的功能
// @Tags user_group
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamLogin true "login: username , password"
// @Success 200 {object} ResToken_
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	// LoginHandler 实现登录功能
	// 1. 获取参数 参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断 err 是不是 validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseErr(c, CodeInvalidParam)
			return
		}
		ResponseErrWithMsg(c, CodeInvalidParam, RemoveTopStruct(errs.Translate(trans)))
		return
	}
	fmt.Printf("%+v\n", *p)

	// 2. 处理业务逻辑
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("login failed ...", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseErr(c, CodeUserNotExist)
		}
		ResponseErr(c, CodeInvalidPassword)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, token)
}
