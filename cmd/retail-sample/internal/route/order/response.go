package order

import (
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
)

type (
	entity struct {
		ID       int `json:"id"`
		RecipeID int `json:"recipeID"`
		Qty      int `json:"qty"`
	}

	single struct {
		Data entity `json:"data"`
	}

	collection struct {
		Data []entity `json:"data"`
	}
)

func toSingleResponse(o order.DTO) single {
	return single{
		Data: entity{
			ID:       o.ID,
			RecipeID: o.RecipeID,
			Qty:      o.Qty,
		},
	}
}

func toCollectionResponse(orders []order.DTO) (response collection) {
	response.Data = make([]entity, 0, len(orders))

	for i := range orders {
		currentOrder := orders[i]

		e := entity{
			ID:       currentOrder.ID,
			RecipeID: currentOrder.RecipeID,
			Qty:      currentOrder.Qty,
		}

		response.Data = append(response.Data, e)
	}

	return response
}

func httpServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	statusText := http.StatusText(status)
	http.Error(w, statusText, http.StatusInternalServerError)
}
