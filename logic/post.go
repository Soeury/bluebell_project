package logic

import (
	"project_bluebell/dao/mysql"
	"project_bluebell/dao/redis"
	"project_bluebell/models"

	"project_bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 将帖子数据插入到数据库中
func CreatePost(p *models.Post) (err error) {

	// 1. 生成 id
	p.ID = int64(snowflake.GetId())

	// 2. 数据入库
	if err := mysql.InsertPost(p); err != nil {
		zap.L().Error("mysql.InsertPost(p) failed", zap.Error(err))
		return err
	}

	// 3. 数据保存到redis
	if err := redis.InsertPost(p, p.CommunityID); err != nil {
		zap.L().Error("redis.InsertPost(p) failed", zap.Error(err))
		return err
	}
	return nil
}

// GetPostByID 根据ID查询某条post的详细详细详细数据(包括作者名字，帖子内容，社区内容)
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {

	// 1. 查询需要的数据
	//  -帖子详细数据
	post, err := mysql.GetPostDetailByID(pid)
	if err != nil {
		zap.L().Error(
			"mysql.GetPostDetailByID(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}

	//  -创建帖子的用户的详细数据
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error(
			"mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}

	//  -社区详细数据
	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error(
			"mysql.GetCommunityByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err),
		)
		return
	}

	// 2. 将我们想要返回的数据拼接在一起，
	// 这个 data 一定要先初始化
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	// 3. 返回数据
	return data, err
}

// GetPostList 获取帖子列表(包括作者名字，帖子内容，社区内容)
func GetPostList(page int64, size int64) (data []*models.ApiPostDetail, err error) {

	// 获取帖子信息列表
	posts := make([]*models.Post, 0, 200)              // make(类型，初始长度，容量)
	err = mysql.GetPostListInMysql(&posts, page, size) // 注意这里要把切片的指针传进去! ! !
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(posts)) // 有多少post就有多少list

	// for 循环把post对应的作者ID ，社区信息拿出来，然后拼接数据，最后追加到data里面
	for _, post := range posts {
		//  -创建帖子的用户的详细数据
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error(
				"mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		//  -社区详细数据
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error(
				"mysql.GetCommunityByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err),
			)
			continue
		}

		//  -将数据拼接到一起
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		//  -向data中追加拼接后的数据
		data = append(data, postDetail)
	}
	return data, err
}

// 按照(时间,分数)获取社区帖子列表
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 1.从redis中获取id列表
	ids, err := redis.GetPostInOrder(p)
	if err != nil {
		return nil, err
	}

	// 细节!
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostInOrder(p) return with 0 data")
		return nil, err
	}

	// 2.根据获取的id列表在mysql中获得post的详细信息 (post详细信息是存储在mysql中的)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 拿到每张帖子的赞成的票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 初始化data切片 , 有多少post就有多少list
	data = make([]*models.ApiPostDetail, 0, len(posts))
	// for循环拼接数据并拿出
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error(
				"mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("uid", post.AuthorID),
				zap.Error(err),
			)
		}

		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error(
				"mysql.GetCommunityByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err),
			)
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 从redis中获取community_id列表
	ids, err := redis.GetCommunityIDsInOrder(p, p.CommunityID)
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostInOrder(p) return with 0 data")
		return
	}

	// 2.根据获取的id列表在mysql中获得post的详细信息 (post详细信息是存储在mysql中的)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	// 拿到每张帖子的赞成的票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 初始化data切片 , 有多少post就有多少list
	data = make([]*models.ApiPostDetail, 0, len(posts))
	// for循环拼接数据并拿出
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error(
				"mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("uid", post.AuthorID),
				zap.Error(err),
			)
		}

		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error(
				"mysql.GetCommunityByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err),
			)
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// 将按照分数|时间返回帖子列表 和 按照社区返回帖子列表整合到一起
// p.communtiy == 0 表示返回全部
// p.community 有值表示按照社区返回
// 将两个查询logic合二为一的函数接口
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	if p.CommunityID == 0 {
		data, err = GetPostList2(p) // 查询并返回所有post
	} else {
		data, err = GetCommunityPostList(p) // 查指定社区中的post
	}

	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
