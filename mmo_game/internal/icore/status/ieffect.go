package status

import "mmo_game/internal/icore/entity"

// IEffect 定义状态效果接口
type IEffect interface {
	GetID() string
	GetType() string
	GetDuration() int
	GetModifier() int
	Tick() bool                    // 返回是否应该移除该效果
	OnTurnStart(unit entity.IUnit) //回合开始时该状态的事件
	OnTurnEnd(unit entity.IUnit)   //回合开始时该状态的事件
	IsExpired() bool               //状态结束后的操作
	Clone() IEffect                //复制状态
}
