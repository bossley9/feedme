package main

import (
	"flag"
	"log"

	"git.sr.ht/~bossley9/feedme/pkg/server"
)

func main() {
	var domain, port, certFile, keyFile string

	flag.StringVar(&domain, "d", "localhost", "server domain name")
	flag.StringVar(&port, "p", "9000", "server port")
	flag.StringVar(&certFile, "c", "", "TLS certificate file")
	flag.StringVar(&certFile, "k", "", "TLS key file")
	flag.Parse()

	log.Fatal(server.New(domain, port, certFile, keyFile))
}
