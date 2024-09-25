package controllers

import (
	"project_bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetPostDetailHandler
// @Summary 查询社区列表
// @Description 实现用户查询所有的社区的详细信息的功能
// @Tags community_group
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Security ApiKeyAuth
// @Success 200 {object} ResCommunity_
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	// 得到社区列表相关的数据(community_id , community_name) 以列表的形式返回
	list, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCoummnityList failed", zap.Error(err))
		ResponseErr(c, CodeServerBusy) // 服务器端的错误不轻易暴露出去，这里使用 CodeServerBusy
		return
	}

	ResponseSuccess(c, list)
}

// GetPostDetailHandler
// @Summary 查询指定社区信息
// @Description 实现用户查询指定的社区的详细信息的功能
// @Tags community_group
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "community_id"
// @Security ApiKeyAuth
// @Success 200 {object} ResCommunity_
// @Router /community/:id [get]
func CommunityDetailHandler(c *gin.Context) {
	// 获取社区 id   c.Param()是获取URL路径上的参数，返回的是string类型
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
	}

	// 获取指定id 的社区的详细信息
	detail, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, detail)
}
