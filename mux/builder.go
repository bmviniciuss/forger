package mux

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmviniciuss/forger-golang/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	utcLayout = "2006-01-02T15:04:05.000Z"
)

func NewStaticRouter(defs []core.RouteDefinition) *chi.Mux {
	router := chi.NewRouter()
	setMiddlewares(router)
	registerRoutes(router, defs)
	setNotFoundHandler(router)
	return router
}

func NewDynamicRouter(loader Loader) *chi.Mux {
	router := chi.NewRouter()
	setMiddlewares(router)
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		defs, err := loader.Load(r)
		if err != nil {
			log.Default().Printf("Error while getting route definitions from provider: %s", err)
			render.JSON(w, r, newInternalErrorResponse("Internal Server Error", "Error while getting route definitions from provider"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Loaded route definitions: %+v\n\n", defs)
		subRouter := chi.NewRouter()
		registerRoutes(subRouter, defs)
		subRouter.ServeHTTP(w, r)
	})
	setNotFoundHandler(router)
	return router
}

func registerRoutes(router *chi.Mux, defs []core.RouteDefinition) {
	baseHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	for _, route := range defs {
		def := route
		fmt.Printf("Registering route [%+v]\n\n", def)
		router.Method(def.Method, def.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Handling route %+v\n\n", def)
			body, err := def.Response.BuildResponseBody(r)
			if err != nil {
				render.JSON(w, r, newInternalErrorResponse("Internal Server Error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resCode := def.Response.BuildResponseStatusCode()
			resHeaders, err := def.Response.BuildHeaders(r)
			if err != nil {
				render.JSON(w, r, newInternalErrorResponse("Internal Server Error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			for k, v := range baseHeaders {
				w.Header().Set(k, v)
			}
			for k, v := range resHeaders {
				w.Header().Set(k, v)
			}
			if def.Response.Delay > 0 {
				time.Sleep(def.Response.Delay)
			}
			w.Header().Set("x-forger-req-end", time.Now().Format(utcLayout))
			w.WriteHeader(resCode)
			w.Write([]byte(*body))
		}))
	}
}

func setNotFoundHandler(router *chi.Mux) {
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-forger-req-end", time.Now().Format(utcLayout))
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, newNotFoundResponse())
	})
}
