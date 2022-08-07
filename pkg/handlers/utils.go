package handlers

import (
	"fmt"
	"net/http"

	"git.sr.ht/~bossley9/feedme/pkg/atom"

	"github.com/gorilla/mux"
)

// errors

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

// usage

func HandleUsage(w http.ResponseWriter, r *http.Request, usage string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "usage: "+usage)
}

func HandleDefaultUsage(w http.ResponseWriter, r *http.Request) {
	usage := `/{type}?{param}={value}

available types are:
* acast
* gemini
* soundcloud
* @solene
`
	HandleUsage(w, r, usage)
}

// success

func HandleSuccess(w http.ResponseWriter, _ *http.Request, feed *atom.AtomFeed) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("Content-Disposition", "inline; filename=\"feed.xml\"")
	fmt.Fprintln(w, feed.String())
}
