package main

import (
	"embed"
	"fmt"
	"os"
	"v2hnch/pkg/config"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist favicon.ico
var assets embed.FS

var Conf *config.Config

func main() {
	// Create an instance of the app structure
	app := NewApp()

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		fmt.Println("argsWithoutProg:", argsWithoutProg)
	}
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "v2hnch",
		Width:  400,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
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
