package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
	"github.com/anatollupacescu/retail-sample/cmd/acceptance/helper"
	"github.com/bxcodec/faker/v3"
)

// group:inventory
func testCreateInventoryItem(t *runner.T) {
	t.Run("given name empty", func(*testing.T) { //error
		post := helper.Post("inventory")
		body := strings.NewReader(`{
			"name":""
		}`)
		r, err := post(body)
		t.Run("assert error", func(*testing.T) {
			if err != nil {
				t.Error("unexpected error", err)
			}
			if http.StatusBadRequest != r.StatusCode {
				t.Errorf("1: status code: want %d, got %d", http.StatusBadRequest, r.StatusCode)
			}
		})
	})
	t.Run("given name non unique", func(*testing.T) {
		post := helper.Post("inventory")
		body := strings.NewReader(`{
			"name":"test"
		}`)
		post(body)
		body = strings.NewReader(`{
			"name":"test"
		}`)
		r, err := post(body)
		t.Run("assert error", func(*testing.T) {
			if err != nil {
				t.Error("request failed")
			}
			if r.StatusCode != 400 {
				t.Errorf("2: status code: want %d, got %d", 400, r.StatusCode)
			}
			defer r.Body.Close()
			res, _ := ioutil.ReadAll(r.Body)
			if string(res) != "create item with name 'test': item type already present: bad request\n" {
				t.Error("expected error response, got", string(res))
			}
		})
	})
	t.Run("given item is saved", func(*testing.T) {
		post := helper.Post("inventory")
		word := faker.Word()
		body := strings.NewReader(fmt.Sprintf(`{
			"name": "%s"
		}`, word))
		r, err := post(body)
		t.Run("assert success", func(*testing.T) {
			if err != nil {
				t.Error("unexpected error", err)
			}
			if r.StatusCode != http.StatusCreated {
				t.Errorf("%s: status code: want %d, got %d", word, http.StatusCreated, r.StatusCode)
			}
		})
	})
}

// group:inventory
func testDisableItem(t *runner.T) {
	t.Run("given item is disabled", func(*testing.T) {
		post := helper.Patch("inventory", 1)
		body := strings.NewReader(`{
			"enabled": false
		}`)
		r, err := post(body)
		t.Run("assert success", func(*testing.T) {
			if err != nil {
				t.Error("unexpected error", err)
			}
			if r.StatusCode != http.StatusAccepted {
				t.Errorf("status code: want %d, got %d", http.StatusAccepted, r.StatusCode)
			}
		})
	})
}

// group:inventory
func testEnableItem(t *runner.T) {
	t.Run("given item is enabled", func(*testing.T) {
		post := helper.Patch("inventory", 1)
		body := strings.NewReader(`{
			"enabled": true
		}`)
		r, err := post(body)
		t.Run("assert success", func(*testing.T) {
			if err != nil {
				t.Error("unexpected error", err)
			}
			if r.StatusCode != http.StatusAccepted {
				t.Errorf("status code: want %d, got %d", http.StatusAccepted, r.StatusCode)
			}
		})
	})
}
