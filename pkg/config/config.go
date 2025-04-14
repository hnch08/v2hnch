package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// 配置目录和文件的路径
var (
	ConfigDir  string // 配置目录
	ConfigFile string // 配置文件
)

// 初始化配置目录和文件路径
func init() {
	ConfigDir = filepath.Join(os.Getenv("APPDATA"), "v2hnch") // 获取应用数据目录
	ConfigFile = filepath.Join(ConfigDir, "config.json")      // 配置文件路径
}

var instance *Config // 单例配置实例
var once sync.Once   // 确保单例初始化的同步机制

// GetConfig 获取配置的单例实例
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{} // 创建新的配置实例
		// 确保配置目录存在
		if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
			err = os.MkdirAll(ConfigDir, 0755) // 创建配置目录
			if err != nil {
				fmt.Printf("创建配置目录失败: %v\n", err)
				return
			}
		}

		// 确保配置文件存在
		if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
			defaultConfig := &Config{
				Username:   "",
				Name:       "",
				RequestURL: "",
			}
			err = Write(defaultConfig) // 写入默认配置
			if err != nil {
				fmt.Printf("写入默认配置失败: %v\n", err)
				return
			}
		}

		config, err := Read() // 读取配置文件
		if err != nil {
			fmt.Printf("读取配置失败: %v\n", err)
			instance = &Config{} // 如果读取失败，返回一个空的配置实例
			return
		}
		instance = config // 设置单例实例为读取的配置
	})

	return instance // 返回单例配置实例
}

// Config 结构体定义配置项
type Config struct {
	Username   string `json:"username"`   // 用户名
	Name       string `json:"name"`       // 名称
	RequestURL string `json:"requestURL"` // 请求URL
}

// Write 将配置写入文件
func Write(config *Config) error {
	// 将配置序列化为JSON格式
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入配置文件
	err = os.WriteFile(ConfigFile, configData, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}
	fmt.Println("写入配置文件成功")
	return nil
}

// Read 从文件中读取配置
func Read() (*Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %v", err)
	}

	// 读取配置文件内容
	file, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON配置
	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}
	fmt.Println("config:", config)

	return config, nil // 返回读取的配置
}

// SetValue 异步设置单个配置值
func SetValue(key string, value string) chan error {
	result := make(chan error) // 创建结果通道

	go func() {
		config, err := Read() // 读取当前配置
		if err != nil {
			fmt.Println("读取配置失败:", err)
			result <- fmt.Errorf("读取配置失败: %v", err)
			return
		}

		// 使用反射动态设置字段值
		v := reflect.ValueOf(config).Elem()
		field := v.FieldByName(cases.Title(language.Und).String(key)) // 获取字段
		fmt.Println("field:", field)
		if !field.IsValid() {
			fmt.Println("未知的配置项:", key)
			result <- fmt.Errorf("未知的配置项: %s", key)
			return
		}

		if field.Kind() == reflect.String {
			fmt.Println("设置配置项:", key, value)
			field.SetString(value) // 设置字段值
		} else {
			fmt.Printf("配置项 %s 类型不是字符串\n", key)
			result <- fmt.Errorf("配置项 %s 类型不是字符串", key)
			return
		}

		if err := Write(config); err != nil { // 写入更新后的配置
			fmt.Println("写入配置失败:", err)
			result <- err
			return
		}

		result <- nil // 返回成功
	}()

	fmt.Println("设置配置完成")
	return result // 返回结果通道
}
