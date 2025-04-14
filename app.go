package main

import (
	"context"
	"fmt"
	"v2hnch/pkg/config" // 导入配置包
	"v2hnch/pkg/server" // 导入服务器包

	"github.com/energye/systray" // 导入系统托盘库

	"github.com/wailsapp/wails/v2/pkg/options" // 导入 Wails 应用选项
	"github.com/wailsapp/wails/v2/pkg/runtime" // 导入 Wails 运行时
)

// App struct 定义应用结构体
type App struct {
	ctx context.Context // 应用上下文
}

// NewApp creates a new App application struct 创建一个新的 App 应用结构体
func NewApp() *App {
	return &App{} // 返回 App 结构体指针
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
// startup 在应用启动时被调用。上下文被保存
// 这样我们就可以调用运行时方法
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx // 保存应用上下文

	conf := config.GetConfig()                        // 获取配置
	if conf.RequestURL == "" && conf.Username != "" { // 如果请求 URL 为空且用户名不为空
		a.ShowWindow() // 显示窗口
	}
	if conf.RequestURL != "" && conf.Username != "" { // 如果请求 URL 不为空且用户名不为空
		a.StartProxy() // 启动代理
	}

	systemTray := func() { // 定义系统托盘函数
		iconBytes, err := assets.ReadFile("favicon.ico") // 读取图标文件
		if err != nil {                                  // 如果读取图标文件失败
			fmt.Println("读取图标文件失败:", err) // 打印错误信息
		}
		systray.SetIcon(iconBytes)                  // 设置托盘图标
		systray.SetTooltip("创合智能")                  // 设置托盘提示
		show := systray.AddMenuItem("显示窗口", "显示窗口") // 添加显示窗口菜单项
		// status := config.GetStatus()
		// toggle := systray.AddMenuItemCheckbox("切换状态", "切换代理状态", status == config.StatusActive)
		systray.AddSeparator()                      // 添加分隔符
		exit := systray.AddMenuItem("退出程序", "退出程序") // 添加退出程序菜单项

		show.Click(func() { a.ShowWindow() }) // 点击显示窗口菜单项，显示窗口
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
		exit.Click(func() { // 点击退出程序菜单项
			systray.Quit() // 退出系统托盘
			a.Quit()       // 退出应用
		})

		systray.SetOnClick(func(menu systray.IMenu) { a.ShowWindow() })   // 点击托盘图标，显示窗口
		systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() }) // 右键点击托盘图标，显示菜单
	}
	systray.Run(systemTray, func() {}) // 运行系统托盘

}

// beforeClose 在应用关闭前被调用
func (a *App) beforeClose(ctx context.Context) bool {
	a.StopProxy() // 停止代理
	return false  // 返回 false，阻止应用关闭
}

// GetConfig 获取配置
func (a *App) GetConfig() *config.Config {
	return config.GetConfig() // 返回配置
}

// SetAddress 设置地址
func (a *App) SetAddress(address string) bool {
	if !server.CheckURL(address) { // 检查 URL 是否有效
		return false // 如果 URL 无效，返回 false
	}
	conf := config.GetConfig() // 获取配置
	conf.RequestURL = address  // 设置请求 URL
	config.Write(conf)         // 写入配置
	a.StartProxy()             // 启动代理
	return true                // 返回 true
}

// CheckURL 检查 URL 是否有效
func (a *App) CheckURL() bool {
	conf := config.GetConfig() // 获取配置
	if conf.RequestURL == "" { // 如果请求 URL 为空
		return false // 返回 false
	}
	return server.CheckURL(conf.RequestURL) // 检查 URL 是否有效
}

// GetStatus 获取状态
func (a *App) GetStatus() int {
	return config.GetStatus() // 返回状态
}

// StartProxy 启动代理
func (a *App) StartProxy() bool {
	return a.toggleProxy(config.StatusActive) // 切换代理状态为 active
}

// StopProxy 停止代理
func (a *App) StopProxy() bool {
	return a.toggleProxy((config.StatusInActive)) // 切换代理状态为 inactive
}

// toggleProxy 切换代理
func (a *App) toggleProxy(status ...int) bool {
	targetStatus := config.StatusAuto // 默认状态为 auto
	if len(status) > 0 {              // 如果有指定状态
		targetStatus = status[0] // 设置目标状态
	}
	_, err := server.Toggle(targetStatus) // 切换代理状态
	return err == nil                     // 如果没有错误，返回 true
}

// onSecondInstanceLaunch 当应用被第二次启动时调用
func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	conf := config.GetConfig()                        // 获取配置
	if conf.RequestURL == "" && conf.Username != "" { // 如果请求 URL 为空且用户名不为空
		a.ShowWindow() // 显示窗口
	}
	if conf.RequestURL != "" && conf.Username != "" { // 如果请求 URL 不为空且用户名不为空
		a.StartProxy()                                             // 启动代理
		new_status := a.GetStatus()                                // 获取状态
		runtime.EventsEmit(a.ctx, "proxyStatusChange", new_status) // 发送代理状态改变事件
	}
}

// HideWindow 隐藏主窗口
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx) // 隐藏主窗口
}

// ShowWindow 显示主窗口
func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)                  // 显示主窗口
	runtime.WindowUnminimise(a.ctx)            // 取消最小化
	runtime.WindowSetAlwaysOnTop(a.ctx, true)  // 设置窗口置顶
	runtime.WindowSetAlwaysOnTop(a.ctx, false) // 取消窗口置顶
}

// Quit 退出应用
func (a *App) Quit() {
	runtime.Quit(a.ctx) // 退出应用
}
