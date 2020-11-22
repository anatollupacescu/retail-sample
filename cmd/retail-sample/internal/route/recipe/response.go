package recipe

import "github.com/anatollupacescu/retail-sample/domain/retail/recipe"

type (
	item struct {
		ID  int `json:"id"`
		Qty int `json:"qty"`
	}

	entity struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Items   []item `json:"items"`
		Enabled bool   `json:"enabled"`
	}

	single struct {
		Data entity `json:"data"`
	}

	collection struct {
		Data []entity `json:"data"`
	}
)

func toResponse(re recipe.Recipe) single {
	return single{
		Data: entity{
			ID:      int(re.ID),
			Name:    string(re.Name),
			Items:   toItems(re.Ingredients),
			Enabled: re.Enabled,
		},
	}
}

func toItems(i []recipe.Ingredient) (items []item) {
	for _, ri := range i {
		items = append(items, item{
			ID:  ri.ID,
			Qty: ri.Qty,
		})
	}

	return
}

func toCollectionResponse(all []recipe.Recipe) collection {
	var response = struct {
		Data []entity `json:"data"`
	}{
		Data: make([]entity, 0, len(all)),
	}

	for i := range all {
		r := all[i]

		response.Data = append(response.Data, entity{
			ID:      int(r.ID),
			Name:    string(r.Name),
			Items:   toItems(r.Ingredients),
			Enabled: r.Enabled,
		})
	}

	return response
}
