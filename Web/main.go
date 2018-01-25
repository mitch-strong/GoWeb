package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"flag"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// RDB: usage should be displayed when someone runs the program without arguments
func usage() {
	// Add a usage statement here, this will be used with the passed arguments
}

// RDB: Add these as command line arguments to the program, look at using the "flag" library
//      Specify required and optional arguments (i.e. localhost is hardcoded, why?  I'm on an older
//      version of docker which uses docker-machine and not localhost).
//Constants for server information and client information
const localhost = "localhost"
const localport = "3000"
const keycloakhost = "keycloak"
const keycloakport = "8080"
const server = "http://" + localhost + ":" + localport
const keycloakserver = "http://" + keycloakhost + ":" + keycloakport
const realm = "demo"
const clientID = "mitchell"
const clientSecret = "6147e7de-67b9-423f-a0a5-4b79ef86e7cb"

var goTest bool // true if unit tests are running
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var oauthStateString string

func main() {
	ctx := context.Background()
	//Gets the provider for authentication (keycloak)
	provider, err = oidc.NewProvider(ctx, keycloakserver+"/auth/realms/"+realm)
	if err != nil {
		fmt.Printf("This is an error with regard to the context: %v", err)
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  server + "/loginCallback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	router := NewRouter()
	//Stats hosting on the constant port
	log.Fatal(http.ListenAndServe(":"+localport, router))
}
