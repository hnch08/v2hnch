package main

import (
	"embed"
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
	// Create an instance of the app structure
	app := NewApp()

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		// 去掉开头的协议和结尾的/
		arg_string := strings.TrimPrefix(argsWithoutProg[0], "v2hnch://")
		arg_string = strings.TrimSuffix(arg_string, "/")
		args := strings.SplitN(arg_string, "_", 2)
		conf := config.GetConfig()
		conf.Username = args[0]
		conf.Name = args[1]
		config.Write(conf)
	}
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "湖南创合",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:         app.startup,
		OnBeforeClose:     app.beforeClose,
		StartHidden:       true,
		HideWindowOnClose: true,
		Bind: []interface{}{
			app,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId:               "123456798",
			OnSecondInstanceLaunch: app.onSecondInstanceLaunch,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "v2hnch",
				Message: "v2hnch Application",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
