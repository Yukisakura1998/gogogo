package entity

import (
	"encoding/json"
	"fmt"
	"mmo_game/internal/icore/entity"
	"sync"
	"zinx/ziface"
)

type Player struct {
	id            int                //该玩家的ID
	name          string             //该玩家的名字
	level         int                //该玩家的等级
	conn          ziface.IConnection //该玩家的连接
	status        entity.PlayerState //该玩家的状态
	currentGameID int                //当前游戏房间ID
	gold          int                //金币
	mutex         sync.RWMutex       //并发控制
}

func NewPlayer(id int) *Player {
	return &Player{id: id}
}

func (p *Player) SetPlayerID(i int) {
	p.id = i
}

func (p *Player) SetName(s string) {
	p.name = s
}

func (p *Player) SetLevel(i int) {
	p.level = i
}

func (p *Player) SetConn(connection ziface.IConnection) {
	p.conn = connection
}

func (p *Player) SetStatus(state entity.PlayerState) {
	p.status = state
}

func (p *Player) SetCurrentGameID(i int) {
	p.currentGameID = i
}

func (p *Player) GetPlayerID() int {
	return p.id
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetLevel() int {
	return p.level
}

func (p *Player) GetConn() ziface.IConnection {
	return p.conn
}

func (p *Player) GetStatus() entity.PlayerState {
	return p.status
}

func (p *Player) GetCurrentGameID() int {
	return p.currentGameID
}

func (p *Player) GetGold() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.gold
}

func (p *Player) AddGold(amount int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.gold += amount
}

func (p *Player) ReduceGold(amount int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.gold -= amount
}

func (p *Player) Connect() error {
	if p.conn == nil {
		return fmt.Errorf("player %d connection is nil", p.id)
	}

	// 设置连接关闭时的回调（自动触发玩家下线）
	p.conn.SetOnConnStop(func(conn ziface.IConnection) {
		err := p.Disconnect()
		if err != nil {
			return
		}
	})

	fmt.Printf("player %d connected successfully\n", p.id)
	return nil
}

func (p *Player) Disconnect() error {
	if p.GetStatus() == entity.PlayerStateOffline {
		return fmt.Errorf("player %d already offline", p.id)
	}

	// 执行登出逻辑
	p.Logout()
	// 清理连接引用
	p.conn = nil
	fmt.Printf("player %d disconnected successfully\n", p.id)
	return nil
}

func (p *Player) SendMsg(msgID uint32, data []byte) error {
	if p.conn == nil {
		return fmt.Errorf("player %d connection is nil\n", p.id) // 增加玩家ID便于定位
	}
	if err := p.conn.SendMsg(msgID, data); err != nil {
		return fmt.Errorf("player %d send msg failed: %w\n", p.id, err) // 包装错误
	}
	return nil
}

func (p *Player) Login() bool {
	// 发送玩家ID同步消息
	syncPid := struct {
		Pid int `json:"pid"` // 玩家ID，JSON序列化标签
	}{
		Pid: p.id,
	}
	data, err := json.Marshal(syncPid)
	if err != nil {
		fmt.Printf("marshal SyncPid error: %v\n", err)
		return false
	}

	// （假设消息ID 1为登录，与客户端保持一致）
	if err := p.SendMsg(entity.MsgIDLogin, data); err != nil {
		fmt.Printf("player send SyncPid error: %v\n", err)
		return false
	}

	fmt.Printf("player %d login successfully\n", p.id)
	p.status = entity.PlayerStateIdle
	return true
}

func (p *Player) Logout() {
	fmt.Printf("[系统]玩家id:%d，已下线\n", p.id)
	// 登出逻辑：设置状态为离线并断开连接
	p.SetStatus(entity.PlayerStateOffline)
	if p.conn != nil {
		p.conn.Stop()
	}
}

func (p *Player) Talk(content string) {
	// 聊天逻辑：构建消息并发送
	if p.GetStatus() == entity.PlayerStateOffline {
		return
	}

	// 构建聊天消息结构
	talkMsg := struct {
		PlayerID int    `json:"player_id"`
		Name     string `json:"name"`
		Content  string `json:"content"`
	}{
		PlayerID: p.GetPlayerID(),
		Name:     p.GetName(),
		Content:  content,
	}
	// 序列化为JSON字节数组
	data, err := json.Marshal(talkMsg)
	if err != nil {
		fmt.Printf("玩家[%d]聊天消息JSON序列化失败: %v\n", p.id, err)
		return
	}

	// 发送聊天消息（假设消息ID 3为聊天消息，与客户端保持一致）
	if err := p.SendMsg(entity.MsgIDTalk, data); err != nil {
		fmt.Printf("玩家[%d]发送聊天消息失败: %v\n", p.id, err)
	}

}
