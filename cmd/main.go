package main

import (
	"log"

	"git.sr.ht/~bossley9/feedme/pkg/server"
)

func main() {
	domain := GetEnv(DOMAIN_NAME, "localhost")
	port := GetEnv(PORT, "9000")
	certFile := GetEnv(CERT_FILE, "")
	keyFile := GetEnv(KEY_FILE, "")

	log.Fatal(server.New(domain, port, certFile, keyFile))
}
