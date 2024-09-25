package mysql

import (
	"project_bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// InsertPost 将帖子数据插入到数据库
func InsertPost(p *models.Post) (err error) {

	sqlStr := "insert into post(post_id , title , content , author_id , community_id) values(?,?,?,?,?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

// GetPostDetailByID 获取某条post的详细数据
func GetPostDetailByID(pid int64) (*models.Post, error) {

	post := new(models.Post)
	sqlStr := "select post_id , title , content , author_id , community_id , create_time from post where post_id = ?"
	err := db.Get(post, sqlStr, pid)
	return post, err
}

// GetPostListInMysql 获取post列表
func GetPostListInMysql(posts *[]*models.Post, page int64, size int64) error {

	// 这里返回的帖子按照时间降序排列的方式返回帖子列表
	sqlStr := "select post_id , title , content , author_id , community_id , create_time from post order by create_time desc limit ?,?"
	err := db.Select(posts, sqlStr, (page-1)*size, size)
	return err
}

// GetPostListByIDs 根据ID列表返回指定顺序的post信息
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {

	// FIND_IN_SET 是MySQL特有的函数，它在一个逗号分隔的字符串列表中查找一个字符串
	sqlStr := "select post_id , title , content , author_id , community_id , create_time from post where post_id in (?) order by FIND_IN_SET (post_id , ?)"
	// sqlx.In(sqlStr, ids)函数会根据输入生成一个新的SQL查询字符串(query)和一个参数切片(args)
	// 因为我们并不知道有多少参数，所以原始sql语句中只传一个?占位符，之后使用sqlx.In()来解析
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	// 重新绑定原 SQL 语句
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
