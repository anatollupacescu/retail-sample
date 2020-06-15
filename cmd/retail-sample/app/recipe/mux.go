package recipe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type (
	webApp struct {
		wrapper wrapper
		logger  middleware.Logger
	}

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

var internalServerError = "internal server error"

type createPayload struct {
	Name  string      `json:"name"` // pointer so we can test for field absence
	Items map[int]int `json:"items"`
}

func (a *webApp) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody createPayload

	if err := d.Decode(&requestBody); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "recipe.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var ingredients []recipe.Ingredient

	for id, qty := range requestBody.Items {
		ingredients = append(ingredients, recipe.Ingredient{
			ID:  id,
			Qty: qty,
		})
	}

	var recipeName = recipe.Name(requestBody.Name)

	re, err := a.wrapper.create(recipeName, ingredients)

	switch err {
	case nil:
		break
	case recipe.ErrEmptyName, recipe.ErrIgredientNotFound, recipe.ErrNoIngredients:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
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
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func Create(name string, items map[int]int, post func(io.Reader) (*http.Response, error)) (r recipe.Recipe, err error) {
	payload := createPayload{
		Name:  name,
		Items: items,
	}

	var data []byte

	data, err = json.Marshal(payload)
	if err != nil {
		return r, err
	}

	body := bytes.NewReader(data)

	response, err := post(body)
	if err != nil {
		return r, err
	}

	if response.StatusCode > 299 {
		return r, fmt.Errorf("unexpected status code: %v", response.StatusCode)
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return r, err
	}

	var responseData singleResponse

	if err = json.Unmarshal(respBody, &responseData); err != nil {
		return r, err
	}

	r = recipe.Recipe{
		ID:          recipe.ID(responseData.Data.ID),
		Name:        recipe.Name(responseData.Data.Name),
		Enabled:     responseData.Data.Enabled,
		Ingredients: toIngredientCollection(responseData.Data.Items),
	}
	return r, nil
}

func toIngredientCollection(ii []item) (r []recipe.Ingredient) {
	for _, i := range ii {
		r = append(r, recipe.Ingredient{
			ID:  i.ID,
			Qty: i.Qty,
		})
	}

	return
}

func (a *webApp) getAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	all, err := a.wrapper.getAll()

	if err != nil {
		http.Error(w, internalServerError, http.StatusBadRequest)
		return
	}

	var response struct {
		Data []entity `json:"data"`
	}

	response.Data = make([]entity, 0)

	for _, r := range all {
		response.Data = append(response.Data, entity{
			ID:      int(r.ID),
			Name:    string(r.Name),
			Items:   toItems(r.Ingredients),
			Enabled: r.Enabled,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.getAll")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

// client
func GetAll(get func() (*http.Response, error)) (rs []recipe.Recipe, err error) {
	response, err := get()
	if err != nil {
		return rs, err
	}

	if response.StatusCode != http.StatusOK {
		return rs, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return rs, err
	}

	var rcps collectionResponse

	if err = json.Unmarshal(respBody, &rcps); err != nil {
		return rs, err
	}

	for _, r := range rcps.Data {
		rs = append(rs, recipe.Recipe{
			Enabled:     r.Enabled,
			ID:          recipe.ID(r.ID),
			Ingredients: toIngredientCollection(r.Items),
			Name:        recipe.Name(r.Name),
		})
	}

	return rs, nil
}

func toItems(i []recipe.Ingredient) (items []item) {
	for _, ri := range i {
		items = append(items, item{
			ID:  int(ri.ID),
			Qty: int(ri.Qty),
		})
	}

	return
}

func (a *webApp) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["recipeID"]

	id, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "recipe.get")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	recipeID := recipe.ID(id)

	rcp, err := a.wrapper.get(recipeID)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:      id,
			Name:    string(rcp.Name),
			Items:   toItems(rcp.Ingredients),
			Enabled: rcp.Enabled,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.get")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

// client
func Get(get func() (*http.Response, error)) (r recipe.Recipe, err error) {
	response, err := get()
	if err != nil {
		return r, err
	}

	if response.StatusCode != http.StatusOK {
		return r, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return r, err
	}

	var one singleResponse

	if err = json.Unmarshal(respBody, &one); err != nil {
		return r, err
	}

	return recipe.Recipe{
		Enabled:     one.Data.Enabled,
		ID:          recipe.ID(one.Data.ID),
		Ingredients: toIngredientCollection(one.Data.Items),
		Name:        recipe.Name(one.Data.Name),
	}, nil
}

type updatePayload struct {
	Enabled bool `json:"enabled"`
}

func (a *webApp) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["recipeID"]

	id, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "recipe.update")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var requestBody updatePayload

	if err := d.Decode(&requestBody); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "recipe.update")
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	re, err := a.wrapper.setStatus(id, requestBody.Enabled)

	switch err {
	case nil:
		break
	case recipe.ErrRecipeNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
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
		a.logger.Log("action", "encode response", "error", err, "method", "recipe.update")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func Update(enabled bool, patch func(io.Reader) (*http.Response, error)) (item recipe.Recipe, err error) {
	payload := updatePayload{
		Enabled: enabled,
	}

	var data []byte

	if data, err = json.Marshal(payload); err != nil {
		return item, err
	}

	body := bytes.NewReader(data)

	response, err := patch(body)
	if err != nil {
		return item, err
	}

	if response.StatusCode > 299 {
		return item, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return item, err
	}

	var responseData singleResponse

	if err = json.Unmarshal(respBody, &responseData); err != nil {
		return item, err
	}

	var ingredients []recipe.Ingredient

	for _, r := range responseData.Data.Items {
		ingredients = append(ingredients, recipe.Ingredient{
			ID:  r.ID,
			Qty: r.Qty,
		})
	}

	return recipe.Recipe{
		ID:          recipe.ID(responseData.Data.ID),
		Name:        recipe.Name(responseData.Data.Name),
		Enabled:     responseData.Data.Enabled,
		Ingredients: ingredients,
	}, nil
}
