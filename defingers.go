package definger

import (
	"github.com/enenisme/definger/finger"
	"github.com/enenisme/definger/logger"
	"github.com/enenisme/definger/utils"
)

type Definger struct {
	URL string
}

func NewDefinger(url string) *Definger {
	return &Definger{URL: url}
}

func (d *Definger) Definger(path string) ([]string, error) {
	// 创建日志记录器
	logger := logger.NewLogger(logger.LogLevel(3))

	config, err := utils.LoadConfig(path)
	if err != nil {
		logger.Warnf("加载指纹规则文件失败: %v", err)
	}

	logger.Infof("加载探针服务配置成功！已识别探针数量: %d", len(config.Probes.Probes))
	logger.Infof("加载指纹服务配置成功！已识别指纹数量: %d", len(config.Tags.Tags))

	finger := finger.NewFinger(config.Probes, config.Tags, logger)
	finger.Run(d.URL)

	return finger.Result, nil
}
