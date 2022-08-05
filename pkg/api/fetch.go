package api

import (
	"io/ioutil"
	"net/http"
)

func FetchGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
