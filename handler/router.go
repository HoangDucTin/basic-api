package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/tinwoan-go/basic-api/handler/check"
)

// This function returns an example
// handler for your service with
// an echo function to check.
// It also provide the applying of
// some self-built middleware to
// set no-cache header and print
// out the request and response of
// each request in JSON format
// (suitable for elastic search).
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(SetNoCacheHeader)
	r.Use(NewLogMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Get("/status", check.Status())
	})
	return r
}

// End-of-file
