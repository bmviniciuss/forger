package core

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/bmviniciuss/forger/core/generators"
	"github.com/go-chi/chi/v5"
	"github.com/tidwall/gjson"
)

func processString(r *http.Request, src string, reqBody *string) (*string, error) {
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
			"requestBody": func(params ...interface{}) (interface{}, error) {
				if len(params) == 0 {
					return *reqBody, nil
				}
				if len(params) > 1 {
					return nil, ErrInvalidRequestBodyKeysAmount
				}
				key, ok := params[0].(string)
				if !ok {
					return nil, ErrInvalidRequestBodyArgumentType
				}
				res := gjson.Get(*reqBody, key)
				if !res.Exists() {
					return "null", nil
				}
				return res.Raw, nil
			},
		}).
		Parse(src)
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
