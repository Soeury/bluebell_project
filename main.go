package main

// @title Bluebell项目接口文档
// @version 1.0
// @description [use bluebell to practice]

// @contact.name Mr_rabbit
// @contact.email 1964475295@qq.com

// @host localhost:8080
// @BasePath /api/v1
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"project_bluebell/controllers"
	"project_bluebell/dao/mysql"
	"project_bluebell/dao/redis"
	"project_bluebell/logger"
	"project_bluebell/pkg/snowflake"
	"project_bluebell/routes"
	"project_bluebell/settings"
	"syscall"
	"time"

	"go.uber.org/zap"
)

//  web_staging 脚手架 - 通过结构体保存配置文件中的信息
//    当其他人接管了我们的项目时，可能不知道我们的配置文件中有哪些内容，
//    这时候通过结构体保存配置文件中的信息就显得非常友好

func main() {

	// 1. 加载配置
	if err := settings.ConfigInit(); err != nil {
		fmt.Printf("settings.ConfigInit failed: %v\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.LogInit(settings.Conf.LogConfig, settings.Conf.StagingConfig.Mode); err != nil {
		fmt.Printf("logger.Loginit failed: %v\n", err)
		return
	}

	defer zap.L().Sync() // 缓冲区的日志刷到磁盘中
	// 3. 初始化mysql连接
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("mysql.MysqlInit failed: %v\n", err)
		return
	}

	defer mysql.Close()

	// 4. 初始化redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("redis.RedisInit failed: %v\n", err)
		return
	}

	defer redis.Close()

	// 雪花算法生成一个唯一的 id
	if err := snowflake.Init(settings.Conf.StagingConfig); err != nil {
		fmt.Printf("snowflake.init failed: %v\n", err)
		return
	}

	// 初始化翻译器 , 翻译数据校验错误信息
	if err := controllers.InitTrans("en"); err != nil {
		fmt.Printf("controllers.initTrans failed: %v\n", err)
		return
	}

	// 5. 注册路由
	r := routes.SetUp(settings.Conf.StagingConfig.Mode)

	// 6. 优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.StagingConfig.Port),
		Handler: r,
	}

	// 开一个 goroutine 执行长时间的任务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Server Shutdown ...")
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.Shutdown(context); err != nil {
		zap.L().Fatal("srv.Shutdown failed : ", zap.Error(err))
	}

	zap.L().Info("server existing ...")
}
