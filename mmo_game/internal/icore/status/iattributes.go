package status

// IAttributes 定义单位基础属性接口
type IAttributes interface {
	//基础属性获取
	GetHP() int           // 生命值
	GetAttack() int       // 攻击力
	GetDefense() int      // 防御力
	GetSpeed() int        // 速度（影响行动顺序）
	GetRange() int        // 攻击范围
	GetMoveRange() int    // 移动范围
	GetBaseMovement() int //基础移动力，在定义时赋值且不更改

	// Clone 克隆方法（用于战斗计算等场景）
	Clone() IAttributes
}
