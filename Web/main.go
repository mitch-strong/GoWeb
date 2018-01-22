package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var goTest bool
var oauthStateString = "random"

func main() {

	ctx := context.Background()
	//fmt.Printf("%v", ctx)
	provider, err = oidc.NewProvider(ctx, "http://localhost:8080/auth/realms/demo")
	if err != nil {
		fmt.Printf("This is an error with regard to the context: %v", err)
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: "mitchell"})

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config = oauth2.Config{
		ClientID:     "mitchell",
		ClientSecret: "6147e7de-67b9-423f-a0a5-4b79ef86e7cb",
		//RedirectURL:  "http://localhost:8080/mitchell/",
		RedirectURL: "http://localhost:3000/loginCallback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":3000", router))
}
