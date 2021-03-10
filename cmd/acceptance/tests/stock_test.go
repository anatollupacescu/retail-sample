package tests

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:stock after:inventory
func testProvision(t *runner.T) {
	t.Run("given quantity is negative", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given position updated", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}

// group:stock
func testExtract(t *runner.T) {
	t.Run("given quantity is negative", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given recipe not found", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given recipe is invalid", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given item not present in stock", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given not enough stock for item", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given position updated", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}
