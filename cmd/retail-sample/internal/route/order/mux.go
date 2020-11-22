package order

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

func errorMsg() string {
	return http.StatusText(http.StatusInternalServerError)
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, errorMsg(), http.StatusInternalServerError)
		return
	}

	dto, err := toCreateDTO(r)
	if err != nil {
		http.Error(w, errorMsg(), http.StatusInternalServerError)
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
		http.Error(w, errorMsg(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	var response = toSingleResponse(newOrder)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, errorMsg(), http.StatusBadRequest)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	order, err := GetByID(r)

	switch err {
	case nil:
		break
	case domain.ErrOrderNotFound, ErrBadItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, errorMsg(), http.StatusInternalServerError)
		return
	}

	response := toSingleResponse(order)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, errorMsg(), http.StatusInternalServerError)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orders, err := List(r)

	if err != nil {
		http.Error(w, errorMsg(), http.StatusInternalServerError)
		return
	}

	response := toCollectionResponse(orders)
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, errorMsg(), http.StatusInternalServerError)
	}
}
