package finger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/enenisme/definger/logger"
	"github.com/enenisme/definger/match"
	"github.com/enenisme/definger/pkg"
	"github.com/enenisme/definger/utils"
)

type Finger struct {
	Url      string   // 目标URL
	Result   []string // 指纹结果
	Title    string   // 标题
	Protocol string   // 协议
	favicon  string   // favicon

	probes *pkg.Probes    // 探针配置
	tags   *pkg.Tags      // 指纹标签
	logger *logger.Logger // 日志对象

	async bool // 是否异步

	// 添加内存控制相关字段
	maxConcurrent   int   // 最大并发数
	maxResponseSize int64 // 最大响应大小
}

// NewFinger 创建Finger对象
// 参数:
//   - probes: 探针配置
//   - tags: 指纹标签
//   - logger: 日志对象
//
// 返回值:
//   - *Finger: 新创建的Finger实例
func NewFinger(probes *pkg.Probes, tags *pkg.Tags, logger *logger.Logger) *Finger {
	return &Finger{
		probes: probes,
		tags:   tags,
		logger: logger,
		Result: make([]string, 0),
	}
}

// Run 运行单个URL的指纹识别
// 参数:
//   - url: 目标URL
func (f *Finger) Run(url string) {
	if finger, err := f.finger(url); err != nil {
		f.logger.Warnf("指纹识别失败: %v", err)
	} else {
		f.logger.Success(finger.Result, finger.Url, finger.Title)
	}
}

// RunAsync 异步执行多URL指纹识别
// 参数:
//   - filePath: 包含URL列表的文件路径
func (f *Finger) RunAsync(filePath string) []Finger {
	fingers, err := f.fingerAsync(filePath)
	if err != nil {
		f.logger.Debugf("指纹识别失败: %v", err)
		return nil
	}
	return fingers
}

// finger 指纹识别核心函数
// 参数:
//   - url: 目标URL
//
// 返回值:
//   - *Finger: 包含识别结果的Finger对象
//   - error: 错误信息
func (f *Finger) finger(url string) (*Finger, error) {
	// 参数校验
	if err := f.validateParams(url); err != nil {
		return nil, fmt.Errorf("参数验证失败: %v", err)
	}

	f.Url = url
	if !f.async {
		f.logger.Infof("探针服务启动成功!")
		f.logger.Debugf("处理URL: %s", url)
	}

	// 发送HTTP请求并收集响应
	resps, err := f.sendProbeRequests(url)
	if err != nil {
		return nil, fmt.Errorf("探针请求失败: %v", err)
	}

	if !f.async {
		f.logger.Infof("指纹识别服务启动成功!")
	}

	// 获取favicon
	favicon, err := f.getFavicon()
	if err != nil {
		f.logger.Debugf("获取favicon失败: %v", err)
	} else {
		f.favicon = favicon
	}

	// 匹配指纹
	if err := f.matchFingerprints(resps); err != nil {
		return nil, fmt.Errorf("指纹匹配失败: %v", err)
	}

	// 获取标题
	if err := f.extractTitle(); err != nil {
		if !f.async {
			f.logger.Debugf("提取标题失败: %v", err)
		}
	}

	return f, nil
}

// validateParams 验证输入参数
// 参数:
//   - url: 目标URL
//
// 返回值:
//   - error: 错误信息
func (f *Finger) validateParams(url string) error {
	if url == "" {
		return fmt.Errorf("URL不能为空")
	}
	if f.probes == nil {
		return fmt.Errorf("探针配置不能为空")
	}
	if f.tags == nil {
		return fmt.Errorf("指纹标签不能为空")
	}
	return nil
}

// sendProbeRequests 并发发送探针请求
// 参数:
//   - url: 目标URL
//
// 返回值:
//   - []*pkg.HttpResponse: HTTP响应列表
//   - error: 错误信息
func (f *Finger) sendProbeRequests(url string) ([]*pkg.HttpResponse, error) {
	var wg sync.WaitGroup
	results := make(chan *pkg.HttpResponse, len(f.probes.Probes))
	errors := make(chan error, len(f.probes.Probes))

	// 并发发送请求
	for _, probe := range f.probes.Probes {
		wg.Add(1)

		// 跳过favicon探针
		if probe.Desc == "favicon" {
			continue
		}

		go func(p pkg.Probe) {
			defer wg.Done()
			resp, err := f.probes.HttpRequest(url, p)
			if err != nil {
				errors <- fmt.Errorf("探针请求失败: %v", err)
				return
			}
			results <- resp
		}(probe)
	}

	// 等待所有请求完成
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	return f.collectResponses(results, errors)
}

// collectResponses 收集HTTP响应
// 参数:
//   - results: 探针结果通道
//   - errors: 错误通道
//
// 返回值:
//   - []*pkg.HttpResponse: HTTP响应列表
//   - error: 错误信息
func (f *Finger) collectResponses(results chan *pkg.HttpResponse, errors chan error) ([]*pkg.HttpResponse, error) {
	var resps []*pkg.HttpResponse
	timeout := time.After(30 * time.Second)
	respCount := 0
	expectedCount := len(f.probes.Probes)

collectLoop:
	for respCount < expectedCount {
		select {
		case err, ok := <-errors:
			if !ok {
				continue
			}
			f.logger.Debugf(err.Error())
			respCount++
		case resp, ok := <-results:
			if !ok {
				break collectLoop
			}
			resps = append(resps, resp)
			respCount++
		case <-timeout:
			f.logger.Debugf("请求超时,开始处理已收到的响应")
			if len(resps) == 0 {
				return nil, fmt.Errorf("请求 %s 超时", f.Url)
			}
			break collectLoop
		}
	}

	if len(resps) == 0 {
		return nil, fmt.Errorf("未收到有效响应")
	}

	return resps, nil
}

// matchFingerprints 匹配指纹
// 参数:
//   - resps: 探针响应列表
//   - favicon: favicon
//
// 返回值:
//   - error: 错误信息
func (f *Finger) matchFingerprints(resps []*pkg.HttpResponse) error {
	var matchWg sync.WaitGroup
	matchResults := make(chan []string, len(resps))
	matchErrors := make(chan error, len(resps))

	// 并发匹配
	for _, resp := range resps {
		matchWg.Add(1)
		go func(r *pkg.HttpResponse) {
			defer matchWg.Done()
			matchedTags, err := match.Match(r, f.tags, f.favicon, f.logger)
			if err != nil {
				matchErrors <- fmt.Errorf("匹配失败: %v", err)
				return
			}
			if len(matchedTags) > 0 {
				matchResults <- matchedTags
			}
		}(resp)
	}

	// 等待所有匹配完成后再关闭通道
	go func() {
		matchWg.Wait()
		close(matchResults)
		close(matchErrors)
	}()

	return f.processMatchResults(matchResults, matchErrors)
}

// processMatchResults 处理匹配结果
// 参数:
//   - matchResults: 匹配结果通道
//   - matchErrors: 错误通道
//
// 返回值:
//   - error: 错误信息
func (f *Finger) processMatchResults(matchResults chan []string, matchErrors chan error) error {
	tagMap := make(map[string]struct{}, 32)
	errCount := 0

	// 使用WaitGroup等待所有结果处理完成
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(2)

	// 处理匹配结果
	go func() {
		defer wg.Done()
		for matchedTags := range matchResults {
			for _, tag := range matchedTags {
				mu.Lock()
				if _, exists := tagMap[tag]; !exists {
					tagMap[tag] = struct{}{}
					f.logger.Debugf("匹配到指纹: %s", tag)
					f.Result = append(f.Result, tag)
				}
				mu.Unlock()
			}
		}
	}()

	// 处理错误
	go func() {
		defer wg.Done()
		for err := range matchErrors {
			f.logger.Debugf(err.Error())
			errCount++
		}
	}()

	wg.Wait()

	if len(f.Result) == 0 {
		f.logger.Debugf("未匹配到指纹")
		return fmt.Errorf("未匹配到指纹")
	}

	if !f.async {
		f.logger.Infof("成功识别到指纹")
	}
	return nil
}

// extractTitle 提取标题
// 返回值:
//   - error: 错误信息
func (f *Finger) extractTitle() error {
	title, err := match.MathTitle(f.probes, f.Url)
	if err != nil {
		return fmt.Errorf("提取标题失败: %v", err)
	}
	f.Title = title
	return nil
}

// getFavicon 获取favicon
// 返回值:
//   - string: favicon
//   - error: 错误信息
func (f *Finger) getFavicon() (string, error) {
	favicon, err := match.MatchFavicon(f.probes, f.Url)
	if err != nil {
		return "", fmt.Errorf("获取favicon失败: %v", err)
	}
	return favicon, nil
}

// fingerAsync 异步处理多个URL的指纹识别
// 参数:
//   - filePath: 目标文件路径
func (f *Finger) fingerAsync(filePath string) ([]Finger, error) {
	var (
		mu      sync.Mutex
		results []Finger
	)

	urls, err := utils.LoadTargetFile(filePath)
	if err != nil {
		f.logger.Debugf("加载目标文件失败: %v", err)
		return nil, fmt.Errorf("加载目标文件失败: %v", err)
	}

	var wg sync.WaitGroup

	// 使用atomic包来保证计数器的原子性
	var successCount, failCount, timeoutCount uint32

	// 设置合理的并发数
	maxWorkers := 100
	if f.maxConcurrent > 0 {
		maxWorkers = f.maxConcurrent
	}

	// 使用带缓冲的通道控制并发
	semaphore := make(chan struct{}, maxWorkers)

	// 使用对象池减少内存分配
	fingerPool := sync.Pool{
		New: func() interface{} {
			return &Finger{
				probes:          f.probes, // 复用探针配置
				tags:            f.tags,   // 复用指纹标签
				logger:          f.logger, // 复用日志对象
				maxConcurrent:   100,
				maxResponseSize: 10 * 1024 * 1024,  // 10MB
				Result:          make([]string, 0), // 每次需要新的结果集
				async:           true,              // 异步模式
			}
		},
	}

	for i, url := range urls {
		if url == "" {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		if f.logger.Level >= logger.LogLevelVerbose {
			log.Printf("第 %d/%d 个URL: %s", i+1, len(urls), url)
			log.Printf("剩余 %d 个goroutine", runtime.NumGoroutine())
		}

		go func(u string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			finger := fingerPool.Get().(*Finger)
			finger.Result = finger.Result[:0]
			finger.Url = u
			finger.Title = ""

			if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
				u = "http://" + u
			}

			if f, err := finger.finger(u); err != nil {
				if strings.Contains(err.Error(), "超时") {
					atomic.AddUint32(&timeoutCount, 1)
				} else {
					finger.logger.Errorf("处理URL %s 失败: %v", u, err)
					atomic.AddUint32(&failCount, 1)
				}
			} else {
				finger.logger.Success(f.Result, f.Url, f.Title)
				atomic.AddUint32(&successCount, 1)
				mu.Lock()
				results = append(results, *f)
				mu.Unlock()
			}

			fingerPool.Put(finger)
		}(url)
	}

	// 等待所有goroutine完成后关闭通道
	wg.Wait()

	f.logger.Infof("指纹识别完成,成功: %d, 失败: %d, 超时: %d, 总数: %d", successCount, failCount, timeoutCount, len(urls))

	return results, nil
}
