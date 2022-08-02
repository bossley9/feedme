package main

var DOMAIN_NAME string
var PORT string
var CERT_FILE string
var KEY_FILE string

func GetEnv(key string, fallback string) string {
	if len(key) == 0 {
		return fallback
	}
	return key
}
