package entity

import (
	"mmo_game/internal/icore/profession"
	"mmo_game/internal/icore/skillList"
	"mmo_game/internal/icore/status"
)

type IUnit interface {
	GetUnitID() int                        //角色ID
	GetPosition() status.IPosition         //位置
	GetUnitType() string                   //类型
	GetTeamID() int                        //阵营
	GetProfession() profession.IProfession //职业

	GetAttributes() status.IAttributes                                  //属性
	GetSkills() []skillList.ISkill                                      //技能list
	TakeDamage(damage int)                                              //受到伤害
	AttackToBySkill(target IUnit, skill skillList.ISkill) (bool, error) // 指定技能攻击目标

	MoveTo(x, y int)                          //移动
	CanMoveTo(position status.IPosition) bool // 判断是否可以移动到指定位置
	RemoveFromMap() bool                      //从地图中移除
	GetActionPoints() int                     //获取移动点数
	SpendActionPoints(cost int) bool          //花费移动点数

	// ========== 状态效果 ==========
	AddStatusEffect(effect status.IEffect) // 添加状态（Buff，DeBuff）
	RemoveStatusEffect(effectID string)    //移除状态

	// ========== 回合事件 ==========
	OnTurnStart() //回合开始
	OnTurnEnd()   //回合结束

	// ========== 序列化 ==========
	Clone() IUnit                  //克隆
	Serialize() ([]byte, error)    //序列化
	Deserialize(data []byte) error //反序列化
}
