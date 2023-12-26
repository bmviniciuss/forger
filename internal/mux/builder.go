package mux

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bmviniciuss/forger-golang/internal/core"
	"github.com/bmviniciuss/forger-golang/internal/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func New(defs []core.RouteDefinition) *chi.Mux {
	r := chi.NewRouter()
	r.Use(startTime)
	r.Use(requestID)
	r.Use(middleware.Logger)

	for _, def := range defs {
		fmt.Println("Registering route", def.Method, def.Path)
		r.Method(def.Method, def.Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := def.Response.BuildResponseBody(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			resCode := def.Response.BuildResponseStatusCode()
			resHeaders, err := def.Response.BuildHeaders(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
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

	return r
}
