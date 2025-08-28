package entity

import (
	"mmo_game/internal/icore/skillList"
)

type Type struct {
	TypeID              int
	TypeBaseHealthPoint int
	TypeBaseMagicPoint  int
	TypeBaseSkillList   []skillList.ISkill
}

func (t *Type) GetTypeID() int {
	return t.TypeID
}
