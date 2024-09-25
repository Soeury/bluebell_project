package mysql

import (
	"fmt"
	"project_bluebell/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

/*
  以下是建表语句 :
  这里为什么设置了 id 还要设置 user_id ?
    - 1. 注册账号的时候其他人可以通过 id 知道我们网站的一个大概的用户访问量
	- 2. 分库分表的时候出现重复
	- 3. uuid 的使用不方便

  由此，使用 分布式ID生成器 : 全局唯一，递增，正确的生成，确保在高并发的环境下也可以表现友好


create table `user` (
      `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
      `user_id` BIGINT(20) not null,
      `username` VARCHAR(64) COLLATE utf8mb4_general_ci NOT NULL,
      `password` VARCHAR(64) COLLATE utf8mb4_general_ci NOT NULL,
      `email` VARCHAR(64) COLLATE utf8mb4_general_ci,
      `gender` TINYINT NOT NULL DEFAULT(0),
      `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
      `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      PRIMARY KEY (`id`),
      UNIQUE KEY `idx_username` (`username`) USING BTREE,
      UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

*/

// 定义一个全局的变量，小写表示只在 dao 层才会使用到这个变量
var db *sqlx.DB

func Init(cfg *settings.MysqlConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("found err in sqlx.Connect", zap.Error(err))
		return
	}

	db.SetMaxIdleConns(cfg.Max_idle_conns)
	db.SetMaxOpenConns(cfg.Max_open_conns)
	return
}

func Close() {
	_ = db.Close()
}
