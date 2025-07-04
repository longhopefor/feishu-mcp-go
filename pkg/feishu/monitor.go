package feishu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// RequestMonitor HTTP请求监控器
type RequestMonitor struct {
	enabled bool
	verbose bool
}

// NewRequestMonitor 创建请求监控器
func NewRequestMonitor(enabled, verbose bool) *RequestMonitor {
	return &RequestMonitor{
		enabled: enabled,
		verbose: verbose,
	}
}

// printRequestBody 安全地打印请求体，处理不同的数据类型
func (m *RequestMonitor) printRequestBody(body interface{}) {
	switch v := body.(type) {
	case []byte:
		fmt.Printf("📄 Body: %s\n", string(v))
	case string:
		fmt.Printf("📄 Body: %s\n", v)
	case map[string]interface{}, map[string]string:
		if jsonBytes, err := json.Marshal(v); err == nil {
			fmt.Printf("📄 Body: %s\n", string(jsonBytes))
		} else {
			fmt.Printf("📄 Body: %+v\n", v)
		}
	default:
		if jsonBytes, err := json.Marshal(v); err == nil {
			fmt.Printf("📄 Body: %s\n", string(jsonBytes))
		} else {
			fmt.Printf("📄 Body: %+v (type: %T)\n", v, v)
		}
	}
}

// EnableMonitoring 为resty客户端启用请求监控
func (m *RequestMonitor) EnableMonitoring(client *resty.Client) *resty.Client {
	if !m.enabled {
		return client
	}

	// 请求前置处理
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		if m.verbose {
			fmt.Printf("🌐 HTTP请求监控\n")
			fmt.Printf("📍 URL: %s %s\n", req.Method, req.URL)
			fmt.Printf("🔑 Headers: %+v\n", req.Header)
			if req.Body != nil {
				m.printRequestBody(req.Body)
			}
			fmt.Printf("⏰ 时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		}
		return nil
	})

	// 响应后处理
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		duration := resp.Time()

		if m.verbose {
			fmt.Printf("📥 HTTP响应监控\n")
			fmt.Printf("📊 状态码: %d\n", resp.StatusCode())
			fmt.Printf("⏱️ 耗时: %v\n", duration)
			fmt.Printf("📏 响应大小: %d bytes\n", len(resp.Body()))
			fmt.Printf("🔖 响应头: %+v\n", resp.Header())

			// 只显示前500个字符的响应内容
			body := string(resp.Body())
			if len(body) > 500 {
				body = body[:500] + "..."
			}
			fmt.Printf("📄 响应内容: %s\n", body)
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		}

		// 检查错误状态
		if resp.StatusCode() >= 400 {
			m.logError(resp)
		}

		return nil
	})

	// 错误处理
	client.OnError(func(req *resty.Request, err error) {
		fmt.Printf("❌ HTTP请求错误\n")
		fmt.Printf("📍 URL: %s %s\n", req.Method, req.URL)
		fmt.Printf("🔥 错误: %v\n", err)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	})

	return client
}

// logError 记录错误详情
func (m *RequestMonitor) logError(resp *resty.Response) {
	fmt.Printf("🚨 API错误响应\n")
	fmt.Printf("📍 URL: %s\n", resp.Request.URL)
	fmt.Printf("📊 状态码: %d\n", resp.StatusCode())
	fmt.Printf("📄 错误内容: %s\n", string(resp.Body()))

	// 根据状态码提供调试建议
	switch resp.StatusCode() {
	case 400:
		fmt.Println("💡 调试建议: 检查请求参数是否正确")
	case 401:
		fmt.Println("💡 调试建议: 检查应用ID和密钥是否正确，token是否过期")
	case 403:
		fmt.Println("💡 调试建议: 检查应用权限配置是否正确")
	case 404:
		fmt.Println("💡 调试建议: 检查API端点是否正确，资源是否存在")
	case 429:
		fmt.Println("💡 调试建议: 请求频率过高，建议增加重试间隔")
	case 500:
		fmt.Println("💡 调试建议: 飞书服务器内部错误，稍后重试")
	default:
		fmt.Println("💡 调试建议: 检查飞书开发者文档获取更多信息")
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// HealthCheck 执行健康检查
func (c *Client) HealthCheck() error {
	fmt.Println("🏥 执行API健康检查...")

	// 检查1: 获取token（这也会检查基础连接）
	fmt.Println("1️⃣ 检查Token获取...")
	token, err := c.GetToken()
	if err != nil {
		return fmt.Errorf("Token获取失败: %w", err)
	}
	fmt.Printf("   ✅ Token获取成功，长度: %d\n", len(token))

	// 检查2: 测试基础API调用 - 使用一个已知的正确API端点
	fmt.Println("2️⃣ 检查API基础连通性...")

	// 尝试一个简单的健康检查API调用
	resp, err := c.client.R().
		SetAuthToken(token).
		Get(c.baseURL + "/auth/v3/info")

	if err != nil {
		fmt.Printf("   ⚠️ API连接测试失败: %v\n", err)
	} else if resp.StatusCode() == 200 {
		fmt.Printf("   ✅ API基础连通性正常\n")
	} else {
		fmt.Printf("   ℹ️ API响应状态码: %d (Token有效，但某些功能可能需要额外权限)\n", resp.StatusCode())
	}

	// 检查3: 验证关键配置
	fmt.Println("3️⃣ 检查应用配置...")
	if c.baseURL == "" {
		return fmt.Errorf("基础URL未配置")
	}
	fmt.Printf("   ✅ 基础URL配置正常: %s\n", c.baseURL)

	fmt.Println("🎉 健康检查完成!")
	return nil
}

// APIValidator API验证器
type APIValidator struct {
	client *Client
}

// NewAPIValidator 创建API验证器
func NewAPIValidator(client *Client) *APIValidator {
	return &APIValidator{client: client}
}

// ValidateEndpoints 验证所有API端点
func (v *APIValidator) ValidateEndpoints() {
	fmt.Println("🔍 开始验证API端点...")

	endpoints := []struct {
		name   string
		method string
		path   string
		test   func() error
	}{
		{
			name:   "获取访问令牌",
			method: "POST",
			path:   "/auth/v3/tenant_access_token/internal",
			test: func() error {
				_, err := v.client.GetToken()
				return err
			},
		},
		{
			name:   "搜索文档",
			method: "GET",
			path:   "/search/v2/data_source",
			test: func() error {
				_, err := v.client.SearchDocuments("test", 1, "")
				return err
			},
		},
		// 可以添加更多端点...
	}

	for _, endpoint := range endpoints {
		fmt.Printf("📍 验证 %s (%s %s)\n", endpoint.name, endpoint.method, endpoint.path)

		err := endpoint.test()
		if err != nil {
			fmt.Printf("   ❌ 验证失败: %v\n", err)
		} else {
			fmt.Printf("   ✅ 验证成功\n")
		}
	}

	fmt.Println("🎯 API端点验证完成!")
}
