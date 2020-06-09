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

	provision := arbor.New("provision stock", testProvision, createOk)
	getOneSP := arbor.New("get single stock position", testGetStockPos, provision)
	getAllSP := arbor.New("get all stock positions", testGetAllStockPos, provision)

	all := arbor.Suite("all", createEmpty, getOne, getAll, noDuplicate, disable, getOneSP, getAllSP)

	all.Run()

	t.Run("succeeds", func(t *testing.T) {
		assert.Equal(t, true, all.Success)
	})

	t.Logf("%s\n", all)

	report := arbor.Marshal(createEmpty, getOne, getAll, noDuplicate, disable, getOneSP, getAllSP)

	arbor.Upload(*arborURL, report)
}
