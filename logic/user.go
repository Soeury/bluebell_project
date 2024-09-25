package logic

import (
	"project_bluebell/dao/mysql"
	"project_bluebell/models"
	"project_bluebell/pkg/jwt"
	"project_bluebell/pkg/snowflake"
)

// SignUp 实现用户注册
func SignUp(p *models.ParamSignUp) error {

	// 1. 判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2. 生成 uid
	UserID := snowflake.GetId()

	// 构造一个用户实例
	user := &models.User{
		UserID:   UserID,
		Username: p.Username,
		Password: p.Password,
	}

	// 4. 保存到数据库
	return mysql.InsertUser(user)
}

// Login 实现用户登录
func Login(p *models.ParamLogin) (token string, err error) {

	// 检验用户是否存在在数据库中 , 这里将 models.ParamLogin 转化成 models.User 类型的对象
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 下面传入的 user 指针已经在 Login 内部修改过值了，可以拿到 user.UserID 的值
	if err := mysql.Login(user); err != nil {
		return "", err
	}

	return jwt.GenToken(user.UserID, user.Username)
}
