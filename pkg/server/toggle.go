package server

import (
	"fmt"
	"v2hnch/pkg/api"
	"v2hnch/pkg/config"
	"v2hnch/pkg/sysproxy"
	"v2hnch/pkg/v2ray"
)

func Toggle(target_status int) (int, error) {
	status := config.GetStatus()
	if status == target_status {
		return target_status, nil
	}
	switch target_status {
	case config.StatusActive:
		return target_status, start()
	case config.StatusInActive:
		return target_status, stop()
	}
	switch status {
	case config.StatusActive:
		return config.StatusInActive, stop()
	case config.StatusInActive:
		return config.StatusActive, start()
	}
	return config.StatusAuto, nil
}

func start() error {
	if err := v2ray.Start(); err != nil {
		return err
	}
	if err := sysproxy.SetProxy(api.PacAddr); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("设置代理成功")
	config.SetStatus(config.StatusActive)
	return nil
}

func stop() error {
	v2ray.Stop()
	if err := sysproxy.UnSetProxy(); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("关闭代理成功")
	config.SetStatus(config.StatusInActive)
	return nil
}
