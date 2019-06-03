package domain_test

import (
	"testing"

	"github.com/anatollupacescu/retail-sample/pkg/retail-sample/domain"
	qt "github.com/frankban/quicktest"
)

func TestDesignValidation(t *testing.T) {
	t.Run("design with empty values returns error", func(t *testing.T) {
		c := qt.New(t)
		d := domain.Design{}
		c.Assert(d.Validate(), qt.Equals, false)
	})
	t.Run("design with invalid name returns error", func(t *testing.T) {
		c := qt.New(t)
		d := domain.Design{
			Name:  "lol",
			Price: 0.1,
		}
		c.Assert(d.Validate(), qt.Equals, false)
	})
	t.Run("valid design succeeds", func(t *testing.T) {
		c := qt.New(t)
		d := domain.Design{
			Name:  "T-Shirt",
			Price: 12.5,
		}
		c.Assert(d.Validate(), qt.Equals, true)
	})
}
