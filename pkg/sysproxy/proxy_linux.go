//go:build linux

package sysproxy

import (
	"fmt"
	"os/exec"
	"strings"
	"v2hnch/pkg/logger"
)

// SetProxy 设置系统代理（仅针对浏览器）
func SetProxy(addr string) error {
	logger.Info("开始设置Linux浏览器代理")

	addr = "127.0.0.1:2081"
	println(addr)

	// 解析代理地址
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}

	// 解析主机和端口
	host, port, err := parseProxyAddr(addr)
	if err != nil {
		return fmt.Errorf("解析代理地址失败: %v", err)
	}

	// 设置GNOME桌面代理
	if err := setGnomeProxy(host, port); err != nil {
		logger.Error("设置GNOME桌面代理失败: %v", err)
		// 如果GNOME设置失败，尝试KDE
		if err := setKDEProxy(host, port); err != nil {
			logger.Error("设置KDE桌面代理失败: %v", err)
			return fmt.Errorf("设置桌面代理失败，请手动在浏览器中设置代理")
		}
	}

	logger.Info("Linux浏览器代理设置完成")
	return nil
}

// UnSetProxy 取消系统代理设置
func UnSetProxy() error {
	logger.Info("开始取消Linux浏览器代理设置")

	// 清除GNOME桌面代理
	if err := unsetGnomeProxy(); err != nil {
		logger.Error("清除GNOME桌面代理失败: %v", err)
		// 如果GNOME清除失败，尝试KDE
		if err := unsetKDEProxy(); err != nil {
			logger.Error("清除KDE桌面代理失败: %v", err)
			return fmt.Errorf("清除桌面代理失败，请手动在浏览器中清除代理")
		}
	}

	logger.Info("Linux浏览器代理取消完成")
	return nil
}

// parseProxyAddr 解析代理地址
func parseProxyAddr(addr string) (host, port string, err error) {
	// 移除协议前缀
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")

	// 分割主机和端口
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("无效的代理地址格式，应为 host:port")
	}

	return parts[0], parts[1], nil
}

// setGnomeProxy 设置GNOME桌面代理
func setGnomeProxy(host, port string) error {
	// 检查是否有gsettings命令
	if _, err := exec.LookPath("gsettings"); err != nil {
		return fmt.Errorf("gsettings命令不存在")
	}

	// 设置GNOME代理
	commands := [][]string{
		{"gsettings", "set", "org.gnome.system.proxy", "mode", "manual"},
		{"gsettings", "set", "org.gnome.system.proxy.http", "host", host},
		{"gsettings", "set", "org.gnome.system.proxy.http", "port", port},
		{"gsettings", "set", "org.gnome.system.proxy.https", "host", host},
		{"gsettings", "set", "org.gnome.system.proxy.https", "port", port},
		{"gsettings", "set", "org.gnome.system.proxy", "ignore-hosts", "['localhost', '127.0.0.0/8', '::1', '10.0.0.0/8', '192.168.0.0/16', '172.16.0.0/12']"},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("执行GNOME代理设置命令失败: %v", err)
		}
	}

	return nil
}

// unsetGnomeProxy 清除GNOME桌面代理
func unsetGnomeProxy() error {
	if _, err := exec.LookPath("gsettings"); err != nil {
		return fmt.Errorf("gsettings命令不存在")
	}

	// 禁用GNOME代理
	cmd := exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "none")
	return cmd.Run()
}

// setKDEProxy 设置KDE桌面代理
func setKDEProxy(host, port string) error {
	// 检查是否有kwriteconfig5命令
	if _, err := exec.LookPath("kwriteconfig5"); err != nil {
		return fmt.Errorf("kwriteconfig5命令不存在")
	}

	// 设置KDE代理
	commands := [][]string{
		{"kwriteconfig5", "--file", "kioslaverc", "--group", "Proxy Settings", "--key", "ProxyType", "1"},
		{"kwriteconfig5", "--file", "kioslaverc", "--group", "Proxy Settings", "--key", "httpProxy", fmt.Sprintf("http://%s:%s", host, port)},
		{"kwriteconfig5", "--file", "kioslaverc", "--group", "Proxy Settings", "--key", "httpsProxy", fmt.Sprintf("http://%s:%s", host, port)},
		{"kwriteconfig5", "--file", "kioslaverc", "--group", "Proxy Settings", "--key", "NoProxyFor", "localhost,127.0.0.1,::1,10.0.0.0/8,192.168.0.0/16,172.16.0.0/12"},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("执行KDE代理设置命令失败: %v", err)
		}
	}

	// 通知KDE重新加载配置
	exec.Command("dbus-send", "--type=signal", "/KIO/Scheduler", "org.kde.KIO.Scheduler.reparseSlaveConfiguration", "string:").Run()

	return nil
}

// unsetKDEProxy 清除KDE桌面代理
func unsetKDEProxy() error {
	if _, err := exec.LookPath("kwriteconfig5"); err != nil {
		return fmt.Errorf("kwriteconfig5命令不存在")
	}

	// 禁用KDE代理
	cmd := exec.Command("kwriteconfig5", "--file", "kioslaverc", "--group", "Proxy Settings", "--key", "ProxyType", "0")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("清除KDE代理失败: %v", err)
	}

	// 通知KDE重新加载配置
	exec.Command("dbus-send", "--type=signal", "/KIO/Scheduler", "org.kde.KIO.Scheduler.reparseSlaveConfiguration", "string:").Run()

	return nil
}

// GetProxyStatus 获取当前代理状态
func GetProxyStatus() (bool, string, error) {
	// 检查GNOME代理状态
	if _, err := exec.LookPath("gsettings"); err == nil {
		cmd := exec.Command("gsettings", "get", "org.gnome.system.proxy", "mode")
		output, err := cmd.Output()
		if err == nil {
			mode := strings.Trim(string(output), "'\n ")
			if mode == "manual" {
				// 获取代理地址
				hostCmd := exec.Command("gsettings", "get", "org.gnome.system.proxy.http", "host")
				portCmd := exec.Command("gsettings", "get", "org.gnome.system.proxy.http", "port")

				hostOutput, hostErr := hostCmd.Output()
				portOutput, portErr := portCmd.Output()

				if hostErr == nil && portErr == nil {
					host := strings.Trim(string(hostOutput), "'\n ")
					port := strings.Trim(string(portOutput), "'\n ")
					return true, fmt.Sprintf("%s:%s", host, port), nil
				}
			}
		}
	}

	return false, "", nil
}
