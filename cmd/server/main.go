package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bmviniciuss/forger-golang/internal/core"
	"github.com/bmviniciuss/forger-golang/internal/mux"
)

func main() {
	defs := []core.RouteDefinition{
		{
			Path:   "/item",
			Method: "GET",
			Response: core.RouteResponse{
				Type:       core.RESPONSE_TYPE_STATIC,
				StatusCode: 200,
				Body:       "{id: 1, name: 'Item 1'}",
				Headers:    map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			Path:   "/item/{id}",
			Method: "GET",
			Response: core.RouteResponse{
				Type:       core.RESPONSE_TYPE_DYNAMIC,
				StatusCode: 200,
				Body: `{
					"id": "$requestVar('id')", 
					"page": "$requestParameter('page')",
					"client_id": "$requestHeader('client-id')",
					"random_uuid": "$uuid"
				}`,
				Delay: time.Duration(0) * time.Millisecond,
				Headers: map[string]string{
					"Content-Type":  "application/json",
					"Item-ID":       "$requestVar('id')",
					"page":          "$requestParameter('page')",
					"res-client-id": "$requestHeader('client-id')",
				},
			},
		},
	}

	r := mux.New(defs)
	fmt.Println("Server started at http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
