package order

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/hlog"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase"
	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dto, err := newCreateDTO(r)
	if err != nil {
		httpServerError(w)
		return
	}

	uc, err := newUseCase(r)
	if err != nil {
		httpServerError(w)
		return
	}

	newOrder, err := uc.PlaceOrder(dto)

	switch err {
	case nil:
		break
	case domain.ErrInvalidRecipe,
		domain.ErrInvalidQuantity,
		stock.ErrNotEnoughStock,
		inventory.ErrDuplicateName,
		inventory.ErrEmptyName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)

	var response = toSingleResponse(newOrder)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

type createPayload struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

func newCreateDTO(r *http.Request) (usecase.PlaceOrderDTO, error) {
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

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	order, err := getByID(r)

	switch err {
	case nil:
		break
	case domain.ErrOrderNotFound, ErrBadItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	response := toSingleResponse(order)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orders, err := getAll(r)

	if err != nil {
		httpServerError(w)
		return
	}

	response := toCollectionResponse(orders)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		httpServerError(w)
	}
}
