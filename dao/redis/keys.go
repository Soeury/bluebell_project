package redis

// 实现投票功能
const (
	KeyPrefix      = "bluebell:"
	KeyPostTime    = "post:time"   // zset ; 帖子及发帖时间
	KeyPostScore   = "post:score"  // zset ; 帖子及分数
	KeyPostVotedPF = "post:voted:" // zset ; 帖子及投票类型 这里是一个前缀，因为每个帖子都有自己的用户的投票情况
	KeyCommunityPF = "community:"  // set  ;
)

// GetRedisKey 给key加上前缀
func GetRedisKey(key string) string {
	return KeyPrefix + key
}
