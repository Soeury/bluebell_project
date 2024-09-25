package models

import "time"

type Community struct {
	ID   int    `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

// json 后面加上 omitempty 表示如果该字段为空，就不展示
type CommunityDetail struct {
	ID           int       `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
