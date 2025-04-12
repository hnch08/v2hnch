package config

import "sync"

const (
	StatusAuto     = 0
	StatusActive   = 1
	StatusInActive = 2
)

var (
	mu     sync.RWMutex
	status int = StatusInActive
)

func GetStatus() int {
	mu.RLock()
	defer mu.RUnlock()

	return status
}
func SetStatus(s int) {
	mu.Lock()
	defer mu.Unlock()
	status = s
}
