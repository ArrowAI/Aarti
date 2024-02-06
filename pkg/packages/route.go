package packages

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerFunc func(repo string) http.HandlerFunc

type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func makeHandler(repo string, h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if repo == "" {
			repo = mux.Vars(r)["repo"]
		}
		h(repo)(w, r)
	}
}
