package items

import (
	"api/app/mock"
  "api/app/models"
	"net/http"
	"net/http/httptest"
	"testing"
  "fmt"
  "strings"
  "encoding/json"

	"github.com/gin-gonic/gin"
)

func TestItem(t *testing.T) {
	router := gin.Default()
	Configure(router, nil)

	// Inject our mock into our handler.
	var is mock.ItemService
	Is = &is

	// Mock our User() call.
	is.ItemFn = func(id string) (*models.Item, error) {
		if id != "100" {
			t.Fatalf(fmt.Sprintf("unexpected id: %s", id))
		}
		return &models.Item{ID: "100", Name: "DaItam", Description: "Elnesto"}, nil
	}

	// Invoke the handler.
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/item/100", nil)
	router.ServeHTTP(w, r)

  result_item := models.Item{}
  result_bytes := w.Body.Bytes()
  err := json.Unmarshal(result_bytes, &result_item)
  if err != nil {
    t.Fatal(err.Error())
  }
  if result_item.Name != "DaItam" || result_item.ID != "100" {
    t.Log("Unexpected item: ")
    t.Fatal(result_item)
  }
	// Validate mock.
	if !is.ItemInvoked {
		t.Fatal("expected Item() to be invoked")
	}
}

func TestItems(t *testing.T) {
	router := gin.Default()
	Configure(router, nil)

	var is mock.ItemService
	Is = &is

	is.ItemsFn = func() ([]*models.Item, error) {
    var result []*models.Item
    var item models.Item
    item.ID = "100"
    item.Description = "Elnesto"
    item.Name = "DaItam"
    result = append(result, &item)
		return result, nil
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/item", nil)
	router.ServeHTTP(w, r)
  if w.Code != 200 {
    t.Log("Expected code 200, got: \n")
    t.Fatal(w.Code)
  }
	if !is.ItemsInvoked {
		t.Fatal("expected Item() to be invoked")
	}
}

func TestPostItem(t *testing.T) {
  router := gin.Default()
  Configure(router, nil)
  var is mock.ItemService
  Is = &is

  is.CreateItemFn = func(i *models.Item) error {
    if i.Description != "Elnesto" {
      t.Log("Unexpected item description: ")
      t.Fatal("Expected item description: Elnesto")
    }
    i.ID = "100"
    return nil
  }

  msg := `{"name":"DaItam", "description":"Elnesto"}`
	w := httptest.NewRecorder()
  reader := strings.NewReader(msg)
	r, _ := http.NewRequest("POST", "/item", reader)
	router.ServeHTTP(w, r)
  if !is.CreateItemInvoked {
    t.Fatal("expected CreateItem() to be invoked")
  }
  if w.Code != 201 {
    t.Log("Bad code. Expected 201, got: \n")
    t.Fatal(w.Code)
  }
}

func TestDeleteItem(t *testing.T) {
  router := gin.Default()
  Configure(router, nil)
  var is mock.ItemService
  Is = &is

  is.DeleteItemFn = func(id string) (*models.Item, error) {
    if id != "100" {
      t.Fatal("Expected id: 100")
    }
    return &models.Item{"100", "DaItam", "Elnesto"}, nil
  }
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/item/100", nil)
	router.ServeHTTP(w, r)
  t.Log(w.Code)
  t.Log(w.Body)
  if !is.DeleteItemInvoked {
    t.Fatal("expected DeleteItem() to be invoked")
  }
  if w.Code != 200 {
    t.Log("Bad code. Expected 200, got: \n")
    t.Fatal(w.Code)
  }
}
