package temp

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int
	MinX      int
	MaxX      int
	MinY      int
	MaxY      int
	playerIDs map[int]bool
	pIDLock   sync.RWMutex
}

func (g *Grid) Add(id int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[id] = true
}

func (g *Grid) Remove(id int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, id)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	for i := range g.playerIDs {
		playerIDs = append(playerIDs, i)
	}
	return
}
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d, MinX:%d, MinY:%d, MaxX:%d, MaxY:%d, playerIDs:%v", g.GID, g.MinX, g.MinY, g.MaxX, g.MaxY, g.playerIDs)
}

func NewGrid(GID int, minX int, maxX int, minY int, maxY int) *Grid {
	return &Grid{
		GID:       GID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}
