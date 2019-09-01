package itemtype

type (
	Repository map[ItemType]bool

	ItemType struct {
		Name string
		Qty  int
	}
)

func (repository *Repository) ListItemTypes() []ItemType {
	types := make([]ItemType, 0, len(*repository))
	for t := range *repository {
		types = append(types, t)
	}
	return types
}

func (repository *Repository) AddItemType(s string, i int) {
	(*repository)[ItemType{Name: s, Qty: i}] = true
}

func (repository *Repository) RemoveItemType(s string, i int) {
	delete(*repository, ItemType{
		Name: s,
		Qty:  i,
	})
}
