package recipe

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type (
	item struct {
		ID  int `json:"id"`
		Qty int `json:"qty"`
	}

	entity struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Items   []item `json:"items"`
		Enabled bool   `json:"enabled"`
	}

	singleResponse struct {
		Data entity `json:"data"`
	}

	collectionResponse struct {
		Data []entity `json:"data"`
	}
)

func internalServerErrorMsg() string {
	return http.StatusText(http.StatusInternalServerError)
}

type createPayload struct {
	Name  string      `json:"name"` // pointer so we can test for field absence
	Items map[int]int `json:"items"`
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

	var requestBody createPayload

	if err := d.Decode(&requestBody); err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var ingredients = make([]recipe.Ingredient, 0, len(requestBody.Items))

	for id, qty := range requestBody.Items {
		ingredients = append(ingredients, recipe.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	var recipeName = recipe.Name(requestBody.Name)

	re, err := uc.Create(recipeName, ingredients)

	switch err {
	case nil:
		break
	case
		recipe.ErrEmptyName,
		recipe.ErrIgredientNotFound,
		recipe.ErrQuantityNotProvided,
		recipe.ErrNoIngredients:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
		Data: entity{
			ID:      int(re.ID),
			Name:    string(re.Name),
			Items:   toItems(re.Ingredients),
			Enabled: re.Enabled,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func toIngredientCollection(items []item) []recipe.Ingredient {
	ingredients := make([]recipe.Ingredient, 0, len(items))

	for i := range items {
		currentItem := items[i]

		ingredients = append(ingredients, recipe.Ingredient{
			ID:  currentItem.ID,
			Qty: currentItem.Qty,
		})
	}

	return ingredients
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	var response struct {
		Data []entity `json:"data"`
	}

	all, err := uc.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.Data = make([]entity, 0, len(all))

	for i := range all {
		r := all[i]

		response.Data = append(response.Data, entity{
			ID:      int(r.ID),
			Name:    string(r.Name),
			Items:   toItems(r.Ingredients),
			Enabled: r.Enabled,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func toItems(i []recipe.Ingredient) (items []item) {
	for _, ri := range i {
		items = append(items, item{
			ID:  ri.ID,
			Qty: ri.Qty,
		})
	}

	return
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	recipeID := vars["recipeID"]

	rcp, err := uc.GetByID(recipeID)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
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
			ID:      int(rcp.ID),
			Name:    string(rcp.Name),
			Items:   toItems(rcp.Ingredients),
			Enabled: rcp.Enabled,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uc, err := useCase(r)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody updatePayload

	if err = d.Decode(&requestBody); err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	vars := mux.Vars(r)
	recipeID := vars["recipeID"]

	re, err := uc.Update(recipeID, requestBody.Enabled)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
		Data: entity{
			ID:      int(re.ID),
			Name:    string(re.Name),
			Items:   toItems(re.Ingredients),
			Enabled: re.Enabled,
		},
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}
