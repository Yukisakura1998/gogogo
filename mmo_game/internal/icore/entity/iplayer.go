package entity

import "zinx/ziface"

type IPlayer interface {
	GetPlayerID() int
	GetName() string
	GetConn() ziface.IConnection
	GetCurrentGameID() int
	GetLevel() int
	IsInGame() bool
	GetStatus() PlayerStatus

	// 阵容管理
	GetUnits() []IUnit           // 玩家持有的战斗单位
	AddUnit(unit IUnit) error    // 添加单位到阵容
	RemoveUnit(unitID int) error // 从阵容移除单位

	// 资源系统
	GetGold() int       // 金币
	AddGold(amount int) // 增加金币
	//GetItems() []item.IItem           // 道具栏

	// 输入处理
	HandleStartTurn() error
	HandleMoveRequest(unitID int, target interface{}) error // 处理移动请求
	HandleAttackRequest(unitID, targetID int) error         // 处理攻击请求
	HandleSkillRequest(unitID, targetID, skillID int) error // 处理技能请求
	HandleEndTurn() error

	// 同步与通知
	SyncGameState(state interface{}) error           // 同步游戏状态到客户端
	Notify(eventType string, data interface{}) error // 发送事件通知

	// 生命周期
	OnDisconnect()                             // 断线处理
	OnReconnect(conn ziface.IConnection) error // 重连处理
	Save() error                               // 保存玩家数据到数据库
}

// PlayerStatus 玩家状态枚举
type PlayerStatus int

const (
	StateIdle           PlayerStatus = iota // 空闲状态
	StateMatching                           // 匹配中
	StateBattleWaiting                      // 战斗中-回合等待
	StateBattleActing                       // 战斗中-行动中
	StateBattleFinished                     // 战斗结束
)
