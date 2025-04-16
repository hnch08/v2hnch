package server

import (
	"v2hnch/pkg/api"
	"v2hnch/pkg/config"
	"v2hnch/pkg/logger"
	"v2hnch/pkg/sysproxy"
	"v2hnch/pkg/v2ray"
)

// Toggle 切换目标状态
func Toggle(target_status int) (int, error) {
	status := config.GetStatus()
	logger.Info("切换代理状态: 当前状态=%d, 目标状态=%d", status, target_status)

	// 如果当前状态与目标状态相同，直接返回
	if status == target_status {
		logger.Info("当前状态与目标状态相同，无需切换")
		return target_status, nil
	}

	// 根据目标状态执行相应的操作
	switch target_status {
	case config.StatusActive:
		logger.Info("切换到活动状态")
		return target_status, start()
	case config.StatusInActive:
		logger.Info("切换到非活动状态")
		return target_status, stop()
	}

	// 根据当前状态切换到目标状态
	switch status {
	case config.StatusActive:
		logger.Info("从活动状态切换到非活动状态")
		return config.StatusInActive, stop()
	case config.StatusInActive:
		logger.Info("从非活动状态切换到活动状态")
		return config.StatusActive, start()
	}
	return config.StatusAuto, nil
}

// start 启动服务并设置代理
func start() error {
	logger.Info("开始启动V2Ray服务")
	if err := v2ray.Start(); err != nil {
		logger.Error("启动V2Ray服务失败: %v", err)
		return err
	}

	logger.Info("开始设置系统代理")
	if err := sysproxy.SetProxy(api.PacAddr); err != nil {
		logger.Error("设置系统代理失败: %v", err)
		return err
	}

	logger.Info("代理服务启动成功")
	config.SetStatus(config.StatusActive)
	return nil
}

// stop 停止服务并取消代理设置
func stop() error {
	logger.Info("开始停止V2Ray服务")
	v2ray.Stop()

	logger.Info("开始取消系统代理设置")
	if err := sysproxy.UnSetProxy(); err != nil {
		logger.Error("取消系统代理设置失败: %v", err)
		return err
	}

	logger.Info("代理服务停止成功")
	config.SetStatus(config.StatusInActive)
	return nil
}
