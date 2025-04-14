package main

import (
	"context"
	"fmt"
	"v2hnch/pkg/config"
	"v2hnch/pkg/server"

	"github.com/energye/systray"

	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	conf := config.GetConfig()
	if conf.RequestURL == "" && conf.Username != "" {
		a.ShowWindow()
	}
	if conf.RequestURL != "" && conf.Username != "" {
		a.StartProxy()
	}

	systemTray := func() {
		iconBytes, err := assets.ReadFile("favicon.ico")
		if err != nil {
			fmt.Println("读取图标文件失败:", err)
		}
		systray.SetIcon(iconBytes)
		systray.SetTooltip("创合智能")
		show := systray.AddMenuItem("显示窗口", "显示窗口")
		// status := config.GetStatus()
		// toggle := systray.AddMenuItemCheckbox("切换状态", "切换代理状态", status == config.StatusActive)
		systray.AddSeparator()
		exit := systray.AddMenuItem("退出程序", "退出程序")

		show.Click(func() { a.ShowWindow() })
		// toggle.Click(func() {
		// 	if toggle.Checked() {
		// 		a.StopProxy()
		// 		toggle.Uncheck()
		// 	} else {
		// 		a.StartProxy()
		// 		toggle.Check()
		// 	}
		// 	new_status := a.GetStatus()
		// 	runtime.EventsEmit(a.ctx, "proxyStatusChange", new_status)
		// })
		exit.Click(func() {
			systray.Quit()
			a.Quit()
		})

		systray.SetOnClick(func(menu systray.IMenu) { a.ShowWindow() })
		systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })
	}
	systray.Run(systemTray, func() {})

}

func (a *App) beforeClose(ctx context.Context) bool {
	a.StopProxy()
	return false
}

func (a *App) GetConfig() *config.Config {
	return config.GetConfig()
}

func (a *App) SetAddress(address string) bool {
	if !server.CheckURL(address) {
		return false
	}
	conf := config.GetConfig()
	conf.RequestURL = address
	config.Write(conf)
	a.StartProxy()
	return true
}

func (a *App) CheckURL() bool {
	conf := config.GetConfig()
	if conf.RequestURL == "" {
		return false
	}
	return server.CheckURL(conf.RequestURL)
}

func (a *App) GetStatus() int {
	return config.GetStatus()
}

func (a *App) StartProxy() bool {
	return a.toggleProxy(config.StatusActive)
}

func (a *App) StopProxy() bool {
	return a.toggleProxy((config.StatusInActive))
}

func (a *App) toggleProxy(status ...int) bool {
	targetStatus := config.StatusAuto
	if len(status) > 0 {
		targetStatus = status[0]
	}
	_, err := server.Toggle(targetStatus)
	return err == nil
}

func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	conf := config.GetConfig()
	if conf.RequestURL == "" && conf.Username != "" {
		a.ShowWindow()
	}
	if conf.RequestURL != "" && conf.Username != "" {
		a.StartProxy()
		new_status := a.GetStatus()
		runtime.EventsEmit(a.ctx, "proxyStatusChange", new_status)
	}
}

// HideWindow 隐藏主窗口
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}

// ShowWindow 显示主窗口
func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
	runtime.WindowUnminimise(a.ctx)
	runtime.WindowSetAlwaysOnTop(a.ctx, true)
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
}

// Quit 退出应用
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}
