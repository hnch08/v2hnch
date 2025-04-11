package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const PacAddr = "http://cdn.hnch.net/pac.js"

// Login 发送登录请求到指定路径
func Login(username, password, loginURL string) (bool, error) {
	// 构造登录请求体
	loginReq := LoginRequest{
		Username: username,
		Password: password,
	}

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return false, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 判断响应状态码
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}
