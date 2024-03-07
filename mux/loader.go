package mux

import (
	"net/http"

	"github.com/bmviniciuss/forger-golang/internal/core"
)

type Loader interface {
	Load(r *http.Request) ([]core.RouteDefinition, error)
}
