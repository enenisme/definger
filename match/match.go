package match

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/enenisme/definger/logger"
	"github.com/enenisme/definger/pkg"
	"github.com/enenisme/definger/utils"
)

// 添加包级别的正则表达式缓存
var (
	regexpCache sync.Map
	// 缓存大小限制 (例如限制2500个正则表达式)
	maxCacheSize = 2500
	// 当前缓存数量
	currentCacheSize atomic.Int32
)

// cacheItem 缓存项
type cacheItem struct {
	pattern string
	regexp  *regexp.Regexp
	lastUse int64
}

// cleanupCache 清理最近最少使用的缓存项
func cleanupCache() {
	if currentCacheSize.Load() < int32(maxCacheSize) {
		return
	}

	var items []cacheItem
	regexpCache.Range(func(key, value interface{}) bool {
		items = append(items, cacheItem{
			pattern: key.(string),
			regexp:  value.(*regexp.Regexp),
			lastUse: time.Now().UnixNano(),
		})
		return true
	})

	// 按最后使用时间排序
	sort.Slice(items, func(i, j int) bool {
		return items[i].lastUse < items[j].lastUse
	})

	// 删除最老的 20% 缓存项
	deleteCount := len(items) / 5
	for i := 0; i < deleteCount; i++ {
		regexpCache.Delete(items[i].pattern)
		currentCacheSize.Add(-1)
	}
}

// getOrCreateRegexp 获取或创建正则表达式
// 参数:
//   - pattern: 正则表达式字符串
//
// 返回值:
//   - *regexp.Regexp: 正则表达式对象
//   - error: 错误信息
func getOrCreateRegexp(pattern string) (*regexp.Regexp, error) {
	pattern = strings.ToLower(pattern)

	// 检查缓存
	if cached, ok := regexpCache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}

	// 编译新的正则表达式
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("编译正则表达式失败: %w", err)
	}

	cleanupCache()

	regexpCache.Store(pattern, compiled)
	currentCacheSize.Add(1)

	return compiled, nil
}

// Match 匹配探针结果和指纹
// 参数:
//   - httpResponse: 探针响应
//   - tags: 指纹
//
// 返回值:
//   - []string: 匹配到的指纹
//   - error: 错误信息
func Match(httpResponse *pkg.HttpResponse, tags *pkg.Tags, favicon string, logger *logger.Logger) ([]string, error) {
	if httpResponse == nil || tags == nil {
		return nil, fmt.Errorf("httpResponse或tags为空")
	}

	// 限制处理的响应体大小为10MB
	const maxBodySize = 10 * 1024 * 1024
	if len(httpResponse.Body) > maxBodySize {
		httpResponse.Body = httpResponse.Body[:maxBodySize]
	}

	headerStr := buildHeaderResponse(httpResponse)
	bodyStr := string(httpResponse.Body)
	matchedTags := make([]string, 0)

	logger.DebugResponsef("HTTP Response Header: %s", headerStr)
	logger.DebugResponsef("HTTP Response Body: %s", bodyStr)

	for _, tag := range tags.Tags {
		matches := 0
		totalMatchers := 0

		for _, http := range tag.HTTP {
			totalMatchers += len(http.Matchers)
			for _, matcher := range http.Matchers {
				switch {
				case matcher.Type == "word" && matcher.Part == "header":
					for _, word := range matcher.Words {
						if matchMode(headerStr, word, matcher.Condition) {
							if http.Mode == "or" || http.Mode == "" {
								matchedTags = append(matchedTags, tag.Info.Name)
							}
							matches++
						}
					}
				case matcher.Type == "word" && matcher.Part == "body":
					for _, word := range matcher.Words {
						if matchMode(bodyStr, word, matcher.Condition) {
							if http.Mode == "or" || http.Mode == "" {
								matchedTags = append(matchedTags, tag.Info.Name)
							}
							matches++
						}
					}
				}
			}
		}
		if matches >= totalMatchers { // 修改这里,要求matches必须等于totalMatchers才算匹配成功
			matchedTags = append(matchedTags, tag.Info.Name)
		}
	}

	return matchedTags, nil
}

// matchPattern 通用的正则匹配函数
// 参数:
//   - targetStr: 目标字符串
//   - patterns: 正则表达式列表
//
// 返回值:
//   - bool: 是否匹配到
func matchPattern(targetStr string, patterns []string) bool {
	for _, pattern := range patterns {
		if re, err := getOrCreateRegexp(pattern); err == nil && re.MatchString(targetStr) {
			return true
		}
	}
	return false
}

// matchMode 处理匹配模式
// 参数:
//   - content: 目标字符串
//   - word: 匹配字段
//   - mode: 匹配模式
//
// 返回值:
//   - bool: 是否匹配到
func matchMode(content, word, mode string) bool {
	// 处理OR模式
	if mode == "or" || mode == "" {
		return matchPattern(content, []string{"(?i)(" + word + ")"})
	}
	// 处理AND模式
	if mode == "and" {
		return matchPattern(content, []string{"(?i)(" + word + ")"})
	}
	return false
}

// buildHeaderResponse 构建HTTP响应头字符串
// 参数:
//   - httpResponse: 探针响应
//
// 返回值:
//   - string: 响应头字符串
func buildHeaderResponse(httpResponse *pkg.HttpResponse) string {
	var headerBuilder strings.Builder
	headerBuilder.Grow(len(httpResponse.Header) * 64) // 预分配内存

	for key, values := range httpResponse.Header {
		headerBuilder.WriteString(strings.ToLower(key))
		headerBuilder.WriteString(": ")
		headerBuilder.WriteString(strings.ToLower(strings.Join(values, ";")))
		headerBuilder.WriteString("\n")
	}
	return headerBuilder.String()
}

// MathTitle 获取title
// 参数:
//   - probes: 探针
//   - url: 目标URL
//
// 返回值:
//   - string: title
//   - error: 错误信息
func MathTitle(probes *pkg.Probes, url string) (string, error) {
	probe := utils.ProbesContent2ProbesStruct(utils.ProbesForGetTitle)
	resp, err := probes.HttpRequest(url, probe.Probes["093561eda8a835f5a01738826c77dbf6"])
	if err != nil {
		return "", fmt.Errorf("获取页面内容失败: %w", err)
	}

	titleRegex := regexp.MustCompile(`<title>(.*?)</title>`)
	matches := titleRegex.FindStringSubmatch(string(resp.Body))

	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("未匹配到title")
}

// MatchFavicon 获取favicon
// 参数:
//   - probes: 探针
//   - url: 目标URL
//
// 返回值:
//   - string: favicon_hash
//   - error: 错误信息
func MatchFavicon(probes *pkg.Probes, url string) (string, error) {
	probe := utils.ProbesContent2ProbesStruct(utils.ProbesForGetFavicon)
	resp, err := probes.HttpRequest(url, probe.Probes["favicon"])
	if err != nil {
		return "", fmt.Errorf("获取页面内容失败: %w", err)
	}
	return fmt.Sprintf("%x", md5.Sum(resp.Body)), nil
}
