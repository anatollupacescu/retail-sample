package order

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type createPayload struct {
	RecipeID int `json:"id"`
	Count    int `json:"qty"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var payload createPayload

	if err := d.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uc, err := newUseCase(r)
	if err != nil {
		httpServerError(w)
		return
	}

	newOrder, err := uc.Create(payload.RecipeID, payload.Count)

	switch err {
	case nil:
		break
	case recipe.ErrDisabled,
		order.ErrInvalidQuantity,
		stock.ErrNotEnoughStock:
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

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ordr, err := getByID(r)

	switch err {
	case nil:
		break
	case order.ErrOrderNotFound, ErrBadItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		httpServerError(w)
		return
	}

	response := toSingleResponse(ordr)
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
