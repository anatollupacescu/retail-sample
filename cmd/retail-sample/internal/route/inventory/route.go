package inventory

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	inventory2 "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
)

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dto, err := newUpdateDTO(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := newUseCase(r)

	if err != nil {
		httpServerError(w)
		return
	}

	item, err := uc.UpdateStatus(dto)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		httpServerError(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	response := toSingleResponse(item)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

var ErrParseItemID = errors.New("could not parse item ID")

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

func newUpdateDTO(r *http.Request) (inventory2.UpdateInventoryItemStatusDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestPayload updatePayload
	if err := d.Decode(&requestPayload); err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update status' payload")
		return inventory2.UpdateInventoryItemStatusDTO{}, ErrParseBody
	}

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["itemID"])
	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("parse 'update status' itemID")
		return inventory2.UpdateInventoryItemStatusDTO{}, ErrParseItemID
	}

	dto := inventory2.UpdateInventoryItemStatusDTO{
		ID:      id,
		Enabled: requestPayload.Enabled,
	}

	return dto, nil
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dto, err := newCreateDTO(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := newUseCase(r)

	if err != nil {
		httpServerError(w)
		return
	}

	newItem, err := uc.Create(dto)

	switch err {
	case nil:
	case inventory.ErrEmptyName, inventory.ErrDuplicateName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)

	response := toSingleResponse(newItem)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

type createPayload struct {
	Name string `json:"name"`
}

var ErrParseBody = errors.New("could not parse body")

func newCreateDTO(r *http.Request) (inventory2.CreateInventoryItemDTO, error) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body createPayload

	if err := d.Decode(&body); err != nil {
		return inventory2.CreateInventoryItemDTO{}, ErrParseBody
	}

	dto := inventory2.CreateInventoryItemDTO{
		Name: body.Name,
	}

	return dto, nil
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	item, err := inventoryItemByID(r)

	switch err {
	case nil:
	case ErrParseItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		httpServerError(w)
		return
	}

	response := toSingleResponse(item)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	all, err := allInventoryItems(r)

	if err != nil {
		httpServerError(w)
		return
	}

	response := toCollectionResponse(all)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}
