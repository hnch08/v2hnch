package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"v2hnch/pkg/logger"

	"github.com/fsnotify/fsnotify"
)

// 配置目录和文件的路径
var (
	ConfigDir  string // 配置目录
	ConfigFile string // 配置文件
)

// 初始化配置目录和文件路径
func init() {
	configDir := os.Getenv("APPDATA")
	if configDir == "" {
		// Linux/macOS 回退，通常使用 ~/.config
		configDir = filepath.Join(os.Getenv("HOME"), ".config", "v2hnch")
	} else {
		configDir = filepath.Join(configDir, "v2hnch")
	}
	ConfigDir = configDir                               // 获取应用数据目录
	ConfigFile = filepath.Join(ConfigDir, "config.enc") // 配置文件路径，使用.enc扩展名表示加密文件
}

// ConfigManager 管理配置，包含锁
type ConfigManager struct {
	config  *Config
	mu      sync.RWMutex      // 读写锁
	watcher *fsnotify.Watcher // 文件变化监听器
}

var (
	instance *ConfigManager
	once     sync.Once
)

// GetConfig 获取配置的单例实例
func GetConfigManager() *ConfigManager {
	once.Do(func() {
		logger.Info("开始初始化配置单例-*-*-*-*-*-*-*-*-*-*-*-*")
		instance = &ConfigManager{
			config: &Config{},
		} // 创建新的配置实例

		// 确保配置目录存在
		if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
			err = os.MkdirAll(ConfigDir, 0755) // 创建配置目录
			if err != nil {
				logger.Error("创建配置目录失败: %v", err)
				return
			}
			logger.Info("创建配置目录成功: %s", ConfigDir)
		}

		// 确保配置文件存在
		if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
			defaultConfig := &Config{
				Username:   "",
				Name:       "",
				RequestURL: "",
				Port:       "",
				Address:    "",
				Id:         "",
			}
			err = Write(defaultConfig) // 写入默认配置（已加密）
			if err != nil {
				logger.Error("写入默认配置失败: %v", err)
				return
			}
			logger.Info("写入默认加密配置成功")
		}

		// 读取加密的配置文件
		config, err := Read()
		if err != nil {
			logger.Error("读取加密配置失败: %v", err)
			instance.config = &Config{} // 如果读取失败，返回一个空的配置实例
			return
		}

		instance.config = config // 设置单例实例为读取的配置
		logger.Info("成功初始化配置单例")

		// 初始化文件变化监听
		err = instance.startFileWatcher()
		if err != nil {
			logger.Error("启动文件监听失败: %v", err)
			// 不终止程序，允许继续运行
		}
	})

	return instance
}

func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return &Config{
		Username:   cm.config.Username,
		Name:       cm.config.Name,
		RequestURL: cm.config.RequestURL,
		Port:       cm.config.Port,
		Address:    cm.config.Address,
		Id:         cm.config.Id,
	}
}

// Config 结构体定义配置项
type Config struct {
	Username   string `json:"username"`   // 用户名
	Name       string `json:"name"`       // 名称
	RequestURL string `json:"requestURL"` // 请求URL
	Port       string `json:"port"`
	Address    string `json:"address"`
	Id         string `json:"id"`
}

// Write 将配置加密写入文件
func Write(config *Config) error {
	// 将配置序列化为JSON格式
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 加密并写入配置文件
	err = WriteEncrypted(ConfigFile, configData)
	if err != nil {
		return fmt.Errorf("写入加密配置文件失败: %v", err)
	}

	logger.Info("写入加密配置文件成功")
	return nil
}

// Read 从加密文件中读取配置
func Read() (*Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %v", err)
	}

	// 读取并解密配置文件内容
	file, err := ReadEncrypted(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("读取加密配置文件失败: %v", err)
	}

	// 解析JSON配置
	config := &Config{}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	logger.Info("成功读取配置")
	return config, nil
}

func (cm *ConfigManager) UpdateConfig(newConfig *Config) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if err := Write(newConfig); err != nil {
		logger.Error("写入配置失败: %v", err)
		return err
	}
	cm.config = newConfig
	logger.Info("配置已更新")
	return nil
}

// startFileWatcher 启动文件变化监听
func (cm *ConfigManager) startFileWatcher() error {
	var err error
	cm.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("创建文件监听器失败: %v", err)
	}

	// 启动协程处理文件事件
	go func() {
		for {
			select {
			case event, ok := <-cm.watcher.Events:
				if !ok {
					logger.Error("文件监听通道已关闭")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logger.Info("检测到配置文件变化: %s", event.Name)
					// 重新加载配置
					if err := cm.reloadConfig(); err != nil {
						logger.Error("重新加载配置失败: %v", err)
					} else {
						logger.Info("配置已重新加载")
					}
				}
			case err, ok := <-cm.watcher.Errors:
				if !ok {
					logger.Error("文件监听错误通道已关闭")
					return
				}
				logger.Error("文件监听错误: %v", err)
			}
		}
	}()

	// 添加配置文件到监听
	err = cm.watcher.Add(ConfigFile)
	if err != nil {
		return fmt.Errorf("添加文件监听失败: %v", err)
	}
	logger.Info("已启动对配置文件 %s 的监听", ConfigFile)
	return nil
}

// reloadConfig 重新加载配置
func (cm *ConfigManager) reloadConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	config, err := Read()
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	cm.config = config
	return nil
}

// Close 关闭文件监听器
func (cm *ConfigManager) Close() error {
	if cm.watcher != nil {
		return cm.watcher.Close()
	}
	return nil
}
