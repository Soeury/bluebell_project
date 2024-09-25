package controllers

import (
	"project_bluebell/logic"
	"project_bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetPostDetailHandler
// @Summary 创建帖子并入库
// @Description 将请求中携带的帖子信息保存到数据库里面
// @Tags post_group
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param object body models.Post true "post : id,title,content,community"
// @Security ApiKeyAuth
// @Success 200 {object} ResNil_
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数 参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("invalid param", zap.Error(err))
		ResponseErr(c, CodeInvalidParam)
		return
	}

	// 获取当前用户的 ID 这个函数是通过 token 解析出来用户ID的，解析不出来说明token有问题，需要重新登录
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseErr(c, CodeNeedLogin)
		return
	}

	p.AuthorID = userID
	// 2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("failed create post", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler
// @Summary 获取指定帖子详细信息
// @Description 获取指定帖子列表的详细信息
// @Tags post_group
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "post_id"
// @Security ApiKeyAuth
// @Success 200 {object} ResPostList_
// @Router /post/:id [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取帖子ID
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("invalid post params", zap.Error(err))
		ResponseErr(c, CodeInvalidParam)
		return
	}

	// 2. 查询数据库
	detail, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, detail)
}

// GetPostListHandler
// @Summary 获取所有帖子详细信息
// @Description 获取所有帖子列表的详细信息
// @Tags post_group
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param page query string false "页码"
// @Param size query string false "每页数据量"
// @Security ApiKeyAuth
// @Success 200 {object} ResPostList_
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := GetPageAndSize(c)

	// 1. 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 根据前端需要(按时间 or 按分数 or 按社区)获取所有帖子列表
// @Tags post_group
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param object query models.ParamPostList false "query string"
// @Security ApiKeyAuth
// @Success 200 {object} ResPostList_
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// 1.获取参数
	// shouldbindquery 是从 queryString 里面获取参数
	// shouldbindjson  是前端传过来的json格式的数据
	// 注意定义的结构体的tag
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, // 默认值
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 failed", zap.Error(err))
		ResponseErr(c, CodeInvalidParam)
		return
	}

	// 像以后如果遇到了多个很相似的处理函数的接口，可以整合一个类似于中介的logic函数
	// 根据请求参数的不同来判断走哪个logic函数
	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("logic.GetPostListNew(p) failed", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, data)
}

/*

// GetCommunityPostListHandlers 查询某个社区下的所有post
func GetCommunityPostListHandler(c *gin.Context) {

	// 下面传入的是默认的queryString参数
	p := &models.ParamCommunityPostList{
		ParamPostList: &models.ParamPostList{
			Page:  1,
			Size:  10,
			Order: models.OrderTime,
		},
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandlers failed", zap.Error(err))
		ResponseErr(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityPostList2(p) failed", zap.Error(err))
		return
	}

	ResponseSuccess(c, data)
}

*/
