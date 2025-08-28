package status

import (
	"mmo_game/internal/icore/status"
)

// ========== Attributes 实现 ==========

type Attributes struct {
	HP        int // 生命值
	Attack    int // 攻击力
	Defense   int // 防御力
	Speed     int // 速度
	Range     int // 攻击范围
	MoveRange int // 移动范围
	Movement  int //移动力
}

func (a *Attributes) GetBaseMovement() int {
	return a.Movement
}

func (a *Attributes) GetHP() int {
	return a.HP
}

func (a *Attributes) GetAttack() int {
	return a.Attack
}

func (a *Attributes) GetDefense() int {
	return a.Defense
}

func (a *Attributes) GetSpeed() int {
	return a.Speed
}

func (a *Attributes) GetRange() int {
	return a.Range
}

func (a *Attributes) GetMoveRange() int {
	return a.MoveRange
}

func (a *Attributes) Clone() status.IAttributes {
	return &Attributes{
		HP:        a.HP,
		Attack:    a.Attack,
		Defense:   a.Defense,
		Speed:     a.Speed,
		Range:     a.Range,
		MoveRange: a.MoveRange,
		Movement:  a.Movement,
	}
}
