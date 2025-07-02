package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
	wrong "v2hnch/pkg/error"
)

type UserResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CheckURL 检查给定地址的健康状态
func CheckURL(address string) bool {
	serverAddress := ""
	result := CheckInputType(address) // 检查输入类型
	fmt.Printf("输入类型: %v\n", result)

	// 根据输入类型构建服务器地址
	switch result {
	case "ipv4", "ipv6": // 如果是 IPv4 或 IPv6 地址
		serverAddress = fmt.Sprintf("http://%s:3060/api/health", address) // 使用 HTTP 协议
	case "domain": // 如果是域名
		serverAddress = fmt.Sprintf("https://%s/api/health", address) // 使用 HTTPS 协议
	default: // 输入无效
		fmt.Printf("输入无效: %v\n", result)
		return false
	}

	// 创建 HTTP 客户端并设置超时时间
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 发送 GET 请求
	resp, err := client.Get(serverAddress)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return false // 请求失败，返回 false
	}

	// 读取响应体
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取失败: %v\n", err)
		return false // 读取失败，返回 false
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	fmt.Printf("响应: %v\n", string(s))
	return string(s) == "\"OK\"" // 检查响应内容是否为 "OK"
}

// CheckURL 检查给定地址的健康状态
func CheckUser(address, username string) error {
	serverAddress := ""
	result := CheckInputType(address) // 检查输入类型
	fmt.Printf("输入类型: %v\n", result)

	// 根据输入类型构建服务器地址
	switch result {
	case "ipv4", "ipv6": // 如果是 IPv4 或 IPv6 地址
		serverAddress = fmt.Sprintf("http://%s:3060/api/auth/vpnstate?username=%s", address, username) // 使用 HTTP 协议
	case "domain": // 如果是域名
		serverAddress = fmt.Sprintf("https://%s/api/auth/vpnstate?username=%s", address, username) // 使用 HTTPS 协议
	default: // 输入无效
		fmt.Printf("输入无效: %v\n", result)
		return wrong.ErrInvalidAddress
	}

	// 创建 HTTP 客户端并设置超时时间
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 发送 GET 请求
	resp, err := client.Get(serverAddress)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return wrong.ErrConnectionFailed // 请求失败，返回 false
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	// 读取响应体
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取失败: %v\n", err)
		return wrong.ErrReadResponseFailed // 读取失败，返回 false
	}

	var userResponse UserResponse
	err = json.Unmarshal(s, &userResponse)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
		fmt.Printf("原始响应: %v\n", string(s))
		return wrong.ErrJsonParseFailed // 解析失败，返回 false
	}
	fmt.Println("====== JSON解析成功 ======")
	fmt.Printf("解析结果: %v\n", userResponse)

	if userResponse.Success {
		fmt.Printf("用户 %v 状态正常", username)
		return nil
	}
	fmt.Printf("用户 %v 状态异常，可能是用户不存在或未激活\n", username)
	return wrong.ErrUserNotActive
}

// CheckInputType 检查输入的类型（IPv4、IPv6、域名或无效）
func CheckInputType(input string) string {
	input = strings.TrimSpace(input) // 去除输入前后的空格
	if input == "" {
		return "empty" // 输入为空
	}

	// IPv4 和域名的正则表达式模式
	ipv4Pattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	domainPattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`

	// 检查是否为 IPv4 地址
	if regexp.MustCompile(ipv4Pattern).MatchString(input) {
		if ip := net.ParseIP(input); ip != nil && ip.To4() != nil {
			return "ipv4" // 返回 IPv4 类型
		}
	}
	// 检查是否为 IPv6 地址
	if ip := net.ParseIP(input); ip != nil {
		return "ipv6" // 返回 IPv6 类型
	}
	// 检查是否为有效的域名
	if regexp.MustCompile(domainPattern).MatchString(input) {
		_, err := net.LookupHost(input) // 查找域名
		if err == nil {
			return "domain" // 返回域名类型
		}
	}
	return "invalid" // 返回无效类型
}
