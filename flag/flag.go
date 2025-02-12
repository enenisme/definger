package flag

import (
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/enenisme/definger/logger"
)

var (
	URL        string          // URL 指定要扫描的目标URL
	RuleFile   string          // RuleFile 指定规则文件的路径
	TargetFile string          // TargetFile 指定目标文件的路径
	LogLevel   logger.LogLevel // LogLevel 指定日志级别
	Timeout    int             // Timeout 指定超时时间(秒)
	OutputFile string          // OutputFile 指定输出文件的路径(excel)

	// util
	Json2Json   bool   // Json2Toml 是否将JSON文件转换为TOML文件
	OldJsonPath string // OldJsonPath 指定旧版指纹(JSON)文件的路径
	NewJsonPath string // NewJsonPath 指定新版指纹(JSON)文件的输出路径
)

// genericLogLevel 用于处理日志级别的设置和获取
type genericLogLevel struct {
	level logger.LogLevel
}

// Set 设置日志级别
func (g *genericLogLevel) Set(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	LogLevel = logger.LogLevel(v)
	g.level = LogLevel
	return nil
}

// String 获取日志级别的字符串表示
func (g *genericLogLevel) String() string {
	return strconv.Itoa(int(g.level))
}

// NewFlag 创建并返回一个新的命令行应用程序
func NewFlag() *cli.App {
	app := cli.NewApp()
	app.Name = "definger"  // 应用名称
	app.Usage = "新版指纹识别工具" // 应用描述
	app.Version = "1.0.1"  // 版本号

	// 定义命令行参数
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Value:   URL,
			Usage:   "指定扫描目标URL, 例如: -u http://example.com",
		},
		&cli.StringFlag{
			Name:        "ruleFile",
			Aliases:     []string{"r"},
			Value:       RuleFile,
			Usage:       "指定指纹规则文件路径",
			Destination: &RuleFile,
		},
		&cli.StringFlag{
			Name:        "targetFile",
			Aliases:     []string{"f"},
			Value:       TargetFile,
			Usage:       "指定目标文件路径",
			Destination: &TargetFile,
		},
		&cli.GenericFlag{
			Name:    "logLevel",
			Aliases: []string{"l"},
			Value: &genericLogLevel{
				level: 3,
			},
			Usage: "设置日志级别(1-5: ERROR, WARN, INFO, DEBUG, VERBOSE)",
		},
		&cli.IntFlag{
			Name:        "timeout",
			Aliases:     []string{"t"},
			Value:       10,
			Usage:       "设置请求超时时间(秒)",
			Destination: &Timeout,
		},
		&cli.StringFlag{
			Name:        "outputFile",
			Aliases:     []string{"o"},
			Value:       OutputFile,
			Usage:       "指定结果输出文件路径(excel)",
			Destination: &OutputFile,
		},
		&cli.BoolFlag{
			Name:        "jsonToJson",
			Aliases:     []string{"j2j"},
			Value:       Json2Json,
			Usage:       "是否将旧版指纹(JSON)文件转换为新版指纹(JSON)文件",
			Destination: &Json2Json,
			Action: func(ctx *cli.Context, b bool) error {
				Json2Json = b
				return nil
			},
		},
		&cli.StringFlag{
			Name:        "oldJsonFile",
			Aliases:     []string{"oj"},
			Value:       OldJsonPath,
			Usage:       "指定旧版指纹(JSON)文件路径",
			Destination: &OldJsonPath,
		},
		&cli.StringFlag{
			Name:        "newJsonFile",
			Aliases:     []string{"nj"},
			Value:       NewJsonPath,
			Usage:       "指定新版指纹(JSON)文件路径",
			Destination: &NewJsonPath,
		},
	}
	return app
}
