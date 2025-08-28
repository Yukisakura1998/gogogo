package entity

import "zinx/ziface"

type IPlayer interface {
	GetPlayerID() int
	GetName() string
	GetLevel() int
	GetConn() ziface.IConnection
	GetStatus() PlayerState
	GetCurrentGameID() int
	GetGold() int // 金币

	SetPlayerID(int)
	SetName(string)
	SetLevel(int)
	SetConn(ziface.IConnection)
	SetStatus(PlayerState)
	SetCurrentGameID(int)
	//SetGold(int) // 金币

	AddGold(amount int)    // 增加金币
	ReduceGold(amount int) // 增加金币

	//Save() error // 保存玩家数据到数据库

	Connect() error
	Disconnect() error
	SendMsg(msgID uint32, data []byte) error

	// Login 业务方法
	Login() bool         // 登录
	Logout()             // 登出
	Talk(content string) // 聊天
}

// PlayerState 玩家状态枚举
type PlayerState int

const (
	PlayerStateIdle     PlayerState = 0 // 空闲状态
	PlayerStateMatching PlayerState = 1 // 匹配中
	PlayerStateFighting PlayerState = 2 // 战斗中
	PlayerStateOffline  PlayerState = 3 // 离线
)

const (
	MsgIDPlayerData uint32 = 1 // 玩家数据消息
	MsgIDLogin      uint32 = 2 // 登录消息
	MsgIDTalk       uint32 = 3 // 聊天消息
)
