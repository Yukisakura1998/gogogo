package temp

import "fmt"

type AOIManager struct {
	MinX    int
	MaxX    int
	CountsX int
	MinY    int
	MaxY    int
	CountsY int
	grids   map[int]*Grid
}

func NewAOIManager(minX int, maxX int, countsX int, minY int, maxY int, countsY int) *AOIManager {
	m := &AOIManager{
		MinX:    minX,
		MaxX:    maxX,
		CountsX: countsX,
		MinY:    minY,
		MaxY:    maxY,
		CountsY: countsY,
		grids:   make(map[int]*Grid),
	}
	for y := 0; y < countsY; y++ {
		for x := 0; x < countsX; x++ {
			gid := y*countsX + x
			m.grids[gid] = NewGrid(gid, minX+x*m.gridWidth(), minX+(x+1)*m.gridWidth(), minY+y*m.gridLength(), minY+(y+1)*m.gridLength())
		}
	}
	return m
}

func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CountsX
}
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CountsY
}
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager MinX:%d, MinY:%d, MaxX:%d, MaxY:%d, CountsX:%d, CountsY:%d\n", m.MinX, m.MinY, m.MaxX, m.MaxY, m.CountsX, m.CountsY)
	for i := range m.grids {
		s += fmt.Sprintln(i)
	}
	return s
}

func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {

	if _, ok := m.grids[gID]; ok {
		return
	}

	grids = append(grids, m.grids[gID])

	idx := gID % m.CountsX
	//左边
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	//右边
	if idx < m.CountsX-1 {
		grids = append(grids, m.grids[gID+1])
	}
	gridsX := make([]int, 0, len(grids))

	for _, v := range grids {
		gridsX = append(gridsX, v.GID)
	}

	for _, v := range gridsX {
		idy := v % m.CountsY
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CountsX])
		}
		//右边
		if idy < m.CountsY-1 {
			grids = append(grids, m.grids[v+m.CountsY])
		}
	}
	return
}

func (m *AOIManager) GetGIDByPosition(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()
	return gy*m.CountsX + gx
}

func (m *AOIManager) GetPIDsByPosition(x, y float32) (pIDs []int) {
	gID := m.GetGIDByPosition(x, y)
	grid := m.GetSurroundGridsByGid(gID)

	for _, i := range grid {
		pIDs = append(pIDs, i.GetPlayerIDs()...)
		fmt.Printf("====>grid ID:%d,playerID:%v", i.GID, i.GetPlayerIDs())
	}
	return
}

func (m *AOIManager) GetPlayersByGID(gID int) (pIDs []int) {
	pIDs = m.grids[gID].GetPlayerIDs()
	return
}

// RemovePlayersFromGID 根据ID移除一个格子中的player
func (m *AOIManager) RemovePlayersFromGID(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// AddPlayersToGrid 添加一个player到指定的格子中
func (m *AOIManager) AddPlayersToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// AddPlayersByPositionToGrid 根据坐标添加一个player到指定的格子中
func (m *AOIManager) AddPlayersByPositionToGrid(pID int, x, y float32) {
	gID := m.GetGIDByPosition(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

// RemovePlayersByPositionToGrid 根据坐标删除一个player到指定的格子中
func (m *AOIManager) RemovePlayersByPositionToGrid(pID int, x, y float32) {
	gID := m.GetGIDByPosition(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)
}
