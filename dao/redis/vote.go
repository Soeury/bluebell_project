package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

var (
	OneWeekInSeconds = int64(7 * 24 * 3600) // post限制投票的时间
	ErrOutofTime     = errors.New("out of limit time")
	ScorePerVote     = float64(432) // 每一票值的分数
	ErrVoteRepeat    = errors.New("repeat voting the same post is not allowed")
)

// VoteForPost 使用redis实现给用户投票的功能
func VoteForPost(uid string, pid string, vote float64) error {

	// 1. 判断投票是否超过该post的时间限制(自发表的一星期之内才能投票)
	// 这里使用Redis，因为是kv存储，所以可以通过 postID 拿到发帖时间，分数，用户的投票情况
	postTime := Client.ZScore(GetRedisKey(KeyPostTime), pid).Val()
	if time.Now().Unix()-int64(postTime) > OneWeekInSeconds {
		return ErrOutofTime
	}

	// 2. 更新post分数 (来到这里说明投票时间没有超过post的时间限制)
	// 更新的分数依赖于用户之前的投票记录,这里要先获取用户之前的投票记录
	ovote := Client.ZScore(GetRedisKey(KeyPostVotedPF+pid), uid).Val()
	// 如果这一次投票和之前的一致，就返回错误(不允许重复投票)
	if ovote == vote {
		return ErrVoteRepeat
	}
	var direct float64 // 设置方向
	if vote > ovote {
		direct = 1
	} else {
		direct = -1
	}

	// 两次投票的差值
	diff := math.Abs(ovote - vote)
	// 更新分数 : 下面两个操作需要同时成功或者失败，开启 redis 事务
	pipeline := Client.TxPipeline()
	pipeline.ZIncrBy(GetRedisKey(KeyPostScore), diff*direct*ScorePerVote, pid)

	// 3. 记录用户为该post的投票类型(支持票，反对票，不投票)
	// 这里好像突然顿悟了! ! ! redis有很多库
	// 每个库里面好像是每一个键值存储组合都有一个类似于title的记号，通过title可以操作里面的kv对
	if vote == 0 {
		pipeline.ZRem(GetRedisKey(KeyPostVotedPF+pid), pid) // 这里直接移除用户记录就可以
	} else {
		pipeline.ZAdd(GetRedisKey(KeyPostVotedPF+pid), redis.Z{
			Score:  vote, // 用户投票情况 赞成 或者 反对
			Member: uid,  // 用户ID
		})
	}

	_, err := pipeline.Exec()
	return err
}
