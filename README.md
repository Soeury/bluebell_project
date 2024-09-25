# bluebell_project
一个练习gin的小作业

## 一些小细节
## Details :

1. 没思路可以写注释，将要做的事情的注释写下来，再翻译成代码。像这里最开始都是从路由写起。要实现一个功能从路由入手

2. route -> controllers (三层) -> logic -> 最底层的实现

3.  登录/注册/... 等需要取出值的功能需要进行参数校验，参数校验先取出json/xml/...的值到struct
    校验分成成功和失败，成功打印输出信息，失败需要将错误信息进行翻译

4.  存储进数据库的密码需要进行加密  
    ****这里很重要啊****
    老师举的例子是使用md5哈希算法实现的，包: "crypto/md5"
    然而，还有其他的方法可以实现  在Go语言中，实现密码加密通常会使用哈希函数  包: "crypto/sha256" 
    盐值（Salt）在密码加密中是一个重要的概念，它主要用于增加密码存储的安全性
        //  md5 加密算法实现密码加密 
        //  MD5是一种单向哈希函数，它可以将任意长度的输入转换成一个固定长度
        //  对于MD5来说是128位，即16字节，通常以32个十六进制字符表示的哈希值
        //  这个转换过程是单向的，无法从已经加密后的密码解析出原来的密码
        //
        //
        // sha256 算法实现密码加密
        // 随机生成一个指定长度的盐值字符串
        // 如果每个密码都有其对应的盐值的话，存储时，通常需要同时存储盐值和哈希值
        // 在验证密码时，需要使用相同的盐值来重新计算哈希值


5.  返回响应的处理重复性太高，可以在处理返回响应的那一层封装自己的 code 和 response ,需要的时候调用即可
        *这里自定义 code 码和 Msg 信息是注意以 Codexxxxxx 定义会更直观更方便，


6.  返回的错误是由最底层定义的错误一直往上传的


7.  最好不要在程序中出现莫名定义的字符串，可以把这些数据定义成全局变量/常量放在最上面，调用即可
        *定义全局的错误类型字符串最好以  Errorxxxxxx 定义会更直接更方便


8.  数据库中ID的问题，最好是通过某些算法(etc: snowflake)生成一些可用性高的ID


9.  gin框架支持多种模式， debug  release  test


10. 客户端携带 Token 的方式有三种， -1.请求头  -2.请求体  -3.URL
    bluebell假设token是放在header的Authorization中，并且使用bearer开头
	具体位置需要根据实际情况来写


11.  JWTAuthMiddleware 中间件，能执行通过表示该用户是一个已经登录的用户
    之后，在 Bluebell 中需要登录才能访问的地方加上这个中间件就可以了


12.  注意 import cycle not allowed 循环导包的问题，主要原因是对 Bluebell 中的模块划分不熟悉


13. access Token 和 refresh Token  : 双token实现无感刷新   
    access Token 是我们访问网站资源时使用的 Token , 还有一种 refresh Token 
    通常 refresh Token 的有效时间会比较长，access Token 的有效时间会比较短
    当 access Token 由于过期而失效时，可以通过 refresh Token 生成一个新的 access Token 
    如果 refresh Token 也失效了，用户只能重新登录
    后端需要定义一个刷新 Token 的接口
    前端需要实现一个拦截器，当 access Token 过期的时候自动请求刷新 Token 的接口来获取新的 access Token


14. token - 限制同一个账号同一段时间只能登录一台设备
    多台设备在一个时间段内登录同一个账号生成的 Token 是不一样的，通过这里可以解决限制的功能
    比如说: 在后端数据库中建立一个数据库，将 token 与 user_id 一一对应存储进去
    服务器在检验 token 没有失效之后，将 token 中携带的自定义的数据(user_id)解析出来(在定义的解析JET的中间件里面)
    通过与数据库中存储的 user_id 对应的 token 进行比较，如果 token 不一致，说明有多个用户同时登录了
    检验到 token 不一致之后，返回让用户重新登录，重新登录之后会生成一个新的 token ，
    这个时候，携带旧的token的用户再次访问时，需要重新登录


15. 借助 Make 我们在编译过程中不再需要每次手动输入编译的命令和编译的参数，可以极大简化项目编译过程
    make 是一个构建自动化的工具，会在当前目录下寻找 MakeFile 或者 MakeFile 文件，
    如果存在相应的文件，它就会依据其中定义好的规则完成构建任务
    Makefile 里面顶格的四个空，只能用 tab 
    (shell 脚本文件 .sh 和 Makefile 文件的区别) ! ! !


16. Air 工具实现实时监听项目的代码文件，在代码发生变化之后自动编译并执行，大大提高gin框架项目的开发效率
    安装方式:   go get -u github.com/cosmtrek/air
    (选择性使用) 即可以不使用


17. 书写结构体时，尽量把数据类型相同的字段放在一起，(内存对齐)


18. 如果后端需要返回多个结构体里面的多个内容，比如 Bluebell 通过帖子ID返回作者姓名，帖子详情，社区详情
    这时候可以定义一个 API 的一个总的结构体，把需要的数据字段或者指针填充进去
    然后将需要的数据查询出来之后，将数据拼接到一起，最后返回即可


19. 注意，Bluebell 中，获取社区或者帖子列表时，使用 db.Select()查询数据库，需要传递的是切片的指针，而不是切片本身


20. 解决后端传给前端的 ID 值失真的问题 : 
    当后端传递给前端的某个int类型的数据过大的时候(即超过了前端数据能够接收的范围)
    这时候，传递过去的数据就会被改变(失真)
    解决办法: 将后端的数据改成结构体类型传递给前端，前端将数据返回过来时改成我们需要的int类型即可
    结构体tag后面加上,string     例:   type Params struct { ID  int64  `json:"id,string"` }


21. 定义路由组的时候可以加上 : v1 := r.Group("/api/v1") 表示我们自己写的v1版本的API接口


22. 记得看一下 Validator 文档，了解参数校验的 binding 后面还可以添加哪些功能， ->  方便
    required   oneof   eqfield   ...


23. 注意一下老师的 go 的一些 json 的技巧


24. shouldbindquery 是从 queryString里面获取参数
    shouldbindjson  是前端传过来的json格式的数据
    shouldbind      是根据前端传过来的数据进行自动的推导
    注意定义的结构体的tag :   form  json  mapstructure  binding 


25. string := c.Query("page")   //     /api/v1/post?page=1&size=1
    string :=c.Param("id")      //     /api/v1/:id


26. 以下是数据库相关内容 :  
    sqlStr := "select post_id , title from post where id in (?) order by FIND_IN_SET (post_id , ?)"
    query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
    query = db.Rebind(query)
    err = db.Select(&postList, query, args...)
    //   FIND_IN_SET 是MySQL特有的函数，它在一个逗号分隔的字符串列表中查找一个字符串
    //   sqlx.In(sqlStr, ids)函数会根据输入生成一个新的SQL查询字符串(query)和一个参数切片(args)
    //   因为我们并不知道有多少参数，所以原始sql语句中只传一个?占位符，之后使用sqlx.In()来解析
    //   sql.In解析之后需要重新绑定，然后再次查询(这时候正常查询就可以了)
    //   再次查询时，args...不要忘记了


27. 使用 pipeline 一次发送多个请求，减少RTT :
    通过 PIPELINE 客户端可以将多个命令打包成一次发送，而不是对每个命令都发送一个请求并等待其响应
    这种方式显著减少了网络往返时间(RTT, Round-Trip Time)


28. *使用 swagger 自动生成接口文档 :
    - 按照要求给代码添加注释  
        1. main中主注释   
        2. controllers中接口注释，每个注册过的路由都需要写
    - 终端输入 swag init , 如果写的注释没有问题的话，会在根目录下自动生成一个docs文件夹
    - main 中需要导入一些包 :
        1. docs 包匿名导入
        2. swaggerFiles "github.com/swaggo/files"
        3. ginSwagger "github.com/swaggo/gin-swagger"
    - 注意 : 写注释的时候都是有格式的，写的时候复制粘贴再改一下就可以了
    - 注意 : 写controllers里面的注释的时候，注意请求中传来的参数的type , 
            Bluebell 中使用到的param_type有 
                1. query 表示url中?后面的kv对
                2. body  这个应该是请求体中的吧...
                3. path  表示url路径中的参数
    - 之后打开浏览器访问 http://localhost:8080/swagger/index.html 即可
    - 实在忘记了的话去看文档吧...


29. go 单元测试(很重要啊) :
    go test ./controllers -v 查看所有以_test.go结尾的测试文件

        安装 : go get github.com/stretchr/testify
        导入包 : "github.com/stretchr/testify/assert"

    去了解了一下 : go 测试应该分成 
        1.  单元测试:  容易发现 bug， 测试文件以_test.go 结尾，测试函数以 Test... 开头传入参数为 *testing.T 类型
                      由于我们时候测试的函数比较复杂，需要进行一些测试前的数据，配置等的初始化操作
        2.  mock测试: 也叫做打桩，作用是降低程序不同模块的耦合度，是的出现错误的来源更单一
        3.  基准测试: 分析测试对象的性能，进行对比分析


30. VScode 里面   go get github.com/.../...   之后如果遇到   bash : xxx command not found   的问题
    注意，这里的 go get 只是导入外部库，要想命令可以使用必须
    要在   go get github.com/.../...   之后再   go install github.com/.../...  最后  go mod tidy


31.  HTTP 压力测试常用工具 :  go-wrk  详细见文档(github.com  go-wrk) 


32. 限流又称为流量控制（流控），通常是指限制到达系统的并发请求数 ， 这里使用用 漏桶 和 令牌桶 来达到限流的效果
    - 漏桶: 无论多少请求，都按照固定的速率处理，缺点 : 不好处理大量的突发请求的情况
    - 令牌桶: 令牌桶按固定的速率往桶里放入令牌，只要能从桶里取出令牌就能通过，有点 : 支持突发流量的快速处理
    (感觉令牌桶会使用的更广泛一点)


33. pprof 性能调优:  pprofiling 是指对应用程序运行过程中 内存 和 cpu 的使用情况的分析
    常配合压测来使用 ，通常只在性能测试的时候才在代码中引入pprof
    go 性能优化的几个方面 :   cpu   memory  goroutine
    gin 框架可以使用第三方的 pprof 库 "github.com/gin-contrib/pprof"
        
    步骤:   1.   pprof.Register(r)  
            2.   配合压测在终端使用  go tool pprof http://127.0.0.1/debug/pprof/xxxxx
            3.   其中，xxxxx 表示 pprof 中支持的想要测试的指标,可以在 127.0.0.1:8080/debug/pprof 查看
            4.   进入 pprof 模式，可以查看程序的运行情况  
    数据可视化 :  go-torch 生成火焰图(安装有点麻烦，还是不实现了)

        
34.  使用  nohup , nginx , supervisor 进行项目部署(后台守护进程)


35. 企业代码发布流程:   CICD : 持续集成，持续交付和持续部署，使用相关工具实现自动化构建，测试，部署和监控软件
    1. 上传   ->   2. 构建，测试     ->     3. 代码评审，配置环境，预发布，灰度发布，上线

