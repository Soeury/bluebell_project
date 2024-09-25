package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	rate2 "github.com/juju/ratelimit"
	rate1 "go.uber.org/ratelimit"
)

// 限制流速的中间件 (令牌桶算法) , 函数参数表示:   fillinterval 秒放 1 个令牌 ， 令牌总容量是 cap
func RateLimitMiddleWare(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := rate2.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// TakeAvailable 表示获取1个令牌 ，返回桶中移除的令牌数量，没有可用数量的令牌就返回 0
		// TakeAvailable 不会阻塞(即不会等待令牌加入到桶中)
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		// 有令牌就继续
		c.Next()
	}
}

// 限制流速的中间件 (漏桶算法) : 源代码写的很好，可以多看看
func LeakyMiddleWare() func(c *gin.Context) {
	rl := rate1.New(100)

	//   r.Take() 返回下一滴水滴下来的时间点 , time.Until(xxx) 表示 xxx 时间减去现在的时间的结果
	//   int(time.Duration.Seconds()) 表示将一个 time.Duration 类型的数据转化成 int 类型
	//   注意 : 教程上第 0 秒的时刻有水 ， t.last 的第一个值是 0 ， 其他的稍微捋顺一下就好
	return func(c *gin.Context) {

		gapTime := time.Until(rl.Take())
		if int(gapTime.Seconds()) > 0 {
			// time.Sleep(gapTime) 需要睡眠的时长
			c.String(http.StatusOK, "rate limit ...")
			c.Abort()
			return
		}
		c.Next()
	}
}
