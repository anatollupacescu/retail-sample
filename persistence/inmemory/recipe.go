package inmemory

import domain "github.com/anatollupacescu/retail-sample/domain/retail-sample/recipe"

type Recipe struct {
	data    map[int]domain.Recipe
	counter *int
}

func NewRecipe() Recipe {
	zero := 0

	return Recipe{
		data:    make(map[int]domain.Recipe),
		counter: &zero,
	}
}

func (m *Recipe) Add(r domain.Recipe) (domain.ID, error) {
	*m.counter++
	id := *m.counter
	m.data[id] = r

	return domain.ID(id), nil
}

func (m *Recipe) List() (r []domain.Recipe, err error) {
	for id := range m.data {
		rp := m.data[id]
		rp.ID = domain.ID(id)
		r = append(r, rp)
	}

	return
}

func (m *Recipe) Get(id domain.ID) (domain.Recipe, error) {
	dict := m.data

	if val, ok := dict[int(id)]; ok {
		val.ID = id
		return val, nil
	}

	return domain.Recipe{}, domain.ErrRecipeNotFound
}

func (m *Recipe) Save(r domain.Recipe) error {
	dict := m.data

	if _, ok := dict[int(r.ID)]; !ok {
		return domain.ErrRecipeNotFound
	}

	m.data[int(r.ID)] = r

	return nil
}
