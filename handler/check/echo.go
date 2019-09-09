package check

import (
	"github.com/go-chi/render"
	"net/http"
)

type echo struct {
	Message string `json:"message"`
}

// This handler responds a message
// to notice the user that everything
// works fine.
func Echo() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, echo{
			Message: "Welcome to my humble library!",
		})
	})
}
