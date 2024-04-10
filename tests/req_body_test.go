package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmviniciuss/forger/core"
	"github.com/bmviniciuss/forger/mux"
	"github.com/stretchr/testify/assert"
)

func Test_RequestBodyFunction(t *testing.T) {
	t.Run("should return interpolated value from request body", func(t *testing.T) {
		defs := []core.RouteDefinition{
			{
				Path:   "/{id}",
				Method: "POST",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: http.StatusTeapot,
					Body: `{
						"id": {{ requestBody "id" }}
					}`,
				},
			},
		}
		r := mux.NewStaticRouter(defs)
		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader([]byte(`{"id": 1}`)))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTeapot, rec.Code)
		assert.JSONEq(t, `{"id": 1}`, rec.Body.String())
	})

	t.Run("should return fallback if provided when key does not exists in body", func(t *testing.T) {
		defs := []core.RouteDefinition{
			{
				Path:   "/{id}",
				Method: "POST",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: http.StatusTeapot,
					Body: `{
						"id": "{{ requestBody "id" "fallback-id" }}"
					}`,
				},
			},
		}
		r := mux.NewStaticRouter(defs)
		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader([]byte(`{"id_name": 1}`)))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTeapot, rec.Code)
		assert.JSONEq(t, `{"id": "fallback-id"}`, rec.Body.String())
	})

	t.Run("should return empty value if fallback is not provided", func(t *testing.T) {
		defs := []core.RouteDefinition{
			{
				Path:   "/{id}",
				Method: "POST",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: http.StatusTeapot,
					Body: `{
						"id": "{{ requestBody "id" }}"
					}`,
				},
			},
		}
		r := mux.NewStaticRouter(defs)
		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader([]byte(`{"id_name": 1}`)))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTeapot, rec.Code)
		assert.JSONEq(t, `{"id": ""}`, rec.Body.String())
	})

	t.Run("should be able to use function as fallback", func(t *testing.T) {
		defs := []core.RouteDefinition{
			{
				Path:   "/{id}",
				Method: "POST",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: http.StatusTeapot,
					Body: `{
						"id": "{{ requestBody "id" (uuid "ulid") }}"
					}`,
				},
			},
		}
		r := mux.NewStaticRouter(defs)
		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader([]byte(`{"id_name": 1}`)))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTeapot, rec.Code)
		fmt.Println(string(rec.Body.String()))
		var data map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &data)
		assert.NotEmpty(t, data["id"])
	})

	t.Run("should return whole body", func(t *testing.T) {
		defs := []core.RouteDefinition{
			{
				Path:   "/{id}",
				Method: "POST",
				Response: core.RouteResponse{
					Type:       core.RESPONSE_TYPE_DYNAMIC,
					StatusCode: http.StatusTeapot,
					Body: `{
						"id": {{ requestBody }}
					}`,
				},
			},
		}
		r := mux.NewStaticRouter(defs)
		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader([]byte(`{"id": 1, "name": "Vinicius"}`)))
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTeapot, rec.Code)
		assert.JSONEq(t, `{"id": {"id": 1, "name": "Vinicius"}}`, rec.Body.String())
	})
}
