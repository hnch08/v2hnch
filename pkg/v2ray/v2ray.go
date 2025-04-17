package v2ray

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"sync"
	"v2hnch/pkg/api"
	"v2hnch/pkg/config"

	core "github.com/v2fly/v2ray-core/v5"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/inbound"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/outbound"
	_ "github.com/v2fly/v2ray-core/v5/main/formats"
)

var (
	mu      sync.Mutex  // 保护对服务器实例的并发访问
	_server core.Server // V2Ray 服务器实例
)

//go:embed v2ray_config.json
var conf embed.FS // 嵌入的 V2Ray 配置文件

// StartV2Ray 启动 V2Ray 服务器
func StartV2Ray(data []byte) (core.Server, error) {
	configFormat := "json"                               // 配置格式为 JSON
	reader := bytes.NewReader(data)                      // 创建字节读取器
	config, err := core.LoadConfig(configFormat, reader) // 加载配置
	if err != nil {
		fmt.Println(err) // 打印错误信息
		return nil, err  // 返回错误
	}

	// 创建新的 V2Ray 服务器实例
	server, err := core.New(config)
	if err != nil {
		fmt.Println(err) // 打印错误信息
		return nil, err  // 返回错误
	}
	return server, nil // 返回服务器实例
}

// Start 启动 V2Ray 服务器
func Start() error {
	api.GetVpnConfig()
	configure, err := conf.ReadFile("v2ray_config.json") // 读取嵌入的配置文件
	if err != nil {
		return err // 返回错误
	}
	mu.Lock()         // 获取锁以保护对服务器实例的访问
	defer mu.Unlock() // 确保在函数结束时释放锁

	configure, err = updateConfig(configure)
	if err != nil {
		return err
	}

	// 如果服务器已经在运行，先关闭它
	if _server != nil {
		_server.Close() // 关闭当前服务器
		_server = nil   // 重置服务器实例
	}

	// 启动新的 V2Ray 服务器
	server, err := StartV2Ray(configure)
	if err != nil {
		return err // 返回错误
	}

	// 启动服务器
	if err := server.Start(); err != nil {
		return err // 返回错误
	}

	_server = server // 保存服务器实例
	return nil       // 返回 nil 表示成功
}

// Stop 停止 V2Ray 服务器
func Stop() {
	mu.Lock()         // 获取锁以保护对服务器实例的访问
	defer mu.Unlock() // 确保在函数结束时释放锁

	// 如果服务器正在运行，关闭它
	if _server != nil {
		_server.Close() // 关闭服务器
		_server = nil   // 重置服务器实例
	}
}

// updateConfig 更新配置文件中的地址、端口和ID
func updateConfig(data []byte) ([]byte, error) {
	// 获取最新的配置
	cm := config.GetConfigManager()

	conf := cm.GetConfig()

	// 将字节数组转换为字符串以便进行替换
	configStr := string(data)

	// 替换地址、端口和ID
	if conf.Address == "" || conf.Port == "" || conf.Id == "" {
		return nil, fmt.Errorf("配置信息不完整")
	}
	configStr = strings.Replace(configStr, `"address": "address"`, fmt.Sprintf(`"address": "%s"`, conf.Address), 1)
	configStr = strings.Replace(configStr, `"port": "port"`, fmt.Sprintf(`"port": %s`, conf.Port), 1)
	configStr = strings.Replace(configStr, `"id": "id"`, fmt.Sprintf(`"id": "%s"`, conf.Id), 1)

	return []byte(configStr), nil
}
