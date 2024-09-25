package models

import "time"

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 社区帖子模型
// 内存对齐: 尽量把结构体中字段类型相同的放在一起
type Post struct {
	ID          int64     `json:"id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id,string" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ApiPostDetail 获取post详细信息的接口模型
type ApiPostDetail struct {
	AuthorName       string `json:"author_name"`
	VoteNum          int64  `json:"vote_num"`
	*Post            `json:"post"`
	*CommunityDetail `json:"community_detail"`
}
