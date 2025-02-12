package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/enenisme/definger/logger"
	"github.com/enenisme/definger/pkg"
)

// JsonData 定义了JSON数据的结构
type JsonData struct {
	ID   string `json:"id"`   // 规则ID
	Name string `json:"name"` // 规则名称
	Type string `json:"type"` // 规则类型
	Mode string `json:"mode"` // 规则模式
	Http Http   `json:"http"` // HTTP请求相关配置
	Rule Rule   `json:"rule"` // 规则匹配条件
}

// Http 定义了HTTP请求的结构
type Http struct {
	ReqMethod string            `json:"reqMethod"` // HTTP请求方法
	ReqPath   string            `json:"reqPath"`   // 请求路径
	ReqHeader map[string]string `json:"reqHeader"` // 请求头
	ReqBody   string            `json:"reqBody"`   // 请求体
}

// Rule 定义了规则匹配条件的结构
type Rule struct {
	InBody   string `json:"inBody"`   // 响应体中的匹配规则
	InHeader string `json:"inHeader"` // 响应头中的匹配规则
	InIcoMd5 string `json:"inIcoMd5"` // favicon.ico的MD5匹配规则
}

// NewInfo 创建新的info字段
func NewInfo(name string) map[string]interface{} {
	return map[string]interface{}{
		"name":     name,
		"author":   "",
		"tags":     "",
		"severity": "info",
		"metadata": map[string]interface{}{
			"product":  "",
			"vendor":   "",
			"verified": true,
		},
	}
}

// NewHttpMatcher 创建新的HTTP匹配器
// 参数:
//   - method: HTTP请求方法
//   - path: HTTP请求路径
//
// 返回值:
//   - map[string]interface{}: 新的HTTP匹配器
func NewHttpMatcher(method, path, mode string) map[string]interface{} {
	if mode == "" {
		mode = "or"
	}
	return map[string]interface{}{
		"method":   method,
		"path":     []string{path},
		"mode":     mode,
		"matchers": []interface{}{},
	}
}

// readJsonFile 从文件中读取JSON数据并解析为JsonData结构
// 参数:
//   - filepath: 文件路径
//
// 返回值:
//   - []JsonData: JsonData结构数组
//   - error: 可能的错误
func readJsonFile(filepath string) ([]JsonData, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var jsonData []JsonData
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

// json2Json 将JsonData结构转换为新的JSON格式
// 参数:
//   - jsonData: JsonData结构数组
//
// 返回值:
//   - map[string]interface{}: 新的JSON格式
func json2Json(jsonData []JsonData) []map[string]interface{} {
	var result []map[string]interface{}

	for _, data := range jsonData {
		newJson := make(map[string]interface{})
		newJson["id"] = data.Name
		newJson["info"] = NewInfo(data.Name)

		httpMatcher := NewHttpMatcher(data.Http.ReqMethod, data.Http.ReqPath, data.Mode)
		matchers := httpMatcher["matchers"].([]interface{})

		// 处理匹配规则
		if data.Rule.InHeader != "" {
			headerArray, mode := parseMatchRule(data.Rule.InHeader)
			matchers = append(matchers, createMatcher(headerArray, "header", mode))
		}
		if data.Rule.InBody != "" {
			bodyArray, mode := parseMatchRule(data.Rule.InBody)
			matchers = append(matchers, createMatcher(bodyArray, "body", mode))
		}
		if data.Rule.InIcoMd5 != "" {
			matchers = append(matchers, createMatcher([]string{data.Rule.InIcoMd5}, "favicon", ""))
		}

		httpMatcher["matchers"] = matchers
		newJson["http"] = []interface{}{httpMatcher}

		result = append(result, newJson)
	}
	return result
}

// parseMatchRule 解析匹配规则字符串,返回规则数组和匹配模式
// 参数:
//   - content: 匹配规则字符串
//
// 返回值:
//   - []string: 规则数组
//   - string: 匹配模式
func parseMatchRule(content string) ([]string, string) {
	switch {
	case strings.Contains(content, "|"):
		return strings.Split(strings.ReplaceAll(strings.ReplaceAll(content, "(", ""), ")", ""), "|"), "or"
	case strings.Contains(content, "&&"):
		return strings.Split(strings.ReplaceAll(strings.ReplaceAll(content, "(", ""), ")", ""), "&&"), "and"
	default:
		return []string{content}, ""
	}
}

// createMatcher 根据给定的值、类型和模式创建匹配器
// 参数:
//   - values: 规则数组
//   - matchType: 匹配类型
//   - mode: 匹配模式
//
// 返回值:
//   - map[string]interface{}: 匹配器
func createMatcher(values []string, matchType, mode string) map[string]interface{} {
	if matchType == "favicon" {
		return map[string]interface{}{
			"type": "favicon",
			"hash": []string{values[0]},
		}
	}

	if values[0] == "()" {
		return nil
	}

	return map[string]interface{}{
		"type":             "word",
		"words":            values,
		"part":             matchType,
		"condition":        mode,
		"case-insensitive": true,
	}
}

// writeToFile 将数据写入指定文件
// 参数:
//   - filepath: 文件路径
//   - data: 数据
//
// 返回值:
//   - error: 可能的错误
func writeToFile(filepath string, data []byte) error {
	return os.WriteFile(filepath, data, 0644)
}

// Json2Json 将一个JSON文件转换为另一种格式的JSON文件
// 参数:
//   - filepath: 源JSON文件路径
//   - newJsonFile: 目标JSON文件路径
//
// 返回值:
//   - []byte: 转换后的JSON字节数组
//   - error: 可能的错误
func Json2Json(filepath, newJsonFile string, logger *logger.Logger) []byte {
	logger.Infof("开始转换JSON文件")
	jsonData, err := readJsonFile(filepath)
	if err != nil {
		logger.Errorf("读取源文件失败: %s", err)
		return nil
	}

	// 使用MarshalIndent进行格式化输出
	data, err := json.Marshal(json2Json(jsonData))
	if err != nil {
		logger.Errorf("格式化JSON失败: %s", err)
		return nil
	}

	// 下面的代码是为了将 map 的字段顺序写入文件
	// 转换为Tag结构
	var tags []pkg.Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		logger.Errorf("验证JSON结构失败: %s", err)
		return nil
	}

	// 顺序输出
	tagsData, err := json.Marshal(tags)
	if err != nil {
		logger.Errorf("验证JSON结构失败: %s", err)
		return nil
	}

	if err := writeToFile(newJsonFile, tagsData); err != nil {
		logger.Errorf("写入目标文件失败: %s", err)
		return nil
	}

	logger.Infof("Json文件转换成功!")

	return data
}
