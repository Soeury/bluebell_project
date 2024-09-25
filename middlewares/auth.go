package middlewares

import (
	"project_bluebell/controllers"
	"project_bluebell/pkg/jwt"

	"strings"

	"github.com/gin-gonic/gin"
)

// 认证 JWT 的一个中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带 Token 的方式有三种， -1.请求头  -2.请求体  -3.URL
		// 这里假设token是放在 header 的 Authorization 中，并且使用 bearer 开头
		// 具体位置需要根据实际情况来写

		// 1. 是否存在
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controllers.ResponseErr(c, controllers.CodeNeedLogin)
			c.Abort()
			return
		}

		// authorization 格式是否正确
		// 下面表示将 autuHeader 字符串按照 " " 分成最多两个字符串并且保存在切片中
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseErr(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}

		// parts[1] 是获取到的 tokenString , 解析JWT
		// 返回的 mc 是我们自定义的结构体，里面存储了我们想要拿到的数据(uid , name)
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseErr(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}

		// 将当前请求的 username 信息保存到请求的上下文中
		// 方便后续通过 Get 获取当前请求的用户信息
		// 这里我把 controllers.CtxUserIDKey 放到 controllers 里面去了，(为了解决循环导包的问题)
		c.Set(controllers.CtxUserIDKey, mc.User_ID)
		c.Next()
	}
}
