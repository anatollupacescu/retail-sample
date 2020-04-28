package web

import (
	"encoding/json"
	"net/http"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

func (a *WebApp) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields

	var requestBody struct {
		Name  *string     `json:"name"` // pointer so we can test for field absence
		Items map[int]int `json:"items"`
	}

	if err := d.Decode(&requestBody); err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	if requestBody.Name == nil {
		http.Error(w, "name is mandatory", http.StatusBadRequest)
		return
	}

	var ingredients []recipe.Ingredient

	for id, qty := range requestBody.Items {
		ingredients = append(ingredients, recipe.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	var recipeName = recipe.Name(*requestBody.Name)

	recipeID, err := a.RecipeBook.Add(recipeName, ingredients)

	switch err {
	case nil:
		break
	case recipe.ErrEmptyName:
		fallthrough
	case recipe.ErrIgredientNotFound:
		fallthrough
	case recipe.ErrNoIngredients:
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	type DescriptorEntity map[recipe.Name]recipe.ID

	var response = struct {
		Data DescriptorEntity `json:"data"`
	}{
		Data: DescriptorEntity{
			recipeName: recipeID,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func (a *WebApp) ListRecipes(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type recipe struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Items []item `json:"items"`
	}

	var response struct {
		Data []recipe `json:"data"`
	}

	response.Data = make([]recipe, 0) //to have  '[]' instead of null

	list, err := a.RecipeBook.List()

	switch err {
	//TODO
	}

	for _, r := range list {
		response.Data = append(response.Data, recipe{
			ID:    int(r.ID),
			Name:  string(r.Name),
			Items: toItems(r.Ingredients),
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

type item struct {
	Id  int `json:"id"`
	Qty int `json:"qty"`
}

func toItems(i []recipe.Ingredient) (items []item) {
	for _, ri := range i {
		items = append(items, item{
			Id:  int(ri.ID),
			Qty: int(ri.Qty),
		})
	}

	return
}

func (a *WebApp) GetRecipe(w http.ResponseWriter, _ *http.Request) {
	panic("not implemented")
}
