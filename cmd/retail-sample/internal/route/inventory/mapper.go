package inventory

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/inventory"
)

var (
	ErrParseBody   = errors.New("could not parse body")
	ErrParseItemID = errors.New("could not parse item ID")
)

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

func toUpdateStatusDTO(r *http.Request) (usecase.UpdateStatusDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestPayload updatePayload
	if err := d.Decode(&requestPayload); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update status' payload")
		return usecase.UpdateStatusDTO{}, ErrParseBody
	}

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update status' itemID")
		return usecase.UpdateStatusDTO{}, ErrParseItemID
	}

	dto := usecase.UpdateStatusDTO{
		ID:      id,
		Enabled: requestPayload.Enabled,
	}

	return dto, nil
}

type createPayload struct {
	Name string `json:"name"`
}

func toCreateDTO(r *http.Request) (usecase.CreateDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body createPayload

	if err := d.Decode(&body); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'create' payload")
		return usecase.CreateDTO{}, ErrParseBody
	}

	dto := usecase.CreateDTO{
		Name: body.Name,
	}

	return dto, nil
}
