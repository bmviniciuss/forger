package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bmviniciuss/forger-golang/internal/core"
	"github.com/bmviniciuss/forger-golang/pkg/path"
)

type PostgresLoader struct {
	db *sql.DB
}

func NewPostgresLoader(db *sql.DB) *PostgresLoader {
	return &PostgresLoader{db}
}

// Ensures PostgresLoader implements core.Loader
var (
	_ core.Loader = (*PostgresLoader)(nil)
)

type dbRouteDefinition struct {
	Name               string `db:"name"`
	Path               string `db:"path"`
	Method             string `db:"method"`
	ResponseType       string `db:"response_type"`
	ResponseStatusCode int    `db:"response_status_code"`
	ResponseBody       string `db:"response_body"`
	ResponseHeaders    string `db:"response_headers"`
	ResponseDelay      int64  `db:"response_delay"`
}

const selectQuery = `
SELECT
	name, path, method,
	response_type, response_status_code, response_body,
	response_headers, response_delay
FROM forger.routes WHERE prefix = $1
`

func (l *PostgresLoader) Load(r *http.Request) ([]core.RouteDefinition, error) {
	var (
		ctx    = r.Context()
		prefix = path.ExtractPrefix(r.URL.Path)
	)
	rows, err := l.db.QueryContext(ctx, selectQuery, prefix)
	if err != nil {
		return []core.RouteDefinition{}, err
	}
	defer rows.Close()

	defs := []dbRouteDefinition{}
	for rows.Next() {
		var route dbRouteDefinition
		err = rows.Scan(
			&route.Name, &route.Path, &route.Method,
			&route.ResponseType, &route.ResponseStatusCode, &route.ResponseBody,
			&route.ResponseHeaders, &route.ResponseDelay,
		)
		if err != nil {
			return []core.RouteDefinition{}, err
		}
		defs = append(defs, route)
	}
	if err := rows.Err(); err != nil {
		return []core.RouteDefinition{}, err
	}

	routeDefs := make([]core.RouteDefinition, len(defs))
	for i, def := range defs {
		responseType, err := core.NewRouteResponseType(def.ResponseType)
		if err != nil {
			return []core.RouteDefinition{}, err
		}

		responseHeaders := make(map[string]string)
		if def.ResponseHeaders != "" {
			err = json.Unmarshal([]byte(def.ResponseHeaders), &responseHeaders)
			if err != nil {
				return []core.RouteDefinition{}, err
			}
		}

		response := core.NewRouteResponse(
			responseType,
			def.ResponseStatusCode,
			def.ResponseBody,
			responseHeaders,
			time.Duration(def.ResponseDelay)*time.Millisecond,
		)
		routeDefs[i] = *core.NewRouteDefinition(def.Path, def.Method, *response)
	}
	return routeDefs, nil
}
