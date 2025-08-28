package ziface

import (
	"net"
)

type IConnection interface {
	// Start 启动连接
	Start()
	// Stop 停止链接
	Stop()
	// GetTCPConnection 获取socketTCPConn
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取连接ID
	GetConnID() string
	// RemoteAddr 获取客户端地址
	RemoteAddr() net.Addr
	// SendMsg 发送数据
	SendMsg(msgId uint32, data []byte) error
	// SendBuffMsg 发送数据,有缓冲
	SendBuffMsg(msgId uint32, data []byte) error
	// SetProperty 设置用户参数
	SetProperty(key string, value interface{})
	// GetProperty 获取用户参数
	GetProperty(key string) (value interface{}, err error)
	// RemoveProperty 移除用户参数
	RemoveProperty(key string)
	// SetOnConnStop  关闭时的回调函数
	SetOnConnStop(f func(IConnection))
	// SetOnConnStart 启动时的回调函数
	SetOnConnStart(f func(IConnection))
}

type HandFunc func(*net.TCPConn, []byte, int) error
