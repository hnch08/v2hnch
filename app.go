package main

import (
	"context"
	"fmt"
	"os"
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

	runtime.WindowSetSize(ctx, 400, 600)
	runtime.WindowCenter(ctx)
	runtime.WindowSetTitle(ctx, "v2hnch")

	Conf = config.GetConfig()
	fmt.Println(Conf)

	systemTray := func() {
		// systray.SetIcon([]byte()) // read the icon from a file
		iconBytes, err := assets.ReadFile("favicon.ico")
		if err != nil {
			fmt.Println("读取图标文件失败:", err)
		}
		systray.SetIcon(iconBytes)
		systray.SetTooltip("创合智能")
		show := systray.AddMenuItem("Show", "Show The Window")
		toggle := systray.AddMenuItem("Toggle", "Toggle The Window")
		systray.AddSeparator()
		exit := systray.AddMenuItem("Exit", "Quit The Program")

		show.Click(func() { runtime.WindowShow(a.ctx) })
		toggle.Click(server.Toggle)
		exit.Click(func() { os.Exit(0) })

		systray.SetOnClick(func(menu systray.IMenu) { runtime.WindowShow(a.ctx) })
		systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })
	}
	systray.Run(systemTray, func() {})
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Login 处理登录请求
func (a *App) Login(username, password string) bool {
	// 这里添加实际的登录验证逻辑
	// 示例中使用简单的判断，实际应用中应该使用更安全的方式
	if username == "admin" && password == "password" {
		config.SetValue("username", username)
		return true
	}
	return false
}

func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	runtime.WindowShow(a.ctx)
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
