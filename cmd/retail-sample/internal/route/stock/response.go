package stock

import (
	"net/http"
	"time"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/machine/stock"
	persistence "github.com/anatollupacescu/retail-sample/cmd/retail-sample/internal/persistence/postgres"
)

type (
	entity struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Qty  int    `json:"qty"`
	}
	single struct {
		Data entity `json:"data"`
	}
	collection struct {
		Data []entity `json:"data"`
	}
)

func toCollectionResponse(entries []stock.Position) collection {
	var response collection
	response.Data = make([]entity, 0, len(entries))

	for i := range entries {
		entry := entries[i]

		response.Data = append(response.Data, entity{
			ID:   entry.ID,
			Name: entry.Name,
			Qty:  entry.Qty,
		})
	}

	return response
}

func toResponse(pos stock.Position) single {
	return single{
		Data: entity{
			ID:   pos.ID,
			Qty:  pos.Qty,
			Name: pos.Name,
		},
	}
}

type provisionLogEntity struct {
	Time time.Time `json:"time"`
	ID   int       `json:"id"`
	Qty  int       `json:"qty"`
}

func toProvisionLog(pl []persistence.ProvisionEntry) interface{} {
	response := struct {
		Data []provisionLogEntity `json:"data"`
	}{
		Data: make([]provisionLogEntity, 0),
	}

	for i := range pl {
		in := pl[i]

		response.Data = append(response.Data, provisionLogEntity{
			ID:  in.ID,
			Qty: in.Qty,
		})
	}

	return response
}

func httpServerError(w http.ResponseWriter) {
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
}
