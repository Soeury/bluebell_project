package mysql

import (
	"database/sql"
	"project_bluebell/models"

	"go.uber.org/zap"
)

// GetCommunityData 实现在数据库中查询社区相关内容
func GetCommunityData() (List []*models.Community, err error) {

	// sqlx 里面单行查询是 get , 多行查询是 select ， 最好不要使用 * 进行返回数据
	sqlStr := "select community_id , community_name from community"
	if err := db.Select(&List, sqlStr); err != nil { // 注意这里要传递的是 List 参数的地址
		if err == sql.ErrNoRows {
			zap.L().Warn(WarnNoRows)
			err = nil
		}
	}
	return List, err
}

// 查询指定id的社区详细内容
func GetCommunityByID(community_id int64) (community *models.CommunityDetail, err error) {

	community = new(models.CommunityDetail)
	sqlStr := "select community_id , community_name , introduction , create_time from community where community_id = ?"
	if err = db.Get(community, sqlStr, community_id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn(WarnNoRows)
			err = ErrorInvalidID
		}
		return nil, err
	}
	return community, err
}
