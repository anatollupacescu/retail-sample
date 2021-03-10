package tests

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

// group:recipe after:inventory
func testCreateRecipe(t *runner.T) {
	t.Run("given empty name", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given no ingredients", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given invalid ingredient quantity", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given duplicate ingredient", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given ingredient not found", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given ingredient disabled", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given recipe name not unique", func(*testing.T) {
		t.Run("assert error", func(*testing.T) {
		})
	})
	t.Run("given recipe created", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}

// group:recipe
func testDisableRecipe(t *runner.T) {
	t.Run("given recipe is disabled", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}

// group:recipe
func testEnableRecipe(t *runner.T) {
	t.Run("given recipe is enabled", func(*testing.T) {
		t.Run("assert success", func(*testing.T) {
		})
	})
}
