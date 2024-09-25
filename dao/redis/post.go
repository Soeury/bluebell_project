package redis

import (
	"project_bluebell/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// getIDs 公用的部分另外封装一个函数
func getIDs(key string, page int64, size int64) ([]string, error) {

	// 确定起始索引和结束索引
	start := (page - 1) * size
	end := start + size - 1 // 因为是按照索引取出来，所以这里要 - 1

	// 返回数据  ZRevRange 表示按照值从高到低排序查询指定数量的post
	return Client.ZRevRange(key, start, end).Result()
}

// * * * * * * * * * * *! ! ! ! ! ! !
// InsertPost 将帖子的创建时间存入库中,后面需要根据创建时间得到分数
func InsertPost(p *models.Post, cid int64) error {

	// 下面两个操作需要同时成功，开启 redis 事务 :
	pipeline := Client.TxPipeline()
	// 帖子创建的时间当作当前的分数
	pipeline.ZAdd(GetRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.ID,
	})

	// 实际分数
	pipeline.ZAdd(GetRedisKey(KeyPostScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.ID,
	})

	// 将社区ID加入到社区的set
	ckey := GetRedisKey(KeyCommunityPF + strconv.Itoa(int(cid)))
	pipeline.SAdd(ckey, p.ID)
	_, err := pipeline.Exec()
	return err
}

// GetPostInOrder 按指定顺序获取postID
func GetPostInOrder(p *models.ParamPostList) ([]string, error) {

	var key string
	// 这里确定需要排序的类型 (按时间|按分数)
	if p.Order == models.OrderScore {
		key = GetRedisKey(KeyPostScore)
	} else {
		key = GetRedisKey(KeyPostTime)
	}

	return getIDs(key, p.Page, p.Size)
}

// 获取某个帖子投赞成票的数量
func GetPostVoteData(ids []string) (data []int64, err error) {

	/*  这里会可能会出现查询次数特别多的情况
	data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := GetRedisKey(KeyPostVotedPF + id) // 拿到key
		// 查找key中元素值是"1"的元素个数
		v := Client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}
	return
	*/

	// 使用pipeline一次发送多个请求，减少RTT
	pipeline := Client.Pipeline()
	for _, id := range ids {
		key := GetRedisKey(KeyPostVotedPF + id)
		pipeline.ZCount(key, "1", "1")
	}

	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityIDsInOrder 按社区查询ids(这里也涉及到时间和分数)
// 即按照时间查询某个社区的所有帖子,还是按照分数查询某个社区的所有帖子
// ? ? ? ? 这里拿不到数据 ? ? ?
// (os:社区key好像没有存数据进去...)
func GetCommunityIDsInOrder(p *models.ParamPostList, cid int64) ([]string, error) {

	// 使用 zinterstore 将分区的set集合和帖子分数的zset的集合取一个交集的zset集合
	// 然后按照之前的逻辑拿出 ids []string
	// 引入缓存key减少zinterstore的执行次数
	var orderKey string
	if p.Order == models.OrderTime {
		orderKey = GetRedisKey(KeyPostTime)
	} else {
		orderKey = GetRedisKey(KeyPostScore)
	}

	ckey := GetRedisKey(KeyCommunityPF + strconv.Itoa(int(cid))) // 社区Key
	key := orderKey + strconv.Itoa(int(cid))                     // 缓存Key
	if Client.Exists(orderKey).Val() < 1 {
		// 不存在，需要计算
		pipeline := Client.Pipeline()
		// 这个操作会将 ckey 和 orderKey 进行交集运算，取最大值(Aggregate: "MAX")，并将结果存储在之前构建的key中
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, orderKey)
		// 60秒后，这个缓存键及其对应的值（即计算结果）将自动从Redis数据库中删除
		pipeline.Expire(key, 60*time.Second)
		// 执行管道中的所有命令
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDs(key, p.Page, p.Size)
}
