package entity

import "mmo_game/internal/core/skillList"

var SkillLibrary = map[string]skillList.Skill{
	"slash": {
		ID:     "slash",
		Name:   "劈砍",
		Damage: 50,
		Range:  1,
	},
	"shield_bash": {
		ID:     "shield_bash",
		Name:   "盾击",
		Damage: 30,
		Range:  2,
	},
}
