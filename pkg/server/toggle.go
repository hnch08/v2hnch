package server

import (
	"fmt"
	"v2hnch/pkg/api"
	"v2hnch/pkg/config"
	"v2hnch/pkg/sysproxy"
	"v2hnch/pkg/v2ray"
)

func Toggle() {
	status := config.GetStatus()

	switch status {
	case config.StatusActive:
		v2ray.Stop()
		if err := sysproxy.UnSetProxy(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("关闭代理成功")
		config.SetStatus(config.StatusInActive)
	case config.StatusInActive:
		if err := v2ray.Start(); err != nil {
			fmt.Println(err)
		}
		if err := sysproxy.SetProxy(api.PacAddr); err != nil {
			fmt.Println(err)
		}
		fmt.Println("设置代理成功")
		config.SetStatus(config.StatusActive)
	}
}
