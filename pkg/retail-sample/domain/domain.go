package domain

type (
	sellableItem struct {
		designID int
	}

	designEntity struct {
		designs []Design
	}
)

func (d *designEntity) All() []Design {
	return d.designs
}

func (d *designEntity) ForId(designID int) Design {
	return Design{
		designID: designID,
	}
}
