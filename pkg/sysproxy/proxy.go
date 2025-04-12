//go:build windows

package sysproxy

import (
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
	// "golang.org/x/sys/windows/registry"
)

// const (
// 	internetSetting = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`
// )

func SetProxy(addr string) error {
	// 启用代理
	enableCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
	enableCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := enableCmd.Run(); err != nil {
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr)
	}

	// 设置代理服务器
	serverCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/t", "REG_SZ", "/d", "127.0.0.1:2081", "/f")
	serverCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := serverCmd.Run(); err != nil {
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr)
	}
	// return nil
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
	return nil
}

func UnSetProxy() error {
	// 禁用代理
	disableCmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
	disableCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := disableCmd.Run(); err != nil {
		return errors.Wrap(err, "清楚系统代理失败, 请手动操作")
	}

	// （可选）清除代理服务器设置
	clearCmd := exec.Command("reg", "delete", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", "/v", "ProxyServer", "/f")
	clearCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := clearCmd.Run(); err != nil {
		return errors.Wrap(err, "清楚系统代理失败, 请手动操作")
	}
	// k, err := registry.OpenKey(registry.CURRENT_USER, internetSetting, registry.ALL_ACCESS)
	// if err != nil {
	// 	return errors.Wrap(err, "清楚系统代理失败, 请手动操作")
	// }
	// defer k.Close()

	// err = k.DeleteValue("AutoConfigURL")
	// if err != nil {
	// 	return errors.Wrap(err, "清除系统代理失败, 请手动清除")
	// }
	// store.SetProxyStatus(false)
	return nil
}
