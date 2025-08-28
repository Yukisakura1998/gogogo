package entity

import (
	"container/list"
	"encoding/json"
	valueobj "mmo_game/internal/core/status"
	"mmo_game/internal/icore/entity"
	"mmo_game/internal/icore/profession"
	"mmo_game/internal/icore/skillList"
	"mmo_game/internal/icore/status"
)

type Unit struct {
	ID           int
	Type         string
	HealthPoint  int
	MagicPoint   int
	Attack       int
	Defence      int
	Movement     int
	OwnerID      int
	Position     status.IPosition
	CanMove      bool
	Profession   profession.IProfession
	Attributes   status.IAttributes
	skillList    []skillList.ISkill
	IsAlive      bool
	StatusEffect []status.IEffect
}

func (u *Unit) GetUnitID() int {
	return u.ID
}

func (u *Unit) GetPosition() status.IPosition {
	return u.Position
}

func (u *Unit) GetUnitType() string {
	return u.Type
}

func (u *Unit) GetTeamID() int {
	return u.OwnerID
}

func (u *Unit) GetProfession() profession.IProfession {
	return u.Profession
}

func (u *Unit) GetAttributes() status.IAttributes {
	return u.Attributes
}

func (u *Unit) GetSkills() []skillList.ISkill {
	return u.skillList
}

func (u *Unit) TakeDamage(damage int) {
	u.HealthPoint -= damage
	if u.HealthPoint <= 0 {
		u.IsAlive = false
		u.CanMove = false
		u.RemoveFromMap()
	}
}

func (u *Unit) AttackToBySkill(target entity.IUnit, skill skillList.ISkill) (bool, error) {
	damage := skill.GetDamage()
	target.TakeDamage(damage)
	//TODO:AttackToBySkill() 后续添加命中判断
	return true, nil
}

func (u *Unit) MoveTo(x, y int) {
	u.Position.SetPosition(valueobj.NewPosition(x, y))
}

func (u *Unit) CanMoveTo(pos status.IPosition) bool {
	maxSteps := u.Movement
	// 定义四方向
	directions := []struct{ DX, DY int }{
		{0, 1}, {1, 0}, {0, -1}, {-1, 0},
	}

	visited := make(map[valueobj.Position]bool)
	queue := list.New()

	// 从 IPosition 获取坐标
	startPos := valueobj.Position{
		X: u.Position.GetX(),
		Y: u.Position.GetY(),
	}

	queue.PushBack(struct {
		Pos   valueobj.Position
		Steps int
	}{
		Pos:   startPos,
		Steps: maxSteps,
	})

	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)
		current := element.Value.(struct {
			Pos   valueobj.Position
			Steps int
		})

		if current.Pos.X == pos.GetX() && current.Pos.Y == pos.GetY() {
			return true
		}

		if current.Steps <= 0 {
			continue
		}

		visited[current.Pos] = true

		for _, dir := range directions {
			nextPos := valueobj.Position{
				X: current.Pos.X + dir.DX,
				Y: current.Pos.Y + dir.DY,
			}

			if visited[nextPos] {
				continue
			}

			queue.PushBack(struct {
				Pos   valueobj.Position
				Steps int
			}{
				Pos:   nextPos,
				Steps: current.Steps - 1,
			})
		}
	}

	return false
}

func (u *Unit) RemoveFromMap() bool {
	//TODO:RemoveFromMap() 在二维地图中移除本单位
	//在地图上移除本单位
	return true
}

func (u *Unit) GetActionPoints() int {
	return u.Movement
}

func (u *Unit) SpendActionPoints(cost int) bool {
	if u.Movement >= cost {
		u.Movement -= cost
		return true
	}
	return false
}

func (u *Unit) AddStatusEffect(effect status.IEffect) {
	u.StatusEffect = append(u.StatusEffect, effect)
}

func (u *Unit) RemoveStatusEffect(effectID string) {
	for i, effect := range u.StatusEffect {
		if effect.GetID() == effectID {
			u.StatusEffect = append(u.StatusEffect[:i], u.StatusEffect[i+1:]...)
			return
		}
	}
}

func (u *Unit) OnTurnStart() {
	// 重置移动点数
	u.Movement = u.Attributes.GetBaseMovement()
	u.CanMove = true

	// 更新状态效果
	for i := 0; i < len(u.StatusEffect); {
		effect := u.StatusEffect[i]
		effect.OnTurnStart(u)
		if effect.IsExpired() {
			u.RemoveStatusEffect(effect.GetID())
		} else {
			i++
		}
	}
}

func (u *Unit) OnTurnEnd() {
	// 更新状态效果
	for i := 0; i < len(u.StatusEffect); {
		effect := u.StatusEffect[i]
		effect.OnTurnEnd(u)
		if effect.IsExpired() {
			u.RemoveStatusEffect(effect.GetID())
		} else {
			i++
		}
	}
}

func (u *Unit) Clone() entity.IUnit {
	// 深拷贝位置
	pos := &valueobj.Position{
		X: u.Position.GetX(),
		Y: u.Position.GetY(),
	}

	// 深拷贝状态效果
	statusEffects := make([]status.IEffect, len(u.StatusEffect))
	for i, effect := range u.StatusEffect {
		statusEffects[i] = effect.Clone()
	}

	// 深拷贝技能列表
	skills := make([]skillList.ISkill, len(u.skillList))
	for i, skill := range u.skillList {
		skills[i] = skill.Clone()
	}

	return &Unit{
		ID:           u.ID,
		Type:         u.Type,
		HealthPoint:  u.HealthPoint,
		MagicPoint:   u.MagicPoint,
		Attack:       u.Attack,
		Defence:      u.Defence,
		Movement:     u.Movement,
		OwnerID:      u.OwnerID,
		Position:     pos,
		CanMove:      u.CanMove,
		Profession:   u.Profession, // 注意: 这里假设Profession是不可变的或已处理深拷贝
		Attributes:   u.Attributes, // 注意: 这里假设Attributes是不可变的或已处理深拷贝
		skillList:    skills,
		IsAlive:      u.IsAlive,
		StatusEffect: statusEffects,
	}
}

func (u *Unit) Serialize() ([]byte, error) {
	// 简单实现: 使用JSON序列化
	data := struct {
		ID          int
		Type        string
		HealthPoint int
		MagicPoint  int
		Attack      int
		Defence     int
		Movement    int
		OwnerID     int
		Position    *valueobj.Position
		CanMove     bool
		IsAlive     bool
		// 注意: 没有序列化Profession, Attributes, skillList和StatusEffect
		// 因为它们可能是接口类型，需要更复杂的序列化逻辑
	}{
		ID:          u.ID,
		Type:        u.Type,
		HealthPoint: u.HealthPoint,
		MagicPoint:  u.MagicPoint,
		Attack:      u.Attack,
		Defence:     u.Defence,
		Movement:    u.Movement,
		OwnerID:     u.OwnerID,
		Position: &valueobj.Position{
			X: u.Position.GetX(),
			Y: u.Position.GetY(),
		},
		CanMove: u.CanMove,
		IsAlive: u.IsAlive,
	}

	return json.Marshal(data)
}

func (u *Unit) Deserialize(data []byte) error {
	var decoded struct {
		ID          int
		Type        string
		HealthPoint int
		MagicPoint  int
		Attack      int
		Defence     int
		Movement    int
		OwnerID     int
		Position    *valueobj.Position
		CanMove     bool
		IsAlive     bool
	}

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	u.ID = decoded.ID
	u.Type = decoded.Type
	u.HealthPoint = decoded.HealthPoint
	u.MagicPoint = decoded.MagicPoint
	u.Attack = decoded.Attack
	u.Defence = decoded.Defence
	u.Movement = decoded.Movement
	u.OwnerID = decoded.OwnerID
	u.Position = decoded.Position
	u.CanMove = decoded.CanMove
	u.IsAlive = decoded.IsAlive

	return nil
}

func NewUnit(ID int, Type string, healthPoint int, magicPoint int, attack int, defence int, movement int, ownerID int, position status.IPosition, canMove bool) *Unit {
	return &Unit{
		ID:          ID,
		Type:        Type,
		HealthPoint: healthPoint,
		MagicPoint:  magicPoint,
		Attack:      attack,
		Defence:     defence,
		Movement:    movement,
		OwnerID:     ownerID,
		Position:    position,
		CanMove:     canMove,
	}
}
