package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/tinwoan-go/basic-api/handler/check"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(setNoCacheHeader)
	r.Use(NewLogMiddleware(logrus.New()))

	r.Route("/", func(r chi.Router) {
		r.Route("/check", func(r chi.Router) {
			r.Post("/echo", check.Echo)
			r.Get("/info", check.Info)
			r.Get("/view", check.View)
		})
	})

	return r
}

// End-of-file
