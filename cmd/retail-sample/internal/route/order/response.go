package order

import "github.com/anatollupacescu/retail-sample/domain/retail/order"

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

func toSingleResponse(o order.Order) single {
	return single{
		Data: entity{
			ID:       int(o.ID),
			RecipeID: o.Entry.RecipeID,
			Qty:      o.Entry.Qty,
		},
	}
}

func toCollectionResponse(orders []order.Order) (response collection) {
	response.Data = make([]entity, 0, len(orders))

	for i := range orders {
		currentOrder := orders[i]

		e := entity{
			ID:       int(currentOrder.ID),
			RecipeID: currentOrder.RecipeID,
			Qty:      currentOrder.Qty,
		}

		response.Data = append(response.Data, e)
	}

	return response
}
