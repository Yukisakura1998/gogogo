package profession

type IProfession interface {
	GetID() string          // 获取职业ID
	GetName() string        // 获取职业名称
	GetSkillPool() []string // 获取职业技能池
	//GetPassiveEffects() []IStatusEffect // 获取职业被动效果
}
