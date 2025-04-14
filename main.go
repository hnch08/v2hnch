package main

import (
	"embed"
	"net/url"
	"os"
	"strings"
	"v2hnch/pkg/config"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist favicon.ico
var assets embed.FS

func main() {
	// 创建应用程序实例
	app := NewApp()

	// 获取命令行参数，去掉程序名称
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		// 去掉开头的协议和结尾的/
		arg_string := strings.TrimPrefix(argsWithoutProg[0], "v2hnch://")
		arg_string = strings.TrimSuffix(arg_string, "/")
		// 分割参数字符串
		args := strings.SplitN(arg_string, "_", 2)

		// 获取配置实例并更新用户名和名称
		conf := config.GetConfig()
		conf.Username = args[0]
		// 解码名称为汉字
		name, err := url.QueryUnescape(args[1])
		if err == nil {
			conf.Name = name
		} else {
			conf.Name = args[1] // 如果解码失败，保留原始值
		}
		// 将更新后的配置写入文件
		config.Write(conf)
	}

	// 创建应用程序并设置选项
	err := wails.Run(&options.App{
		Title:  "湖南创合", // 应用程序标题
		Width:  800,    // 窗口宽度
		Height: 600,    // 窗口高度
		AssetServer: &assetserver.Options{
			Assets: assets, // 嵌入的静态资源
		},
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 1}, // 窗口背景颜色
		OnStartup:         app.startup,                              // 启动时调用的函数
		OnBeforeClose:     app.beforeClose,                          // 关闭前调用的函数
		StartHidden:       true,                                     // 启动时隐藏窗口
		HideWindowOnClose: true,                                     // 关闭时隐藏窗口
		Bind: []interface{}{
			app, // 绑定应用程序实例
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "123456798",                // 唯一标识符
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch, // 第二个实例启动时调用的函数
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(), // 隐藏标题栏
			About: &mac.AboutInfo{
				Title:   "v2hnch",             // 关于窗口标题
				Message: "v2hnch Application", // 关于窗口消息
			},
		},
	})

	// 错误处理
	if err != nil {
		println("Error:", err.Error())
	}
}
