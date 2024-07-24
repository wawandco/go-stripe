package slices

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Show(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	err := rw.Render("slices/slice_one.slide")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
