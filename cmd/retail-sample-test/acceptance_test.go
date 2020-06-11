// +build acceptance

package acceptance_test

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anatollupacescu/retail-sample/internal/arbor"
)

var arborURL = flag.String("arborURL", "", "graph server URL")

func TestAcceptance(t *testing.T) {
	createOk := arbor.New("can create", testCreate)
	createEmpty := arbor.New("rejects empty name", testCreateWithEmptyName, createOk)

	getOne := arbor.New("get one", testGetOne, createOk)
	getAll := arbor.New("get all", testGetAll, createOk)
	noDuplicate := arbor.New("no duplicate", testDuplicate, createOk)
	disable := arbor.New("disable", testDisable, createOk)

	getOneSP := arbor.New("get single stock position", testGetStockPos, createOk)
	provision := arbor.New("provision stock", testProvision, getOneSP)
	getAllSP := arbor.New("get all stock positions", testGetAllStockPos, createOk)

	all, success := arbor.Run(createEmpty, getOne, getAll, noDuplicate, disable, provision, getAllSP)

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, true, success)
	})

	t.Logf("\n%s\n", all)

	report := arbor.Marshal(all...)

	arbor.Upload(*arborURL, report)
}
