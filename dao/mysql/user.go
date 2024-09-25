package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"project_bluebell/models"
)

// 把每一步数据库操作封装成函数
//	让其他的包来调用这些函数即可

const secret = "encryption"

// InsertUser 像数据库中插入一条新的数据
func InsertUser(user *models.User) (err error) {

	// 对密码进行加密(密码不能以明文的方式存储在数据库中)
	user.Password = EncryptPassword(user.Password)

	// 执行SQL 将数据保存入库
	sqlStr := "insert into user(user_id , username , password) values(?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return err
}

// CheckUserExist 查询用户是否存在
func CheckUserExist(username string) (err error) {

	// 注意这里的 SQL 语句 : 不要写习惯了 name 和 age ! 这里是 username!
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}

	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// encryptPassword 对密码进行加密 ? ? ? ! ! !
func EncryptPassword(opassword string) string {

	combined := fmt.Sprintf("%s%s", secret, opassword) // 将 secret 和 opassword 拼接
	h := md5.New()                                     // 创建一个新的 md5 哈希对象
	h.Write([]byte(combined))                          // 将拼接后字符串的字节序列写入哈希对象
	return hex.EncodeToString(h.Sum(nil))              // 获取哈希值的十六进制字符串表示
}

// Login 验证登录信息
func Login(user *models.User) (err error) {

	// 先将原始密码拿出来
	opassword := user.Password // 用户输入的密码
	sqlStr := "select user_id , username , password from user where username = ?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}

	// 数据库查询失败
	if err != nil {
		return err
	}

	// 判断密码是否正确
	// 用户输入的加密后的密码和数据库中存储的加密后的密码是否一致
	password := EncryptPassword(opassword)
	if user.Password != password {
		return ErrorWrongPassword
	}
	return nil
}

// GetUserByID 根据ID获取用户信息
func GetUserByID(uid int64) (user *models.User, err error) {

	user = new(models.User)
	sqlStr := "select user_id , username , password from user where user_id = ?"
	err = db.Get(user, sqlStr, uid)
	if err == sql.ErrNoRows {
		return nil, ErrorUserNotExist
	}

	if err != nil {
		return nil, err
	}
	return user, err
}
