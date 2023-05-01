package server

import (
	"fmt"
	"net/http"
	"time"

	h "github.com/bossley9/feedme/pkg/handlers"
)

func New(domain string, port string, certFile string, keyFile string) error {
	r := h.SetupRouter()
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
