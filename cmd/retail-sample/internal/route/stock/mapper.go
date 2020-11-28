package stock

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"
)

var ErrParseItemID = errors.New("could not parse item ID")

func toProvisionDTO(r *http.Request) (usecase.ProvisionDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body provisionPayload

	if err := d.Decode(&body); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'provision' payload")
		return usecase.ProvisionDTO{}, err
	}

	vars := mux.Vars(r)
	itemID := vars["itemID"]

	id, err := strconv.Atoi(itemID)

	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update status' itemID")
		return usecase.ProvisionDTO{}, ErrParseItemID
	}

	dto := usecase.ProvisionDTO{
		Qty:             body.Qty,
		InventoryItemID: id,
	}

	return dto, nil
}
