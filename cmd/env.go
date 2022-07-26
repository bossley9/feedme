package main

var DOMAIN_NAME string

func GetEnv(key string, fallback string) string {
	if len(key) == 0 {
		return fallback
	}
	return key
}
