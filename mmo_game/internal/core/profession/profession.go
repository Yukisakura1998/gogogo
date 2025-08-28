package profession

type Profession struct {
	ID        string
	Name      string
	SkillPool []string
}

func NewProfession(id, name string, skillPool []string) *Profession {
	return &Profession{
		ID:        id,
		Name:      name,
		SkillPool: skillPool,
	}
}

func (p *Profession) GetID() string {
	return p.ID
}

func (p *Profession) GetName() string {
	return p.Name
}

func (p *Profession) GetSkillPool() []string {
	return p.SkillPool
}
