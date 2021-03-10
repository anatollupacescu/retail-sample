package tests

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:order after:recipe
func testCreateOrder(t *runner.T) {
	t.Run("given quantity is invalid", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given recipe is not valid", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given stock is updated and order created", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}
