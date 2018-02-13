package main

import (
	"flag"
	"log"
	"net/http"

	keycloak "github.com/mitch-strong/keycloakgo"
)

//hosting variables defined in flags
var localhost string
var localport string
var keycloakhost string
var keycloakport string
var server string
var keycloakserver string

var goTest bool // true if unit tests are running

func main() {
	flag.StringVar(&localport, "p", "3000", "Specify which port to use")
	flag.StringVar(&localhost, "host", "localhost", "Specify the name of the host")
	flag.StringVar(&keycloakhost, "keycloakh", "localhost", "Specify the name of the keycloak host")
	flag.StringVar(&keycloakport, "keycloakp", "8080", "Specify the port of keycloak")
	flag.Parse()

	server = "http://" + localhost + ":" + localport
	keycloakserver = "http://" + keycloakhost + ":" + keycloakport

	addKeycloak(keycloakserver, server)

	router := NewRouter()
<<<<<<< HEAD
	log.Fatal(http.ListenAndServe(":8000", router))
=======
	//Stats hosting on the constant port
	log.Fatal(http.ListenAndServe(":"+localport, router))
}

func addKeycloak(keycloakserver, server string) {
	keycloak.Init(keycloakserver, server)
>>>>>>> 000513c4a31478cf641b6479f84a0ee2c1f91894
}
