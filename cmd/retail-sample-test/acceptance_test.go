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
	createEmpty := arbor.New("rejects empty name", testCreateWithEmptyName)
	createOk := arbor.New("can create", testCreate)

	create := arbor.Suite("create", createEmpty, createOk)

	getOne := arbor.New("get one", testGetOne, create)
	getAll := arbor.New("get all", testGetAll, create)
	noDuplicate := arbor.New("no duplicate", testDuplicate, create)
	disable := arbor.New("disable", testDisable, create)

	inv := arbor.Suite("inv", getAll, getOne, noDuplicate, disable)

	provision := arbor.New("provision stock", testProvision, create)
	getOneSP := arbor.New("get single stock position", testGetStockPos, provision)
	getAllSP := arbor.New("get all stock positions", testGetAllStockPos, provision)

	stock := arbor.Suite("stock", getOneSP, getAllSP)

	all := arbor.Suite("all", inv, stock)

	all.Run()

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, true, all.Success)
	})

	t.Logf("%s\n", all)

	report := arbor.Marshal(create, createEmpty, createOk, getOne, getAll, noDuplicate, disable, provision, getOneSP, getAllSP)

	arbor.Upload(*arborURL, report)
}
