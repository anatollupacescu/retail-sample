package order

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	usecase "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/usecase/order"

	"github.com/anatollupacescu/retail-sample/domain/retail/inventory"
	domain "github.com/anatollupacescu/retail-sample/domain/retail/order"
	"github.com/anatollupacescu/retail-sample/domain/retail/stock"
)

type (
	entity struct {
		ID       int `json:"id"`
		RecipeID int `json:"recipeID"`
		Qty      int `json:"qty"`
	}

	singleResponse struct {
		Data entity `json:"data"`
	}
)

func internalServerErrorMsg() string {
	return http.StatusText(http.StatusInternalServerError)
}

type createPayload struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var payload createPayload

	if err := d.Decode(&payload); err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	recipeID := payload.ID
	orderQty := payload.Qty

	entryID, err := uc.Create(recipeID, orderQty)

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
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
		Data: entity{
			ID:       int(entryID),
			RecipeID: recipeID,
			Qty:      orderQty,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	orderID := vars["orderID"]

	ordr, err := uc.GetByID(orderID)

	switch err {
	case nil:
		break
	case domain.ErrOrderNotFound, usecase.ErrBadItemID:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:       int(ordr.ID),
			RecipeID: ordr.RecipeID,
			Qty:      ordr.Qty,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	all, err := uc.GetAll()

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	var response struct {
		Data []entity `json:"data"`
	}

	response.Data = make([]entity, 0)

	for i := range all {
		currentOrder := all[i]

		e := entity{
			ID:       int(currentOrder.ID),
			RecipeID: currentOrder.RecipeID,
			Qty:      currentOrder.Qty,
		}

		response.Data = append(response.Data, e)
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}
