package check

import (
	"github.com/go-chi/render"
	"net/http"
)

type echoFailed struct {
	Failure string `json:"failure"`
}

func Echo(w http.ResponseWriter, r *http.Request) {
	var request interface{}

	// Responsed with error (if any)
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, echoFailed{
			Failure: err.Error(),
		})
	}

	// Echo back what we got
	render.Status(r, http.StatusOK)
	render.JSON(w, r, request)
}

// End-of-file
