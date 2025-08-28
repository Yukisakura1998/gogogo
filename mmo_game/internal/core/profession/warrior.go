package profession

import (
	"mmo_game/internal/core/entity"
	"mmo_game/internal/icore/skillList"
)

type Warrior struct {
	Profession
}

func (w *Warrior) GetID() string {
	return "Warrior"
}

func (w *Warrior) GetName() string {
	return "战士"
}

func (w *Warrior) GetSkillPool() []string {
	return []string{"slash", "shield_bash"}
}

func (w *Warrior) GetSkills() []skillList.ISkill {
	skills := make([]skillList.ISkill, 0)
	for _, skillID := range w.GetSkillPool() {
		if skill, ok := entity.SkillLibrary[skillID]; ok {
			skills = append(skills, skill.Clone()) // 克隆技能实例
		}
	}
	return skills
}
