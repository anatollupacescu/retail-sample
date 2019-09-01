package itemtype

type (
	ItemTypeDB interface {
		Add(string) uint64
		Get(uint64) Entity
		Remove(uint64)
		List() []Entity
	}

	repository struct {
		db ItemTypeDB
	}

	ItemType struct {
		Name string
	}
)

func (r *repository) List() []ItemType {
	var v []ItemType
	for _, i := range r.db.List() {
		v = append(v, ItemType{
			Name: i.name,
		})
	}
	return v
}

func (r *repository) Add(name string) uint64 {
	return r.db.Add(name)
}

func (r *repository) Remove(id uint64) {
	r.db.Remove(id)
}

func (r *repository) Get(i uint64) ItemType {
	entity := r.db.Get(i)
	return ItemType{
		Name: entity.name,
	}
}
