package server

import (
	"fmt"
	"net/http"
	"time"

	h "git.sr.ht/~bossley9/feedme/pkg/handlers"

	"github.com/gorilla/mux"
)

func New(domain string, port string, certFile string, keyFile string) error {
	r := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	r.HandleFunc("/", h.HandleDefaultUsage)
	r.HandleFunc("/{type}", handleFeed)
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         domain + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      r,
	}

	fmt.Println("starting server...")
	if len(certFile) > 0 && len(keyFile) > 0 {
		return srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		return srv.ListenAndServe()
	}
}
