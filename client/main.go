package main

import (
	"io"
	"net/http"
)

func main() {
	req, err := http.Get("http://server-go-client-server-1:8081/")
	if err != nil {
		panic(err)
	}
	_, err = io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	req.Body.Close()
}
