package utils

import (
	"github.com/BurntSushi/toml"
	"github.com/enenisme/definger/pkg"
)

type ProbesConfig struct {
	Probes *pkg.Probes `toml:"probes"`
}

// ProbesContent2ProbesStruct 将probes内容转换为Probes结构体
func ProbesContent2ProbesStruct(content string) *pkg.Probes {
	var config ProbesConfig
	toml.Decode(content, &config.Probes)
	return config.Probes
}
