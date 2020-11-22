package order

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/hlog"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/order"
)

type createPayload struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

func toCreateDTO(r *http.Request) (usecase.PlaceOrderDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var payload createPayload

	if err := d.Decode(&payload); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'create order' payload")
		return usecase.PlaceOrderDTO{}, err
	}

	dto := usecase.PlaceOrderDTO{
		RecipeID: payload.ID,
		OrderQty: payload.Qty,
	}

	return dto, nil
}
