package pkg

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/enenisme/definger/logger"
)

// Probe 定义探针的结构
type Probe struct {
	Data            string `toml:"data"`
	Desc            string `toml:"desc"`
	Format          string `toml:"format"`
	Timeout         int    `toml:"timeout"`
	WriteExpression string `toml:"write_expression"`
}

// Probes 定义配置结构
type Probes struct {
	Probes map[string]Probe `toml:"probes"`
}

// HttpResponse 定义HTTP响应结构
type HttpResponse struct {
	Status     string
	StatusCode int
	Header     http.Header
	Body       []byte
}

// HttpRequest 发送HTTP请求到指定URL
// 参数:
//   - url: 目标URL地址
//   - probe: 探针配置信息
//
// 返回:
//   - *HttpResponse: 自定义的HTTP响应结构
//   - error: 错误信息
func (p *Probes) HttpRequest(url string, probe Probe) (*HttpResponse, error) {
	return p.sendHTTPRequest(url, probe)
}

// sendHTTPRequest 发送探针请求到指定URL
// 参数:
//   - url: 目标URL地址s
//   - probe: 探针配置信息
//
// 返回:
//   - *HttpResponse: 自定义的HTTP响应结构
//   - error: 错误信息
func (p *Probes) sendHTTPRequest(url string, probe Probe) (*HttpResponse, error) {
	// 预分配合适大小的切片避免多次扩容
	lines := strings.Split(probe.Data, "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("探针数据无效")
	}

	// 解析请求行
	requestLine := lines[0]
	parts := strings.SplitN(requestLine, " ", 3) // 限制分割次数提高性能
	if len(parts) < 2 {
		return nil, fmt.Errorf("请求行格式无效")
	}

	// 使用strings.Builder拼接URL,避免字符串拼接的内存分配
	var urlBuilder strings.Builder
	urlBuilder.WriteString(url)
	urlBuilder.WriteString(parts[1])
	finalURL := urlBuilder.String()

	// 复用resty客户端以减少资源消耗
	client := resty.New().
		SetTimeout(time.Duration(probe.Timeout) * time.Second).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS13,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			},
		}).
		// 设置重试策略
		SetRetryCount(1).                     // 最多重试1次
		SetRetryWaitTime(2 * time.Second).    // 重试等待1秒
		SetRetryMaxWaitTime(5 * time.Second). // 最大重试等待3秒
		SetLogger(&logger.Logger{Level: logger.LogLevelError}).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, nil
		}).
		// 设置重试条件
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// 增加对TLS握手超时的重试判断
			if err != nil {
				if strings.Contains(err.Error(), "TLS handshake timeout") {
					return true
				}
				return true
			}
			return r.StatusCode() >= 500 // 服务器错误时重试
		}).
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))

	if probe.Timeout <= 0 {
		client.SetTimeout(30 * time.Second)
	}

	// 预分配headers容量,避免map扩容
	headers := make(map[string]string, len(lines)-1)
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		if headerParts := strings.SplitN(line, ": ", 2); len(headerParts) == 2 {
			headers[headerParts[0]] = headerParts[1]
		}
	}

	// 执行请求
	resp, err := client.R().
		SetHeaders(headers).
		Execute(parts[0], finalURL)
	if err != nil {
		return nil, fmt.Errorf("请求执行失败: %v", err)
	}

	return requestHandle(resp)
}

// requestHandle 处理HTTP响应
// 参数:
//   - resp: resty的HTTP响应
//
// 返回:
//   - *HttpResponse: 自定义的HTTP响应结构
//   - error: 错误信息
func requestHandle(resp *resty.Response) (*HttpResponse, error) {
	// 创建新的http.Response
	httpResp := &HttpResponse{
		Status:     resp.Status(),
		StatusCode: resp.StatusCode(),
		Header:     resp.Header(),
		Body:       resp.Body(),
	}
	return httpResp, nil
}
