package stock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/inventory"
	"github.com/anatollupacescu/retail-sample/domain/retail-sample/stock"
)

type (
	webApp struct {
		logger  middleware.Logger
		wrapper wrapper
	}
	entity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Qty  int    `json:"qty"`
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

func (a *webApp) getAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	entries, err := a.wrapper.currentStock()

	if err != nil {
		http.Error(w, internalServerErrorMsg(), http.StatusBadRequest)
		return
	}

	var response collectionResponse
	response.Data = make([]entity, 0, len(entries))

	for i := range entries {
		entry := entries[i]

		response.Data = append(response.Data, entity{
			ID:   entry.ID,
			Name: entry.Name,
			Qty:  entry.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "stock.getAll")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

type getFunc func() (*http.Response, error)

func GetAll(get getFunc) (d []stock.Position, err error) {
	response, err := get()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	defer func() {
		_ = response.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var pos collectionResponse

	if err = json.Unmarshal(respBody, &pos); err != nil {
		return nil, err
	}

	for i := range pos.Data {
		positionData := pos.Data[i]

		d = append(d, stock.Position{
			ID:   positionData.ID,
			Name: positionData.Name,
			Qty:  positionData.Qty,
		})
	}

	return d, nil
}

func (a *webApp) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	itemID, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "stock.get")
		http.Error(w, "could not parse id", http.StatusBadRequest)

		return
	}

	pos, err := a.wrapper.quantity(itemID)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = singleResponse{
		Data: entity{
			ID:   pos.ID,
			Qty:  pos.Qty,
			Name: pos.Name,
		},
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "stock.get")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

func Get(get getFunc) (pos stock.Position, err error) {
	response, err := get()
	if err != nil {
		return pos, err
	}

	if response.StatusCode != http.StatusOK {
		return pos, errors.New("unexpected status code")
	}

	defer func() {
		_ = response.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return pos, err
	}

	var e singleResponse
	if err = json.Unmarshal(respBody, &e); err != nil {
		return pos, err
	}

	return stock.Position{
		ID:   e.Data.ID,
		Name: e.Data.Name,
		Qty:  e.Data.Qty,
	}, nil
}

type provisionPayload struct {
	Qty int `json:"qty"`
}

type provisionResponseData struct {
	ID  int `json:"id"`
	Qty int `json:"qty"`
}

type provisionResponse struct {
	Data provisionResponseData `json:"data"`
}

func (a *webApp) provision(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	rid := vars["itemID"]

	itemID, err := strconv.Atoi(rid)

	if err != nil {
		a.logger.Log("action", "parse id", "error", err, "method", "stock.provision")
		http.Error(w, "could not parse id", http.StatusBadRequest)

		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var body provisionPayload

	if err = d.Decode(&body); err != nil {
		a.logger.Log("action", "decode request", "error", err, "method", "stock.provision")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	newQty, err := a.wrapper.provision(itemID, body.Qty)

	switch err {
	case nil:
		break
	case inventory.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
		return
	}

	var response = provisionResponse{
		Data: provisionResponseData{
			ID:  itemID,
			Qty: newQty,
		},
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "stock.provision")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}

type postFunc func(io.Reader) (*http.Response, error)

func Provision(qty int, post postFunc) (newQty int, err error) {
	payload := provisionPayload{
		Qty: qty,
	}

	var data []byte

	if data, err = json.Marshal(payload); err != nil {
		return 0, err
	}

	body := bytes.NewReader(data)

	response, err := post(body)
	if err != nil {
		return 0, err
	}

	switch response.StatusCode {
	case http.StatusCreated, http.StatusAccepted:
		break
	default:
		return 0, fmt.Errorf("unexpected status code: %v", response.Status)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var responseData provisionResponse

	if err = json.Unmarshal(respBody, &responseData); err != nil {
		return 0, err
	}

	return responseData.Data.Qty, nil
}

func (a *webApp) getProvisionLog(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	provisionLog, err := a.wrapper.getProvisionLog()

	if err != nil {
		a.logger.Log("action", "call application", "error", err, "method", "stock.provisionlog")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)

		return
	}

	type provisionLogEntity struct {
		Time time.Time `json:"time"`
		ID   int       `json:"id"`
		Qty  int       `json:"qty"`
	}

	var response struct {
		Data []provisionLogEntity `json:"data"`
	}

	response.Data = make([]provisionLogEntity, 0)

	for _, in := range provisionLog {
		response.Data = append(response.Data, provisionLogEntity{
			ID:  in.ID,
			Qty: in.Qty,
		})
	}

	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		a.logger.Log("action", "encode response", "error", err, "method", "stock.provisionlog")
		http.Error(w, internalServerErrorMsg(), http.StatusInternalServerError)
	}
}
