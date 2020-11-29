package inventory

import (
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

type (
	entity struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	single struct {
		Data entity `json:"data"`
	}
	collection struct {
		Data []entity `json:"data"`
	}
)

func toCollectionResponse(items []inventory.Item) collection {
	var response = collection{
		Data: make([]entity, 0, len(items)),
	}

	for _, item := range items {
		response.Data = append(response.Data, entity{
			ID:      item.ID,
			Name:    item.Name,
			Enabled: item.Enabled,
		})
	}

	return response
}

func toSingleResponse(i inventory.Item) single {
	return single{
		Data: entity{
			ID:      i.ID,
			Name:    i.Name,
			Enabled: i.Enabled,
		},
	}
}

func httpServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	statusText := http.StatusText(status)
	http.Error(w, statusText, http.StatusInternalServerError)
}
