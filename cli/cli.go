package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/enenisme/definger/finger"
	"github.com/enenisme/definger/logger"
	"github.com/enenisme/definger/pkg"
	"github.com/enenisme/definger/utils"
)

// Args 定义命令行参数结构
type Args struct {
	URL        string // 目标URL
	RuleFile   string // 规则文件路径
	TargetFile string // 目标文件路径
	LogLevel   int    // 日志级别
	Timeout    int    // 超时时间
	OutputFile string // 输出文件路径(excel)

	Json2Json   bool   // 是否将旧版指纹(JSON)文件转换为新版指纹(JSON)文件
	OldJsonFile string // JSON文件
	NewJsonFile string // JSON文件
}

// NewArgs 创建新的Args对象
// 参数:
//   - c: CLI上下文
//
// 返回:
//   - *Args: 命令行参数对象
func NewArgs(c *cli.Context) *Args {
	return &Args{
		URL:        c.String("url"),
		RuleFile:   c.String("ruleFile"),
		TargetFile: c.String("targetFile"),
		LogLevel:   c.Int("logLevel"),
		Timeout:    c.Int("timeout"),
		OutputFile: c.String("outputFile"),

		Json2Json:   c.Bool("jsonToJson"),
		OldJsonFile: c.String("oldJsonFile"),
		NewJsonFile: c.String("newJsonFile"),
	}
}

// Run 运行主程序入口
// 参数:
//   - c: CLI上下文
//
// 返回:
//   - error: 错误信息
func Run(c *cli.Context) error {
	args := NewArgs(c)
	return args.run()
}

// run 执行指纹识别主逻辑
// 返回:
//   - error: 错误信息
func (a *Args) run() error {
	// 创建日志记录器
	logger := logger.NewLogger(logger.LogLevel(a.LogLevel))

	// 参数校验
	if err := a.validateArgs(logger); err != nil {
		return err
	}

	// 如果需要将旧指纹JSON文件转换为新指纹JSON文件，则执行转换
	if a.Json2Json {
		utils.Json2Json(a.OldJsonFile, a.NewJsonFile, logger)
	}

	// 如果目标文件不为空，则执行异步指纹识别
	if a.TargetFile != "" {
		logger.Infof("执行异步指纹识别")
		return a.runAsync(logger)
	}

	// 加载配置文件
	config, err := a.loadConfig(logger)
	if err != nil {
		return err
	}

	// 执行指纹识别
	return a.runFingerprint(logger, config)
}

func (a *Args) runAsync(logger *logger.Logger) error {
	// 加载配置文件
	config, err := a.loadConfig(logger)
	if err != nil {
		return err
	}

	// 执行指纹识别
	return a.runAsyncFingerprint(logger, config, a.TargetFile)
}

// validateArgs 验证必要的参数
// 参数:
//   - logger: 日志对象
//
// 返回:
//   - error: 错误信息
func (a *Args) validateArgs(logger *logger.Logger) error {
	if !a.Json2Json {
		if a.RuleFile == "" {
			logger.Warnf("指纹规则文件未指定")
			return fmt.Errorf("指纹规则文件未指定")
		}

		if a.TargetFile == "" && a.URL == "" {
			logger.Warnf("目标文件或URL未指定")
			return fmt.Errorf("目标文件或URL未指定")
		}
	} else {
		if a.OldJsonFile == "" {
			logger.Warnf("旧版指纹(JSON)文件未指定")
			return fmt.Errorf("旧版指纹(JSON)文件未指定")
		}

		if a.NewJsonFile == "" {
			logger.Warnf("新版指纹(JSON)文件未指定")
			return fmt.Errorf("新版指纹(JSON)文件未指定")
		}
	}

	return nil
}

// loadConfig 加载配置文件
// 参数:
//   - logger: 日志对象
//
// 返回:
//   - *pkg.Config: 配置对象
//   - error: 错误信息
func (a *Args) loadConfig(logger *logger.Logger) (*pkg.Config, error) {
	config, err := utils.LoadConfig(a.RuleFile)
	if err != nil {
		logger.Warnf("加载指纹规则文件失败: %v", err)
		return nil, err
	}

	logger.Infof("加载探针服务配置成功！已识别探针数量: %d", len(config.Probes.Probes))
	logger.Infof("加载指纹服务配置成功！已识别指纹数量: %d", len(config.Tags.Tags))

	return config, nil
}

// runFingerprint 执行指纹识别
// 参数:
//   - logger: 日志对象
//   - config: 配置对象
//
// 返回:
//   - error: 错误信息
func (a *Args) runFingerprint(logger *logger.Logger, config *pkg.Config) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		logger.Infof("总执行时间: %s", elapsed)
	}()

	finger := finger.NewFinger(config.Probes, config.Tags, logger)
	finger.Run(a.URL)
	return nil
}

func (a *Args) runAsyncFingerprint(logger *logger.Logger, config *pkg.Config, filePath string) error {
	startTime := time.Now()
	defer func() {
		elapsed := time.Since(startTime)
		logger.Infof("总执行时间: %s", elapsed)
	}()

	finger := finger.NewFinger(config.Probes, config.Tags, logger)
	fingers := finger.RunAsync(filePath)

	if a.OutputFile != "" {
		if strings.HasSuffix(a.OutputFile, ".xlsx") {
			fingerData := make(map[string]utils.FingerData)

			for _, finger := range fingers {
				fingerData[finger.Url] = utils.FingerData{
					Protocol: "TCP/HTTP",
					Url:      finger.Url,
					Result:   finger.Result,
					Title:    finger.Title,
				}
			}

			err := utils.SaveExecl(fingerData, a.OutputFile)
			if err != nil {
				logger.Warnf("保存Excel文件失败: %v", err)
				return err
			}
			logger.Infof("保存Excel文件成功: %s", a.OutputFile)
		}
	}

	return nil
}
