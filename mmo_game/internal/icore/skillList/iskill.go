package skillList

// ISkill 定义技能接口
type ISkill interface {
	// 基础信息
	GetID() string   // 技能唯一标识
	GetName() string // 技能名称

	// 战斗属性
	GetDamage() int // 基础伤害
	GetRange() int  // 技能范围

	Clone() ISkill //复制
}
