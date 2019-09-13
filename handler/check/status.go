package check

import (
	"github.com/go-chi/render"
	"net/http"
)

// This HandlerFunc will response
// HTTP status 200 for the purpose
// of health checking.
func Status() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		render.Status(r, http.StatusOK)
	})
}