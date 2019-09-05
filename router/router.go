package router

import (
	"github.com/basic-api/router/check"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Suggested by the security team
func setNoCacheHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate") // HTTP 1.1
		w.Header().Set("Pragma", "no-cache")                                 // HTTP 1.0
		w.Header().Set("Expires", "0")                                       // Proxies
		next.ServeHTTP(w, r)
	})
}

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(setNoCacheHeader)

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
