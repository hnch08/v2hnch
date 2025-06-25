package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"v2hnch/pkg/config"
	"v2hnch/pkg/logger"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VpnConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	UUID string `json:"uuid"`
}

type VpnConfigResponse struct {
	Data []VpnConfig `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

const PacAddr = "http://cdn.hnch.net/pac.js"

const VpnURL = "http://124.232.157.50:8008/api/vpn/"

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

func GetVpnConfig() {

	logger.Info("正在从服务器获取vpn配置")
	resp, err := http.Get(VpnURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 处理响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	cm := config.GetConfigManager()
	var vpnConfig VpnConfigResponse
	err = json.Unmarshal(body, &vpnConfig)
	if err != nil {
		logger.Error("解析VPN配置失败: %v", err)
		return
	}
	if len(vpnConfig.Data) == 0 {
		logger.Error("获取到的配置为空")
		return
	}

	conf := cm.GetConfig() // 获取配置实例
	conf.Port = vpnConfig.Data[0].Port
	conf.Address = vpnConfig.Data[0].Host
	conf.Id = vpnConfig.Data[0].UUID

	if err := cm.UpdateConfig(conf); err != nil {
		logger.Error("写入配置失败: %v", err)
	}
}
