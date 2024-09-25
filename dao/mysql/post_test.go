package mysql

import (
	"project_bluebell/models"
	"project_bluebell/settings"
	"testing"
)

// 这里容易出现一个空指针的问题，原因是测试仅执行以下测试函数中的内容
// 我们调用的InsertPost函数里面调用了 db , db 初始化的时候为一个 *sqlx.DB 类型的空指针
// 所以这里的测试需要先进行 db 初始化的一个操作

func init() {
	dbCfg := settings.MysqlConfig{
		Host:           "localhost",
		Port:           3306,
		User:           "root",
		Password:       "123456",
		Dbname:         "bluebell_test",
		Max_open_conns: 10,
		Max_idle_conns: 10,
	}

	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestInsertPost(t *testing.T) {

	p := models.Post{
		ID:          1,
		AuthorID:    int64(23456),
		CommunityID: int64(1),
		Title:       "title",
		Content:     "i am content",
	}

	err := InsertPost(&p)
	if err != nil {
		t.Fatalf("failed : %s\n", err)
	}
	t.Logf("success")
}

func TestGetPostDetailByID(t *testing.T) {

}

func TestGetPostListInMysql(t *testing.T) {

}

func TestGetPostListByIDs(t *testing.T) {

}
