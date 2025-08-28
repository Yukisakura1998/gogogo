package factory

import (
	"mmo_game/internal/icore/entity"
	"mmo_game/internal/icore/status"
)

type IUnitFactory interface {
	CreateUnit(id, teamID string, pos status.IPosition) entity.IUnit
}
