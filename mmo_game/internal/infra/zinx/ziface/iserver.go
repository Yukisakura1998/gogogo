package ziface

// IServer 定义服务器接口
type IServer interface {
	// Start 启动服务器方法
	Start()
	// Stop 停止服务器方法
	Stop()
	// Serve 开启业务服务方法
	Serve()
	// AddRouter 注册路由功能
	AddRouter(msgId uint32, router IRouter)
	//GetConnMgr 获取连接管理
	GetConnMgr() IConnManager
	// SetOnConnStart 创建之后的Hook函数
	SetOnConnStart(func(c IConnection))
	// SetOnConnStop 停止之前的Hook函数
	SetOnConnStop(func(c IConnection))
	//CallOnConnStart 执行创建之后的Hook函数
	CallOnConnStart(c IConnection)
	//CallOnConnStop 执行停止之前的Hook函数
	CallOnConnStop(c IConnection)
}
