package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"project_bluebell/dao/mysql"
	"project_bluebell/settings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// controllers 层的测试主要是针对参数校验的测试
// dao 层的测试主要是针对数据插入或者查询的测试
// logic 层的测试主要是 模拟得到数据，并将数据进行处理的过程

// 注意 ， 测试中的代码会调用我们之前写过的函数或者方法
// 如果其中有没有初始化的参数或者变量(尤其是空指针)，需要提前初始化配置信息

// 测试代码不仅要测试成功的例子，可能出现的所有情况也需要测试 ! ! !

// 初始化 db 信息
func init() {
	dbCfg := settings.MysqlConfig{
		Host:           "localhost",
		Port:           3306,
		User:           "root",
		Password:       "123456",
		Dbname:         "bluebell_test",
		Max_open_conns: 10,
		Max_idle_conns: 10,
	}

	err := mysql.Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

// 这里会卡在 getCurrentUser 函数那里，因为这里是测试，没有登录，所以返回 codeNeedLogin
func TestCreatePostHandler(t *testing.T) {

	// 为了防止循环导包的问题，这里新建了一个 router，重新模拟路由注册的过程
	gin.SetMode(gin.TestMode) // 需要设置为 test 模式
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"community_id": 1 ,
		"title": "title" , 
		"content": "content" 
	}`

	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder() // 一个接收响应的对象
	r.ServeHTTP(w, req)         // 模拟处理请求 ， 将返回的响应放入 w 中

	// 这里判断状态码是否是 200
	assert.Equal(t, 200, w.Code)

	// 这里判断响应的内容是否是我们程序中定义好的
	assert.Contains(t, w.Body.String(), CodeMsgMap[CodeNeedLogin])
}

// 这里是查询指定 id 的 Post, 只要这个 post 存在，测试函数没问题，就可以查询到，所以这里返回 success
func TestGetPostDetailHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post/:id"
	r.GET(url, GetPostDetailHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/post/1001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), CodeMsgMap[CodeSuccess])
}

// 这里只是一个查询(第一条数据，只要数据存在,测试函数没问题，就可以查到，所以这里是返回 success
func TestGetPostListHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/posts"
	r.GET(url, GetPostListHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=1&size=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), CodeMsgMap[CodeSuccess])
}

// 这边用到的 redis 记得开，记得初始化配置信息。。。
func TestGetPostListHandler2(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/posts2"
	r.GET(url, GetPostListHandler2)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts2?page=1&size=1&community_id=0", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), CodeMsgMap[CodeServerBusy])
}
