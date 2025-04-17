package config

import "sync"

// 定义状态常量
const (
	StatusAuto     = 0 // 自动状态
	StatusActive   = 1 // 活动状态
	StatusInActive = 2 // 非活动状态
)

// 保护状态变量的读写锁
var (
	mu     sync.RWMutex                  // 读写锁，用于保护状态的并发访问
	status int          = StatusInActive // 当前状态，初始为活动状态
)

// GetStatus 获取当前状态
func GetStatus() int {
	mu.RLock()         // 获取读锁
	defer mu.RUnlock() // 确保在函数结束时释放读锁

	return status // 返回当前状态
}

// SetStatus 设置当前状态
func SetStatus(s int) {
	mu.Lock()         // 获取写锁
	defer mu.Unlock() // 确保在函数结束时释放写锁
	status = s        // 更新当前状态
}
