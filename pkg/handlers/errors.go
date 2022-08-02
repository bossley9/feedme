package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	feedType := mux.Vars(r)["type"]
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "feed of type '%s' not found.\n", feedType)
}

func HandleUnimplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Not yet implemented.")
}

func HandleBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, err)
}
