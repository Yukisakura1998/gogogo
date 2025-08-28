package factory

import (
	"fmt"
	"mmo_game/internal/icore/entity"
	"mmo_game/internal/icore/factory"
	"mmo_game/internal/icore/status"
)

type GlobalUnitFactory struct {
	factories map[string]factory.IUnitFactory
}

func NewGlobalUnitFactory() *GlobalUnitFactory {
	return &GlobalUnitFactory{factories: make(map[string]factory.IUnitFactory)}
}

func (g *GlobalUnitFactory) Register(key string, factory factory.IUnitFactory) {
	g.factories[key] = factory
}

func (g *GlobalUnitFactory) CreateUnit(key, id, teamID string, pos status.IPosition) (entity.IUnit, error) {
	if factory, ok := g.factories[key]; ok {
		return factory.CreateUnit(id, teamID, pos), nil
	}
	return nil, fmt.Errorf("factory not found for key: %s", key)
}
