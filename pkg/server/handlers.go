package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "usage: /{type}?{param}={value}")
}

func handleFeed(w http.ResponseWriter, r *http.Request) {
	feedType := mux.Vars(r)["type"]

	fmt.Fprintln(w, "type is '"+feedType+"'")
	fmt.Fprintln(w, "query string 'query' val is '"+r.FormValue("query")+"'")
}
