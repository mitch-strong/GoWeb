package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

//Global vairable definitions
var oauth2Config oauth2.Config
var provider *oidc.Provider
var err error
var verifier *oidc.IDTokenVerifier
var token *oauth2.Token

//Index returns when the main page is called and returns HTML indicating the availale paths
var indexHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("./templates/index.html")
	tpl.Execute(w, nil)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return
})

//PersonList returns a readable list of people
var personListHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, _ := template.ParseFiles("./templates/people.html")

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	data := struct {
		Title      string
		PeopleList []Person
	}{
		Title:      "List of People",
		PeopleList: people,
	}
	err := tpl.Execute(w, data)
	check(err)
	return
})

//PersonListJSON returns a comma seperated list of people as raw JSON
var personListJSONHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(people); err != nil {
		panic(err)
	}
})

//GenericListJSON returns a comma seperated list of generic JSON objects stored in objects array
var genericListJSONHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(objects); err != nil {
		panic(err)
	}
})

//PersonCreate is a POST method that converts JSON objects into people objects and stores them
var personCreateHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var person Person
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &person); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreatePerson(person)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
})

//GenericJSON is a POST method that converts JSON objects into empty interface objects and stores them
var genericJSONHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var f interface{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &f); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	m := f.(map[string]interface{})
	objects = append(objects, m)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
})

//handleLogin is the login function, redirects to the loginCallback function
func handleLogin(w http.ResponseWriter, r *http.Request) {
	//create a random string for oath2 verification
	oauthStateString = randSeq(20)
	//Uses random gnerated string to verify keyclock security
	url := oauth2Config.AuthCodeURL(oauthStateString)
	//redirects to loginCallback
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//handleLoginCallback is a fuction that verifies login success and forwards to index
func handleLoginCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	//Checks that the strings are in a consistent state
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	//Gets the code from keycloak
	code := r.FormValue("code")
	//Exchanges code for token
	token, err = oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%v'\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	client := &http.Client{}
	url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	//Sends the token to get user info
	response, err := client.Do(req)
	//Checks if token and authentication were successful
	if err != nil || response.Status == "401 Unauthorized" {
		//forwards back to login if not successful
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else {
		//forwards to index if login sucessful
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	return
}

//authMiddleware is a middlefuntion that verifies authentication before each redirect
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//If running unit tests skip authentication (temp)
		if goTest {
			next.ServeHTTP(w, r)
		}
		client := &http.Client{}
		url := keycloakserver + "/auth/realms/" + realm + "/protocol/openid-connect/userinfo"
		req, _ := http.NewRequest("GET", url, nil)
		if token == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		//Check if token is still valid
		response, err := client.Do(req)
		if err != nil || response.Status != "200 OK" {
			//Go to login if token is no longer valid
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		} else {
			//Go to redirect if token is still valid
			next.ServeHTTP(w, r)
		}
	})
	//return function for page handling
	return handler
}

//Logout logs the user out
func Logout(w http.ResponseWriter, r *http.Request) {
	//Makes the logout page redirect to login page
	URI := server + "/login"
	//Logout using endpoint and redirect to login page
	http.Redirect(w, r, keycloakserver+"/auth/realms/"+realm+"/protocol/openid-connect/logout?redirect_uri="+URI, http.StatusTemporaryRedirect)

}

//randSeq generates a random string of letters of the given length (Helper function)
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
