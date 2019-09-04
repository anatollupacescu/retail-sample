//go:generate mockgen -source=repository.go -package mocks -destination mocks/repository.go
package itemtype

type (

	Store interface {
		Add(string) DTO
		Get(uint64) DTO
		Remove(uint64)
		List() []DTO
	}

	Repository interface {
		Add(string) uint64
		Get(uint64) string
		Remove(uint64)
		List() []string
	}

	DTO struct {
		Id uint64
		Name string
	}

	repository struct {
		store Store
	}
)

func (r *repository) List() []string {
	var v []string
	for _, i := range r.store.List() {
		v = append(v, i.Name)
	}
	return v
}

func (r *repository) Add(name string) uint64 {
	return r.store.Add(name).Id
}

func (r *repository) Remove(id uint64) {
	r.store.Remove(id)
}

func (r *repository) Get(i uint64) string {
	return r.store.Get(i).Name
}
