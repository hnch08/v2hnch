package wrong

import "errors"

var (
	ErrInvalidAddress     = errors.New("无效的服务器地址")
	ErrConnectionFailed   = errors.New("无法连接到服务器")
	ErrReadResponseFailed = errors.New("读取服务器响应失败")
	ErrJsonParseFailed    = errors.New("解析服务器响应失败")
	ErrUserNotActive      = errors.New("用户不存在或未激活")
	ErrUserStatusNormal   = errors.New("用户状态正常，但响应格式不符合预期") // 一个例子，用于更细致的错误
)
