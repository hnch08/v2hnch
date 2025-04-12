package main

import (
	"context"
	"fmt"
	"net/http"
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

	a.StartProxy()

	conf := config.GetConfig()
	if conf.RequestURL == "" {
		a.ShowWindow()
	}

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

		show.Click(func() { a.ShowWindow() })
		toggle.Click(func() { a.toggleProxy() })
		exit.Click(func() { a.Quit() })

		systray.SetOnClick(func(menu systray.IMenu) { a.ShowWindow() })
		systray.SetOnRClick(func(menu systray.IMenu) { menu.ShowMenu() })
	}
	systray.Run(systemTray, func() {})
}

func (a *App) beforeClose(ctx context.Context) bool {
	a.StopProxy()
	return false
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetConfig() *config.Config {
	return config.GetConfig()
}

func (a *App) SetAddress(address string) bool {
	url := fmt.Sprintf("http://%s:3060/api/health", address)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return false
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	conf := config.GetConfig()
	conf.RequestURL = address
	conf.Username = conf.Username + "1"
	config.Write(conf)
	return true
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
	if err != nil {
		return false
	}
	// runtime.EventsEmit(a.ctx, "proxyStatusChange", new_status)
	return true
}

func (a *App) onSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	a.ShowWindow()
	a.StartProxy()
	// runtime.WindowShow(a.ctx)
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
