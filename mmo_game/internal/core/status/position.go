package status

import (
	"mmo_game/internal/icore/status"
)

// Position 表示二维坐标位置
type Position struct {
	X int // 水平坐标
	Y int // 垂直坐标
}

// NewPosition 创建新位置实例
func NewPosition(x, y int) *Position {
	return &Position{X: x, Y: y}
}

// GetX 实现IPosition接口
func (p *Position) GetX() int {
	return p.X
}

// GetY 实现IPosition接口
func (p *Position) GetY() int {
	return p.Y
}

// DistanceTo 计算曼哈顿距离（棋盘距离）
func (p *Position) DistanceTo(target status.IPosition) int {
	dx := p.X - target.GetX()
	dy := p.Y - target.GetY()
	return abs(dx) + abs(dy)
}

// IsAdjacent 判断是否相邻（包括斜向相邻）
func (p *Position) IsAdjacent(target status.IPosition) bool {
	return p.DistanceTo(target) == 1
}

func (p *Position) SetPosition(pos status.IPosition) {
	p.X = pos.GetX()
	p.Y = pos.GetY()
}

// 辅助函数：绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
