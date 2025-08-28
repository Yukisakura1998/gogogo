package status

import (
	"mmo_game/internal/icore/entity"
	"mmo_game/internal/icore/status"
)

// ========== StatusEffect 实现 ==========

type Effect struct {
	ID       string
	Type     string
	Duration int
	Modifier int
}

func (e *Effect) OnTurnStart(unit entity.IUnit) {
	//TODO implement me
	panic("implement me")
}

func (e *Effect) OnTurnEnd(unit entity.IUnit) {
	//TODO implement me
	panic("implement me")
}

func (e *Effect) IsExpired() bool {
	//TODO implement me
	panic("implement me")
}

func (e *Effect) Clone() status.IEffect {
	//TODO implement me
	panic("implement me")
}

func (e *Effect) GetID() string {
	return e.ID
}

func (e *Effect) GetType() string {
	return e.Type
}

func (e *Effect) GetDuration() int {
	return e.Duration
}

func (e *Effect) GetModifier() int {
	return e.Modifier
}

func (e *Effect) Tick() bool {
	e.Duration--
	return e.Duration <= 0
}

func GetMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
