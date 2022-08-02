package handlers

import (
	"fmt"
	"net/http"
)

func HandleUsage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "usage: /{type}?{param}={value}")
}
