package controllers

import (
	"project_bluebell/logic"
	"project_bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GetPostDetailHandler
// @Summary 实现帖子投票功能
// @Description 实现用户对帖子进行投票的功能
// @Tags vote_group
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param object body models.ParamVoteData true "vote : post_id direction"
// @Security ApiKeyAuth
// @Success 200 {object} ResNil_
// @Router /vote [post]
func PostVoteHandler(c *gin.Context) {
	// 请求参数获取 校验
	// 注意这里前端传过来的两个数据 post_id , direction 都是 string 类型的数据(json里面指定了)
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Post vote with invalid params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 类型断言，是否触发validator里面的错误
		if !ok {
			ResponseErr(c, CodeInvalidParam)
			return
		}
		errsData := RemoveTopStruct(errs.Translate(trans))
		ResponseErrWithMsg(c, CodeInvalidParam, errsData)
		return
	}

	// 获取当前用户ID
	uid, err := GetCurrentUser(c)
	if err != nil {
		ResponseErr(c, CodeNeedLogin)
		return
	}

	// ***实现投票功能***
	if err := logic.VoteForPost(uid, p); err != nil {
		zap.L().Error("logic.VoteForPost(uid, p) failed", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, nil)
}
