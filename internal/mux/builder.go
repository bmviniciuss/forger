package mux

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bmviniciuss/forger-golang/internal/core"
	"github.com/bmviniciuss/forger-golang/internal/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const (
	utcLayout = "2006-01-02T15:04:05.000Z"
)

func requestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := getReqID(r)
		w.Header().Set("request-id", reqID)
		h.ServeHTTP(w, r.WithContext(ctx.WithRequestID(r.Context(), reqID)))
	})
}

func getReqID(r *http.Request) string {
	reqID := r.Header.Get("request-id")
	if reqID != "" {
		return reqID
	}
	return uuid.NewString()
}

func startTime(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-forger-req-start", time.Now().Format(utcLayout))
		h.ServeHTTP(w, r)
	})
}

func NewStaticRouter(defs []core.RouteDefinition) *chi.Mux {
	r := chi.NewRouter()
	r.Use(startTime)
	r.Use(requestID)
	r.Use(middleware.Logger)
	registerRoutes(r, defs)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("x-forger-req-end", time.Now().Format(utcLayout))
		render.JSON(w, r, NewNotFoundResponse())
	})
	return r
}

func NewDynamicRouter(providerFn func(r *http.Request) ([]core.RouteDefinition, error)) *chi.Mux {
	router := chi.NewRouter()
	router.Use(startTime)
	router.Use(requestID)
	router.Use(middleware.Logger)
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		defs, err := providerFn(r)
		if err != nil {
			render.JSON(w, r, NewInternalErrorResponse("Internal Server Error", "Error while getting route definitions from provider"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		subRouter := chi.NewRouter()
		registerRoutes(subRouter, defs)
		subRouter.ServeHTTP(w, r)
	})
	return router
}

func registerRoutes(r *chi.Mux, defs []core.RouteDefinition) {
	for _, def := range defs {
		fmt.Println("Registering route", def.Method, def.Path)
		r.Method(def.Method, def.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := def.Response.BuildResponseBody(r)
			if err != nil {
				render.JSON(w, r, NewInternalErrorResponse("Internal Server Error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resCode := def.Response.BuildResponseStatusCode()
			resHeaders, err := def.Response.BuildHeaders(r)
			if err != nil {
				render.JSON(w, r, NewInternalErrorResponse("Internal Server Error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
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
