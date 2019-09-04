package check

import (
	"github.com/HoangDucTin/basic-api/internal/mongo"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/go-chi/render"

)

type viewResponse struct {
	Count int    `json:"count" bson:"count"`
	Error string `json:"error" bson:"-"`
}

func countView() (result viewResponse) {
	err := mongo.Change(
		"test",
		"counter",
		bson.M{"name": "views"},
		bson.M{"$inc": bson.M{"count": 1}},
		&result)

	if err != nil {
		result.Error = err.Error()
	}

	return result
}

func View(w http.ResponseWriter, r *http.Request) {

	result := countView()
	if result.Error != "" {
		render.Status(r, http.StatusBadRequest)
	} else {
		render.Status(r, http.StatusOK)
	}

	render.JSON(w, r, result)
}

// End-of-file
