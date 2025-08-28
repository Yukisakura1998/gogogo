package status

// IPosition 定义单位坐标接口
type IPosition interface {
	// GetX 获取X坐标（水平方向）
	GetX() int

	// GetY 获取Y坐标（垂直方向）
	GetY() int

	// DistanceTo 计算到目标位置的距离（棋盘距离）
	DistanceTo(target IPosition) int

	// IsAdjacent 判断是否与目标位置相邻（距离为1）
	IsAdjacent(target IPosition) bool

	//SetPosition 设置位置
	SetPosition(pos IPosition)
}
