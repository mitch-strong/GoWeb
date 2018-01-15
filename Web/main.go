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

//How to implement generic interface for random JSON objects
//https://blog.golang.org/json-and-go
