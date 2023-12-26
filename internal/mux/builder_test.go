package mux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmviniciuss/forger-golang/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func StringToMap(str string) map[string]interface{} {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &m)
	if err != nil {
		panic(err)
	}
	return m
}

func TestMux_UUID_Generator(t *testing.T) {
	defs := []core.RouteDefinition{
		{
			Path:   "/item/{id}",
			Method: "GET",
			Response: core.RouteResponse{
				Type:       core.RESPONSE_TYPE_DYNAMIC,
				StatusCode: 200,
				Body: `{
					"id": "{{ uuid "ulid" }}"
				}`,
			},
		},
	}

	t.Run("[uuid] should generate uuid in request", func(t *testing.T) {
		mux := New(defs)
		itemID := uuid.NewString()
		req, err := http.NewRequest("GET", fmt.Sprintf("/item/%s", itemID), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		m := StringToMap(rr.Body.String())
		_, ok := m["id"]
		assert.True(t, ok)
		_, ok = m["id"].(string)
		assert.True(t, ok)
		_, err = uuid.Parse(m["id"].(string))
		assert.Nil(t, err)
	})

	t.Run("[uuid] should generate uuid of type ulid in request", func(t *testing.T) {
		mux := New([]core.RouteDefinition{
			{
				Path:   "/item/{id}",
				Method: "GET",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: 200,
					Body: `{
						"id": "{{ uuid "ulid" }}"
					}`,
				},
			},
		})
		itemID := uuid.NewString()
		req, err := http.NewRequest("GET", fmt.Sprintf("/item/%s", itemID), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		m := StringToMap(rr.Body.String())
		_, ok := m["id"]
		assert.True(t, ok)
		_, ok = m["id"].(string)
		assert.True(t, ok)
		_, err = uuid.Parse(m["id"].(string))
		assert.Nil(t, err)
	})

}

func TestMux_RequestVar(t *testing.T) {
	t.Run("[requestVar] should access request var", func(t *testing.T) {
		mux := New([]core.RouteDefinition{
			{
				Path:   "/item/{id}",
				Method: "GET",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: 200,
					Body: `{
						"id": "{{ requestVar "id" }}"
					}`,
				},
			},
		})
		itemID := uuid.NewString()
		req, err := http.NewRequest("GET", fmt.Sprintf("/item/%s", itemID), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		m := StringToMap(rr.Body.String())
		_, ok := m["id"]
		assert.True(t, ok)
		_, ok = m["id"].(string)
		assert.True(t, ok)
		assert.Equal(t, itemID, m["id"])
	})

	t.Run("[requestVar] should return empty value if parameter does not exists", func(t *testing.T) {
		mux := New([]core.RouteDefinition{
			{
				Path:   "/item/{id}",
				Method: "GET",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: 200,
					Body: `{
						"id": "{{ requestVar "nid" }}"
					}`,
				},
			},
		})
		itemID := uuid.NewString()
		req, err := http.NewRequest("GET", fmt.Sprintf("/item/%s", itemID), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		m := StringToMap(rr.Body.String())
		_, ok := m["id"]
		assert.True(t, ok)
		assert.Equal(t, "", m["id"])
	})
}
