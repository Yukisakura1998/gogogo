package entity

import (
	"errors"
	"mmo_game/internal/core/status"
	"mmo_game/internal/icore/entity"
	"sync"
	"zinx/ziface"
)

type Player struct {
	ID          int                //该玩家的ID
	Name        string             //该玩家的名字
	Level       int                //该玩家的等级
	Conn        ziface.IConnection //该玩家的连接
	InGame      bool               //是否正在进行游戏
	Status      entity.PlayerStatus
	CurrentGame string // 当前游戏房间ID
	// 游戏数据
	Units []*Unit      // 战斗单位阵容
	Gold  int          // 金币
	mutex sync.RWMutex // 并发控制
}

func (p *Player) GetName() string {
	return p.Name
}
func (p *Player) GetLevel() int {
	return p.Level
}

func (p *Player) GetConn() ziface.IConnection {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Conn
}

func (p *Player) GetCurrentGameID() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.InGame && p.CurrentGame != "" {
		//TODO GetCurrentGameID(): 简单返回1，实际项目应返回具体游戏ID
		return 1
	}
	return 0
}

func (p *Player) IsInGame() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.InGame
}

func (p *Player) GetStatus() entity.PlayerStatus {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Status
}

func (p *Player) GetUnits() []entity.IUnit {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	units := make([]entity.IUnit, len(p.Units))
	for i, u := range p.Units {
		units[i] = u
	}
	return units
}

func (p *Player) AddUnit(unit entity.IUnit) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.Units) >= 5 { // TODO AddUnit():根据不同地图修改阵容上限
		return errors.New("unit slot full")
	}

	// 类型断言确保是*Unit类型
	if u, ok := unit.(*Unit); ok {
		p.Units = append(p.Units, u)
		return nil
	}
	return errors.New("invalid unit type")
}

func (p *Player) RemoveUnit(unitID int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i, u := range p.Units {
		if u.GetUnitID() == unitID {
			p.Units = append(p.Units[:i], p.Units[i+1:]...)
			return nil
		}
	}
	return errors.New("unit not found")
}

func (p *Player) GetGold() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Gold
}

func (p *Player) AddGold(amount int) {
	if amount == 0 {
		return
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Gold += amount
}

func (p *Player) HandleMoveRequest(unitID int, target interface{}) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// 检查状态
	if p.Status != entity.StateBattleActing {
		return errors.New("not your turn")
	}

	// 转换目标坐标
	pos, ok := target.(status.Position)
	if !ok {
		return errors.New("invalid target coordinate")
	}

	// 查找单位
	var unit *Unit
	for _, u := range p.Units {
		if u.GetUnitID() == unitID {
			unit = u
			break
		}
	}
	if unit == nil {
		return errors.New("unit not found")
	}

	// 实际项目中这里应该调用GameManager验证移动范围
	// 这里仅作示例，假设移动总是合法的
	unit.Position.SetPosition(&pos)
	return nil
}

func (p *Player) HandleAttackRequest(unitID, targetID int) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.Status != entity.StateBattleActing {
		return errors.New("not your turn")
	}

	// 查找攻击单位和目标单位
	var attacker, target *Unit
	for _, u := range p.Units {
		if u.GetUnitID() == unitID {
			attacker = u
			break
		}
	}
	if attacker == nil {
		return errors.New("attacker unit not found")
	}

	// 实际项目中需要通过GameManager查找目标单位
	// 这里简化处理
	_ = targetID

	/*示例：简单扣血
	_, err := attacker.AttackToBySkill(target, attacker.skillList[0])
	if err != nil {
		return err
	}

	*/
	target.TakeDamage(attacker.Attack)
	return nil
}

func (p *Player) HandleSkillRequest(unitID, targetID, skillID int) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.Status != entity.StateBattleActing {
		return errors.New("not your turn")
	}

	// 查找攻击单位和目标单位
	var attacker, target *Unit
	for _, u := range p.Units {
		if u.GetUnitID() == unitID {
			attacker = u
			break
		}
	}
	if attacker == nil {
		return errors.New("attacker unit not found")
	}

	/* 验证技能是否存在
	skill := unit.GetSkillByID(skillID)
	if skill == nil {
		return errors.New("skill not found")
	}*/

	// 实际项目中需要验证技能范围、消耗等
	// 这里简化处理
	_ = targetID

	// 示例：触发技能效果
	_, err := attacker.AttackToBySkill(target, attacker.skillList[skillID])
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) HandleStartTurn() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.Status != entity.StateBattleActing {
		return errors.New("not your turn")
	}

	// 重置单位行动状态
	for _, u := range p.Units {
		u.OnTurnStart()
	}

	p.Status = entity.StateBattleActing
	return nil
}

func (p *Player) HandleEndTurn() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.Status != entity.StateBattleActing {
		return errors.New("not your turn")
	}

	// 重置单位行动状态
	for _, u := range p.Units {
		u.OnTurnEnd()
	}

	p.Status = entity.StateBattleWaiting
	return nil
}

func (p *Player) SyncGameState(state interface{}) error {
	conn := p.GetConn()
	if conn == nil {
		return errors.New("connection lost")
	}

	// 实际项目中应该定义明确的协议格式
	msg := []byte("PingRouter...")

	// 这里假设有统一的消息发送方法
	return conn.SendMsg(0, msg)
}

func (p *Player) Notify(eventType string, data interface{}) error {
	conn := p.GetConn()
	if conn == nil {
		return errors.New("connection lost")
	}

	msg := []byte("PingRouter...")

	return conn.SendMsg(1, msg)
}

func (p *Player) OnDisconnect() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Conn = nil
	p.Status = entity.StateIdle
	p.InGame = false

	// 实际项目中应该通知GameManager处理玩家掉线
	// 例如：GameManager.OnPlayerDisconnect(p.id)
}

func (p *Player) OnReconnect(conn ziface.IConnection) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.Conn != nil {
		return errors.New("already connected")
	}

	p.Conn = conn
	p.Status = entity.StateBattleWaiting // 根据实际情况设置状态

	// 同步当前状态给客户端
	return p.SyncGameState(map[string]interface{}{
		"units": p.Units,
		"gold":  p.Gold,
	})
}

func (p *Player) Save() error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// 构建保存数据 data
	_ = map[string]interface{}{
		"id":    p.ID,
		"name":  p.Name,
		"level": p.Level,
		"gold":  p.Gold,
		"units": p.Units,
	}

	// 实际项目中应该调用数据库接口
	// 这里简化处理
	// return database.SavePlayerData(data)

	return nil
}

// NewPlayer 创建新玩家
func NewPlayer(ID int, name string, conn ziface.IConnection, inGame bool) *Player {
	return &Player{
		ID:     ID,
		Name:   name,
		Conn:   conn,
		InGame: inGame,
		Units:  make([]*Unit, 0),
		Status: entity.StateIdle,
	}
}

func (p *Player) GetPlayerID() int {
	return p.ID
}
