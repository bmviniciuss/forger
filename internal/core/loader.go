package core

import (
	"net/http"
)

type Loader interface {
	Load(r *http.Request) ([]RouteDefinition, error)
}
