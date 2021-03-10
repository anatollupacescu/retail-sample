package tests

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:inventory
func testCreateInventoryItem(t *runner.T) {
	t.Run("given name empty", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given name non unique", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given item is saved", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}

// group:inventory
func testDisableItem(t *runner.T) {
	t.Run("given item is disabled", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}

// group:inventory
func testEnableItem(t *runner.T) {
	t.Run("given item is enabled", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}
