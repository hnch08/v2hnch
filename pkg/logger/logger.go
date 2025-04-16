package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	logFile *os.File
	logger  *log.Logger
	once    sync.Once
)

// Init 初始化日志
func Init(logDir string) error {
	var err error
	once.Do(func() {
		// 确保日志目录存在
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return
		}

		// 生成日志文件名
		logPath := filepath.Join(logDir, fmt.Sprintf("v2hnch_%s.log",
			time.Now().Format("2006-01-02")))

		// 打开日志文件
		logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return
		}

		// 创建logger实例，不包含文件信息，我们将手动添加
		logger = log.New(logFile, "", log.Ldate|log.Ltime)

	})
	return err
}

// getCallerInfo 获取调用者信息
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2) // 跳过两层调用栈以获取实际调用位置
	if !ok {
		file = "unknown"
		line = 0
	}
	// 获取文件的短名称（不包含完整路径）
	short := filepath.Base(file)
	return fmt.Sprintf("%s:%d", short, line)
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if logger != nil {
		caller := getCallerInfo()
		logger.Printf("%s [INFO] "+format, append([]interface{}{caller}, v...)...)
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if logger != nil {
		caller := getCallerInfo()
		logger.Printf("%s [ERROR] "+format, append([]interface{}{caller}, v...)...)
	}
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	if logger != nil {
		caller := getCallerInfo()
		logger.Printf("%s [DEBUG] "+format, append([]interface{}{caller}, v...)...)
	}
}

// Close 关闭日志文件
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
