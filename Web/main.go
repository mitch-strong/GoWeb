package main

import (
	"log"
	"net/http"
)

var goTest bool

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
