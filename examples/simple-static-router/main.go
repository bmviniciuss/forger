package main

import (
	"fmt"
	"net/http"

	"github.com/bmviniciuss/forger/internal/core"
	"github.com/bmviniciuss/forger/mux"
)

func main() {
	defs := []core.RouteDefinition{
		{
			Path:   "/items",
			Method: "GET",
			Response: core.RouteResponse{
				Type:       core.RESPONSE_TYPE_STATIC,
				StatusCode: 200,
				Body:       `{"id": 1, "name": "Item 1"}`,
				Headers:    map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			Path:   "/items/{id}",
			Method: "POST",
			Response: core.RouteResponse{
				Type:       core.RESPONSE_TYPE_DYNAMIC,
				StatusCode: http.StatusTeapot,
				Body: `{
					"id": "{{ requestVar "id"}}",
					"item": {{ requestBody "item.prices.0.value" }},
					"page": "{{ requestQuery "page" }}",
					"client_id": "{{ requestHeader "client-id" }}",
					"random_uuid": "{{ uuid "ulid" }}",
					"time": "{{ time "iso8601" }}"
				}`,
				Headers: map[string]string{
					"Content-Type":  "application/json",
					"Item-ID":       `{{ requestVar "id"}}`,
					"page":          `{{ requestQuery "page" }}`,
					"res-client-id": `{{ requestHeader "client-id" }}`,
				},
			},
		},
	}

	r := mux.NewStaticRouter(defs)
	fmt.Println("Server started at http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
