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

var (
	ConfigDir  string
	ConfigFile string
)

func init() {
	ConfigDir = filepath.Join(os.Getenv("APPDATA"), "v2hnch")
	ConfigFile = filepath.Join(ConfigDir, "config.json")
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		// 确保配置目录存在
		if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
			err = os.MkdirAll(ConfigDir, 0755)
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
			err = Write(defaultConfig)
			if err != nil {
				fmt.Printf("写入默认配置失败: %v\n", err)
				return
			}
		}

		config, err := Read()
		if err != nil {
			fmt.Printf("读取配置失败: %v\n", err)
			instance = &Config{}
			return
		}
		instance = config
	})

	return instance
}

type Config struct {
	Username   string `json:"username"`
	Name       string `json:"name"`
	RequestURL string `json:"requestURL"`
}

// Write 将配置写入文件
func Write(config *Config) error {
	// 将配置序列化为JSON
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

	// 读取配置文件
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

	return config, nil
}

// SetValue 异步设置单个配置值
func SetValue(key string, value string) chan error {
	result := make(chan error)

	go func() {
		config, err := Read()
		if err != nil {
			fmt.Println("读取配置失败:", err)
			result <- fmt.Errorf("读取配置失败: %v", err)
			return
		}

		// 使用反射动态设置字段值
		v := reflect.ValueOf(config).Elem()
		field := v.FieldByName(cases.Title(language.Und).String(key))
		fmt.Println("field:", field)
		if !field.IsValid() {
			fmt.Println("未知的配置项:", key)
			result <- fmt.Errorf("未知的配置项: %s", key)
			return
		}

		if field.Kind() == reflect.String {
			fmt.Println("设置配置项:", key, value)
			field.SetString(value)
		} else {
			fmt.Printf("配置项 %s 类型不是字符串\n", key)
			result <- fmt.Errorf("配置项 %s 类型不是字符串", key)
			return
		}

		if err := Write(config); err != nil {
			fmt.Println("写入配置失败:", err)
			result <- err
			return
		}

		result <- nil
	}()

	fmt.Println("设置配置完成")
	return result
}
