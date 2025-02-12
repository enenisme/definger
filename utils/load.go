package utils

import (
	"encoding/json"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/enenisme/definger/pkg"
)

// LoadConfig 一次性加载所有配置
// 参数:
//   - filePath: 配置文件路径
//
// 返回:
//   - *Config: 配置结构
//   - error: 错误信息
func LoadConfig(filePath string) (*pkg.Config, error) {
	var config pkg.Config

	// 加载指纹配置
	tagsData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tags []pkg.Tag
	if err := json.Unmarshal(tagsData, &tags); err != nil {
		return nil, err
	}

	config = pkg.Config{
		Tags:   &pkg.Tags{Tags: tags},
		Probes: ProbesContent2ProbesStruct(ProbesContent),
	}

	return &config, nil
}

func LoadTargetFile(filePath string) ([]string, error) {
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var urls []string
	switch runtime.GOOS {
	case "linux", "darwin":
		urls = strings.Split(string(content), "\n")
	case "windows":
		urls = strings.Split(string(content), "\r\n")
	}
	return urls, nil
}
