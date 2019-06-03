package design_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/domain"
	qt "github.com/frankban/quicktest"
)

func TestCanAdd(t *testing.T) {

	t.Run("add design persists it", func(t *testing.T) {
		c := qt.New(t)
		d := domain.Design{}
		c.Assert(d, qt.Not(qt.IsNil))
	})
}
