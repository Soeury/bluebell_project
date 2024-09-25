package logic

import (
	"project_bluebell/dao/mysql"
	"project_bluebell/models"
)

// GetCommunityList 获得社区列表
func GetCommunityList() ([]*models.Community, error) {
	// 注意这里的数据类型是  []*name  不是  *[]name
	return mysql.GetCommunityData()
}

// GetCommunityDetail 获得某个社区的详细信息
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityByID(id)
}
