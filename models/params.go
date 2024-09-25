package models

// params 里面放的都是和请求参数相关的结构体

// tag 使用 json 是为了将数据从 json 格式的文件中取出来
// tag 使用 binding 为了进行参数的校验
//  * tag 的格式一定要正确，否则不能被识别!  格式:    `key1:"value1" key2:"value2"`

//  1. 定义请求的参数结构体
//  这里加了 bind:"required" 之后 ShouldBindJson 会校验该字段的值是否为空
//   eqfield=Password 表示这个字段的值必须与Password的值保持一致
//   json:"name,string" 表示后端将字符串类型的数据传递过去，前端将字符串类型的数据传递过来，但是得到的是我们需要的int类型
//   * 注意这里的大小写eqfield= 后面的值必须是结构体里面的字段，不能是我们自定义的

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`                    // 用户名
	Password   string `json:"password" binding:"required"`                    // 密码
	RePassword string `json:"repassword" binding:"required,eqfield=Password"` // 确认密码
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// ParamVoteData 投票请求参数
// binding 里面的 oneof=1 0 -1 表示这个字段的值只能是 1 0 -1 里面其中的一个
// oneof 里面的参数如果设置了0的话就不需要再binding里面写required字段了,不然传入0值会报错
type ParamVoteData struct {
	// UserID 可以从登陆的用户中获取即可
	PostID    int64 `json:"post_id,string" binding:"required"`       // 帖子id
	Direction int   `json:"direction,string" binding:"oneof=1 0 -1"` // 投票类型 : 赞成(1)反对(-1)取消(0)
}

// 注意这里结构体的tag需要根据前端传过来的数据决定使用哪个tag
type ParamPostList struct {
	Page        int64  `json:"page" form:"page" example:"1"`      // 分页页码
	Size        int64  `json:"size" form:"size" example:"10"`     // 每页数据量
	Order       string `json:"order" form:"order" example:"time"` // 排序方式(time | score)
	CommunityID int64  `json:"community_id" form:"community_id"`  // 可以为空
}

// // ParamCommunityPostList 某个社区下的所有post信息
// type ParamCommunityPostList struct {
// 	*ParamPostList
// }
