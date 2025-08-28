package main

import (
	"encoding/json"
	"fmt"
	"mmo_game/internal/core/entity"
	ientity "mmo_game/internal/icore/entity"
	"zinx/ziface"
	"zinx/znet"
)

// PlayerRouter 自定义玩家路由（对应Player相关业务）
type PlayerRouter struct {
	znet.BaseRouter
}

func (r *PlayerRouter) Handle(request ziface.IRequest) {
	// 解析客户端发送的玩家数据
	data := string(request.GetData())
	fmt.Printf("[收到玩家消息]: %s\n", data)

	// 回复客户端（例如同步玩家状态）
	err := request.GetConnection().SendMsg(ientity.MsgIDPlayerData, []byte("[玩家数据已处理]: "+data))
	if err != nil {
		fmt.Println("[回复失败]:", err)
		request.GetConnection().Stop() // 发送失败时关闭连接
	}
}

type LoginRouter struct {
	znet.BaseRouter
}

func (r *LoginRouter) Handle(request ziface.IRequest) {
	// 解析客户端发送的登录数据（JSON格式）
	var loginReq struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
	if err := json.Unmarshal(request.GetData(), &loginReq); err != nil {
		fmt.Println("login request unmarshal error:", err)
		err := request.GetConnection().SendMsg(ientity.MsgIDLogin, []byte("invalid login format"))
		if err != nil {
			return
		}
	}
	// 模拟玩家创建（实际项目中应从数据库验证用户）
	player := entity.NewPlayer(10001)
	player.SetConn(request.GetConnection())
	if err := player.Connect(); err != nil {
		fmt.Println("player connect error:", err)
		err := request.GetConnection().SendMsg(ientity.MsgIDLogin, []byte("connect failed"))
		if err != nil {
			return
		}
	}

	// 执行登录逻辑
	if player.Login() {
		err := request.GetConnection().SendMsg(ientity.MsgIDLogin, []byte("login success"))
		if err != nil {
			return
		}
	} else {
		err := request.GetConnection().SendMsg(ientity.MsgIDLogin, []byte("login failed"))
		if err != nil {
			return
		}
	}
}

// TalkRouter 新增聊天消息路由
type TalkRouter struct {
	znet.BaseRouter
}

func (r *TalkRouter) Handle(request ziface.IRequest) {
	// 实际项目中应通过连接映射获取玩家（此处简化）
	player := entity.NewPlayer(10001)

	player.SetConn(request.GetConnection())

	// 读取聊天内容并处理
	content := string(request.GetData())
	if content == "" {
		err := request.GetConnection().SendMsg(ientity.MsgIDTalk, []byte("empty content"))
		if err != nil {
			return
		}
	}

	// 执行聊天逻辑（内部会广播消息）
	player.Talk(content)
}

func main() {
	// 创建服务器实例
	s := znet.NewServer()

	// 注册玩家相关的消息路由（假设消息ID=1对应玩家业务）
	s.AddRouter(ientity.MsgIDPlayerData, &PlayerRouter{})
	s.AddRouter(ientity.MsgIDLogin, &LoginRouter{})
	s.AddRouter(ientity.MsgIDTalk, &TalkRouter{})

	// 启动服务器
	s.Serve()
}
