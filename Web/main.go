package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

const localhost = "localhost"
const localport = "3000"
const keycloakhost = "localhost"
const keycloakport = "8080"
const server = "http://" + localhost + ":" + localport
const keycloakserver = "http://" + keycloakhost + ":" + keycloakport
const realm = "demo"
const clientID = "mitchell"
const clientSecret = "6147e7de-67b9-423f-a0a5-4b79ef86e7cb"

var goTest bool
var oauthStateString = "random"

func main() {

	ctx := context.Background()
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

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":"+localport, router))
}
