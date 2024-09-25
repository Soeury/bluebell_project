package logic

import (
	"project_bluebell/dao/redis"
	"project_bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

// 投票的几种情况 :
// direction = 1
//   -1. 之前投反对票，现在改投赞成票    -> 更新分数和投票记录  差值 2
//   -2. 之前没有投票，现在改投赞成票    -> 更新分数和投票记录  差值的绝对值 1
// direction = 0
//   -1. 之前投反对票，现在改不投票      -> 更新分数和投票记录  差值的绝对值 1
//   -2. 之前投赞成票，现在改不投票      -> 更新分数和投票记录  差值的绝对值 1
// direction = -1
//   -1. 之前投赞成票，现在改投反对票    -> 更新分数和投票记录  差值的绝对值 2
//   -2. 之前没有投票，现在改投反对票    -> 更新分数和投票记录  差值的绝对值 1

//  投票限制 :
//    -1. post自发表之后一个星期之内允许用户投票，过期之后不允许投票
//    -2. 到期之后将每个post的赞成票数和反对票数存储到mysql中
//    -3. 到期之后删除 Redis 中的 KeyPostVotedPF

//  Bluebell 在这里使用简化版的投票方法 :
// 投票要考虑到帖子发布的时间(时间戳)，比如说分数相同的两个不火的帖子，后发的要比先发的排在前面
// 投一票就加上 432 分 ， 86400/200 = 432 ， 200张票就可以让帖子续一天

// VoteForPost 实现post投票功能
func VoteForPost(uid int64, p *models.ParamVoteData) error {

	// 记录一个 Debug 的日志
	zap.L().Debug(
		"vote for post",
		zap.Int64("uid", uid),
		zap.Int64("pid", p.PostID),
		zap.Int("direction", p.Direction),
	)

	uidStr := strconv.Itoa(int(uid))      // uid string
	pidStr := strconv.Itoa(int(p.PostID)) // pid string
	vote := float64(p.Direction)          // direction float64
	return redis.VoteForPost(uidStr, pidStr, vote)
}
