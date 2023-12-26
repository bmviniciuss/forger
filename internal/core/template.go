package core

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/bmviniciuss/forger-golang/internal/core/generators"
	"github.com/go-chi/chi/v5"
)

func ProcessString(r *http.Request, rr RouteResponse) (*string, error) {
	t, err := template.New("").
		Funcs(template.FuncMap{
			"uuid": func(options ...interface{}) (string, error) {
				return generators.UUID(r.Context(), options...)
			},
			"requestVar": func(name string) string {
				val := chi.URLParam(r, name)
				return val
			},
			"requestHeader": func(key string) string {
				return r.Header.Get(key)
			},
			"requestQuery": func(key string) string {
				return r.URL.Query().Get(key)
			},
			"time": func(options ...interface{}) (string, error) {
				return generators.Time(r.Context(), options...)
			},
		}).
		Parse(rr.Body)
	if err != nil {
		return nil, err
	}

	builder := &strings.Builder{}
	err = t.Execute(builder, r)
	if err != nil {
		return nil, err
	}
	result := builder.String()
	return &result, nil
}
