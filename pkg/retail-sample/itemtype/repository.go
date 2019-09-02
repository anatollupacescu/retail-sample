//go:generate mockgen -source=repository.go -package mocks -destination mocks/repository.go
package itemtype

type (

	Entity struct {
		Name string
	}

	ItemTypeStore interface {
		Add(string) uint64
		Get(uint64) Entity
		Remove(uint64)
		List() []Entity
	}

	Repository struct {
		DB ItemTypeStore
	}

	ItemType struct {
		Name string
	}
)

func (r *Repository) List() []ItemType {
	var v []ItemType
	for _, i := range r.DB.List() {
		v = append(v, ItemType{
			Name: i.Name,
		})
	}
	return v
}

func (r *Repository) Add(name string) uint64 {
	return r.DB.Add(name)
}

func (r *Repository) Remove(id uint64) {
	r.DB.Remove(id)
}

func (r *Repository) Get(i uint64) ItemType {
	entity := r.DB.Get(i)
	return ItemType{
		Name: entity.Name,
	}
}
