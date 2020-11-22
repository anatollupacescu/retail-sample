package order

import (
	"encoding/json"
	"net/http"

	usecase2 "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"

	"github.com/rs/zerolog/hlog"
)

type createPayload struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

func toCreateDTO(r *http.Request) (usecase2.PlaceOrderDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var payload createPayload

	if err := d.Decode(&payload); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'create order' payload")
		return usecase2.PlaceOrderDTO{}, err
	}

	dto := usecase2.PlaceOrderDTO{
		RecipeID: payload.ID,
		OrderQty: payload.Qty,
	}

	return dto, nil
}
