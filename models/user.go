package models

// 定义一个用户实例
// 注意这里的 UserID 是 int64 类型的
type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
