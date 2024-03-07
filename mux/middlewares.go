package mux

import (
	"net/http"
	"time"

	"github.com/bmviniciuss/forger-golang/internal/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func setMiddlewares(router *chi.Mux) {
	router.Use(startTime)
	router.Use(requestID)
	router.Use(middleware.Logger)
}

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
