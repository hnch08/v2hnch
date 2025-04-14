package server

import (
	"fmt"
	"v2hnch/pkg/api"
	"v2hnch/pkg/config"
	"v2hnch/pkg/sysproxy"
	"v2hnch/pkg/v2ray"
)

// Toggle 切换目标状态
func Toggle(target_status int) (int, error) {
	status := config.GetStatus() // 获取当前状态
	// 如果当前状态与目标状态相同，直接返回
	if status == target_status {
		return target_status, nil
	}
	// 根据目标状态执行相应的操作
	switch target_status {
	case config.StatusActive: // 如果目标状态是活动状态
		return target_status, start() // 启动服务
	case config.StatusInActive: // 如果目标状态是非活动状态
		return target_status, stop() // 停止服务
	}
	// 根据当前状态切换到目标状态
	switch status {
	case config.StatusActive: // 当前状态是活动状态
		return config.StatusInActive, stop() // 切换到非活动状态并停止服务
	case config.StatusInActive: // 当前状态是非活动状态
		return config.StatusActive, start() // 切换到活动状态并启动服务
	}
	return config.StatusAuto, nil // 默认返回自动状态
}

// start 启动服务并设置代理
func start() error {
	if err := v2ray.Start(); err != nil { // 启动 V2Ray
		return err // 返回错误
	}
	if err := sysproxy.SetProxy(api.PacAddr); err != nil { // 设置系统代理
		fmt.Println(err) // 打印错误信息
		return err       // 返回错误
	}
	fmt.Println("设置代理成功")                 // 打印成功信息
	config.SetStatus(config.StatusActive) // 更新状态为活动状态
	return nil                            // 返回 nil 表示成功
}

// stop 停止服务并取消代理设置
func stop() error {
	v2ray.Stop()                                  // 停止 V2Ray
	if err := sysproxy.UnSetProxy(); err != nil { // 取消系统代理
		fmt.Println(err) // 打印错误信息
		return err       // 返回错误
	}
	fmt.Println("关闭代理成功")                   // 打印成功信息
	config.SetStatus(config.StatusInActive) // 更新状态为非活动状态
	return nil                              // 返回 nil 表示成功
}
