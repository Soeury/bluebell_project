package snowflake

import (
	"project_bluebell/settings"
	"time"

	sf "github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *sf.Node

func Init(cfg *settings.StagingConfig) (err error) {

	var st time.Time
	st, err = time.Parse("2006-01-02", cfg.StartTime)
	if err != nil {
		zap.L().Error("[time.Parse failed ...]", zap.Error(err))
		return
	}

	// 解析到的时间转换成一个毫秒级的时间戳
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(cfg.MachineId) // 通过machineId 的值创建一个新的 node实例
	zap.L().Error("[sf.NewNode failed ...]", zap.Error(err))
	return
}

// 上面的函数执行成功后，后续要拿到 id 只需要调用下面这个函数就可以了
func GetId() int64 {
	return node.Generate().Int64()
}
