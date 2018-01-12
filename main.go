package main

import (
	"log"
	"net/http"
)

func main() {

	router := NewRouter()
	go runUnitTests()
	log.Fatal(http.ListenAndServe(":8080", router))

}
