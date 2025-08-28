package skillList

import (
	"mmo_game/internal/icore/skillList"
)

// ========== Skill 实现 ==========

type Skill struct {
	ID     string
	Name   string
	Damage int
	Range  int
}

func (s *Skill) Clone() skillList.ISkill {
	return &Skill{
		ID:     s.ID,
		Name:   s.Name,
		Damage: s.Damage,
		Range:  s.Range,
	}
}

func (s *Skill) GetID() string {
	return s.ID
}

func (s *Skill) GetName() string {
	return s.Name
}

func (s *Skill) GetDamage() int {
	return s.Damage
}

func (s *Skill) GetRange() int {
	return s.Range
}
