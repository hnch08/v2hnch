//go:build windows

package sysproxy

import (
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
	// "golang.org/x/sys/windows/registry"
)

// SetProxy 设置系统代理
func SetProxy(addr string) error {
	// 启用代理
	enableCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
	enableCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏命令窗口
	if err := enableCmd.Run(); err != nil {
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr) // 返回错误信息
	}

	// 设置代理服务器地址
	serverCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/t", "REG_SZ", "/d", "127.0.0.1:2081", "/f")
	serverCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏命令窗口
	if err := serverCmd.Run(); err != nil {
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr) // 返回错误信息
	}

	// 下面的代码是使用 Windows 注册表 API 设置代理的示例（已注释）
	// k, err := registry.OpenKey(registry.CURRENT_USER, internetSetting, registry.ALL_ACCESS)
	// if err != nil {
	// 	return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr)
	// }
	// defer k.Close()

	// err = k.SetStringValue("AutoConfigURL", addr)
	// if err != nil {
	// 	return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr)
	// }
	// store.SetProxyStatus(true)
	return nil // 返回 nil 表示成功
}

// UnSetProxy 禁用系统代理
func UnSetProxy() error {
	// 禁用代理
	disableCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
	disableCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏命令窗口
	if err := disableCmd.Run(); err != nil {
		return errors.Wrap(err, "清除系统代理失败, 请手动操作") // 返回错误信息
	}

	// （可选）清除代理服务器设置
	clearCmd := exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/f")
	clearCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 隐藏命令窗口
	if err := clearCmd.Run(); err != nil {
		return errors.Wrap(err, "清除系统代理失败, 请手动操作") // 返回错误信息
	}

	// 下面的代码是使用 Windows 注册表 API 清除代理的示例（已注释）
	// k, err := registry.OpenKey(registry.CURRENT_USER, internetSetting, registry.ALL_ACCESS)
	// if err != nil {
	// 	return errors.Wrap(err, "清除系统代理失败, 请手动操作")
	// }
	// defer k.Close()

	// err = k.DeleteValue("AutoConfigURL")
	// if err != nil {
	// 	return errors.Wrap(err, "清除系统代理失败, 请手动清除")
	// }
	// store.SetProxyStatus(false)
	return nil // 返回 nil 表示成功
}
