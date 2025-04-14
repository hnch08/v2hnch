package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func CheckURL(address string) bool {
	serverAddress := ""
	// isDomain := false
	result := CheckInputType(address)
	fmt.Printf("输入类型: %v\n", result)
	switch result {
	case "ipv4", "ipv6":
		serverAddress = fmt.Sprintf("http://%s:3060/api/health", address)
	case "domain":
		serverAddress = fmt.Sprintf("https://%s/api/health", address)
	default:
		fmt.Printf("输入无效: %v\n", result)
		return false
	}
	// url := fmt.Sprintf("http://%s:3060/api/health", address)

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(serverAddress)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return false
	}
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()
	fmt.Printf("响应: %v\n", string(s))
	return string(s) == "\"OK\""
}

func CheckInputType(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return "empty"
	}

	ipv4Pattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	domainPattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`

	if regexp.MustCompile(ipv4Pattern).MatchString(input) {
		if ip := net.ParseIP(input); ip != nil && ip.To4() != nil {
			return "ipv4"
		}
	}
	if ip := net.ParseIP(input); ip != nil {
		return "ipv6"
	}
	if regexp.MustCompile(domainPattern).MatchString(input) {
		_, err := net.LookupHost(input)
		if err == nil {
			return "domain"
		}
	}
	return "invalid"

}
