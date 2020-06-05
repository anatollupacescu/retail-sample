package inventory

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gojektech/heimdall/httpclient"
	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type (
	webApp struct {
		logger  middleware.Logger
		wrapper wrapper
	}
	entity struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	responseData struct {
		Data entity `json:"data"`
	}
)

var internalServerError = "internal server error"

func (a *webApp) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	id, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "inventory.update")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	type payload struct {
		Enabled bool `json:"enabled"`
	}

	var requestPayload payload

	if err := d.Decode(&requestPayload); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "inventory.update")
		http.Error(w, "parse body", http.StatusBadRequest)
		return
	}

	item, err := a.wrapper.setStatus(id, requestPayload.Enabled)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
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
			Name:    item.Name,
			Enabled: item.Enabled,
		},
	}

	w.WriteHeader(http.StatusAccepted)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "inventory.update")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

func Update() (item inventory.Item, err error) {
	return
}

type createPayload struct {
	Name string `json:"name"`
}

func (a *webApp) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body createPayload

	if err := d.Decode(&body); err != nil {
		a.logger.Log("action", "decode request payload", "error", err, "method", "inventory.create")
		http.Error(w, "could not parse body", http.StatusBadRequest)
		return
	}

	itemName := body.Name
	newItem, err := a.wrapper.create(itemName)

	switch err {
	case nil:
		break
	case inventory.ErrEmptyName:
		fallthrough
	case inventory.ErrDuplicateName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = responseData{
		Data: entity{
			ID:      int(newItem.ID),
			Name:    newItem.Name,
			Enabled: newItem.Enabled,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "inventory.create")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}

// client
func Create(url, name string) (item inventory.Item, err error) {
	timeout := 100 * time.Millisecond

	httpClient := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
	)
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	payload := createPayload{
		Name: name,
	}

	var data []byte

	data, err = json.Marshal(payload)
	if err != nil {
		return item, err
	}

	body := bytes.NewReader(data)

	response, err := httpClient.Post(url, body, headers)
	if err != nil {
		return item, err
	}

	if response.StatusCode != http.StatusCreated {
		return item, errors.New("unexpected status code")
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return item, err
	}

	var responseData responseData

	if err = json.Unmarshal(respBody, &responseData); err != nil {
		return item, err
	}

	return inventory.Item{
		ID:      responseData.Data.ID,
		Name:    responseData.Data.Name,
		Enabled: responseData.Data.Enabled,
	}, err
}

func (a *webApp) getAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	list, err := a.wrapper.getAll()

	if err != nil {
		http.Error(w, internalServerError, http.StatusBadRequest)
	}

	var response struct {
		Data []entity `json:"data"`
	}

	response.Data = make([]entity, 0)

	for _, tp := range list {
		response.Data = append(response.Data, entity{
			ID:      int(tp.ID),
			Name:    string(tp.Name),
			Enabled: tp.Enabled,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "inventory.getAll")
		http.Error(w, internalServerError, http.StatusBadRequest)
	}
}

func (a *webApp) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	reqID := vars["itemID"]

	id, err := strconv.Atoi(reqID)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "inventory.get")
		http.Error(w, "could not parse id", http.StatusBadRequest)
		return
	}

	inventoryItem, err := a.wrapper.getOne(id)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	var response = struct {
		Data entity `json:"data"`
	}{
		Data: entity{
			ID:      int(inventoryItem.ID),
			Name:    string(inventoryItem.Name),
			Enabled: inventoryItem.Enabled,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "inventory.get")
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}
}
