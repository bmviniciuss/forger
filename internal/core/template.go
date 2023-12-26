package core

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/bmviniciuss/forger-golang/internal/core/generators"
	"github.com/go-chi/chi/v5"
)

func BuildBody(r *http.Request, rr RouteResponse) (*string, error) {
	t, err := template.New("").
		Funcs(template.FuncMap{
			"uuid": func(options ...interface{}) (string, error) {
				return generators.UUID(r.Context(), options...)
			},
			"requestVar": func(name string) string {
				val := chi.URLParam(r, name)
				return val
			},
		}).
		Parse(rr.Body)
	if err != nil {
		fmt.Println("err", err)

		return nil, err
	}

	builder := &strings.Builder{}
	err = t.Execute(builder, r)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	result := builder.String()
	return &result, nil
}
