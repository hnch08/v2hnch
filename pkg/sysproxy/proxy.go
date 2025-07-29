//go:build windows

package sysproxy

import (
	"os/exec"
	"syscall"
	"v2hnch/pkg/logger"

	"github.com/pkg/errors"
)

// SetProxy 设置系统代理
func SetProxy(addr string) error {
	logger.Info("开始设置Windows系统代理")

	addr = "127.0.0.1:2081"

	// 启用代理
	enableCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
	enableCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := enableCmd.Run(); err != nil {
		logger.Error("启用系统代理失败: %v", err)
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr) // 返回错误信息
	}

	// 设置代理服务器地址
	serverCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/t", "REG_SZ", "/d", addr, "/f")
	serverCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏命令窗口
	if err := serverCmd.Run(); err != nil {
		logger.Error("设置系统代理服务器地址失败: %v", err)
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr) // 返回错误信息
	}

	logger.Info("系统代理设置成功")
	return nil
}

// UnSetProxy 取消系统代理设置
func UnSetProxy() error {
	logger.Info("开始取消Windows系统代理设置")

	// 禁用代理
	disableCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
	disableCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := disableCmd.Run(); err != nil {
		logger.Error("禁用系统代理失败: %v", err)
		return err
	}

	logger.Info("系统代理取消成功")
	return nil
}
