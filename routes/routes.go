package routes

import (
	"net/http"
	"project_bluebell/controllers"
	"project_bluebell/logger"
	"project_bluebell/middlewares"
	"time"

	_ "project_bluebell/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetUp(mode string) *gin.Engine {

	// 如果 mode 设置成 release 那么就设置 gin 框架为 release 模式，其他默认为 debug 模式
	// gin 框架总共有三种模式， release   debug   test
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// 这里的 r := gin.New() 吗?
	r := gin.Default()
	//  下面这个路由加了 RateLimitMiddleWare 中间件 ， 令牌桶算法限制流速
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleWare(time.Second*1, 1))

	v1 := r.Group("/api/v1")

	// 测试
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 1. 用户注册
	v1.POST("/signup", controllers.SignUpHandler)
	// 2. 用户登录
	v1.POST("/login", controllers.LoginHandler)

	//  下面这个路由加了 JWTAuthMiddleware 中间件，能执行通过表示该用户是一个已经登录的用户
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		// *实现社区列表接口(返回所有的社区数据)
		v1.GET("/community", controllers.CommunityHandler)
		// *获取社区详情接口(返回指定 id 的社区数据)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		// *创建帖子接口(将用户传入的帖子数据存到数据库中)
		v1.POST("/post", controllers.CreatePostHandler)
		// *获取帖子详情接口(返回指定 id 的帖子)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		// *实现帖子列表接口(返回所有的 posts)
		v1.GET("/posts", controllers.GetPostListHandler)
		// *实现用户投票功能接口
		v1.POST("/vote", controllers.PostVoteHandler)
		// *实现按照分数或时间或社区获取帖子列表接口
		v1.GET("/posts2", controllers.GetPostListHandler2)
	}

	// 注册 swagger api相关路由 (路由处理都放在routes里面了)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 注册 pprof 相关路由(因为是第三方库，里面封装好了相关路由，使用 go 自带的 pprof 需要手动配置路由)
	pprof.Register(r)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "page not found!",
		})
	})

	return r
}
