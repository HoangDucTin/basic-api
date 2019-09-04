package check

import (
	"net/http"

	"github.com/go-chi/render"
)

type infoResponse struct {
	Version string `json:"version"`
}

func Info(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, infoResponse{
		Version: "0.0.1",
	})
}

// End-of-file
