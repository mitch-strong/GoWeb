package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var oauth2Config oauth2.Config
var provider *oidc.Provider
var err error
var verifier *oidc.IDTokenVerifier

//Index returns when the main page is called and returns HTML indicating the availale paths
func Index(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("./templates/index.html")
	tpl.Execute(w, nil)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	oauth2Token, _ := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	http.Redirect(w, r, oauth2Config.AuthCodeURL(oauth2Token.AccessToken), http.StatusFound)
}

func handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// Verify state and errors
	ctx := context.Background()

	oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		// handle error
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// handle missing token
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		// handle error
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		// handle error
	}
}

//PersonList returns a readable list of people
func PersonList(w http.ResponseWriter, r *http.Request) {
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
}

//PersonListJSON returns a comma seperated list of people as raw JSON
func PersonListJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(people); err != nil {
		panic(err)
	}
}

//GenericListJSON returns a comma seperated list of generic JSON objects stored in objects array
func GenericListJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(objects); err != nil {
		panic(err)
	}
}

//PersonCreate is a POST method that converts JSON objects into people objects and stores them
func PersonCreate(w http.ResponseWriter, r *http.Request) {
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
}

//GenericJSON is a POST method that converts JSON objects into empty interface objects and stores them
func GenericJSON(w http.ResponseWriter, r *http.Request) {
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

	// for k, v := range m {
	// 	switch vv := v.(type) {
	// 	case string:
	// 		fmt.Println(k, "is string", vv)
	// 	case float64:
	// 		fmt.Println(k, "is float64", vv)
	// 	case []interface{}:
	// 		fmt.Println(k, "is an array:")
	// 		for i, u := range vv {
	// 			fmt.Println(i, u)
	// 		}
	// 	default:
	// 		fmt.Println(k, "is of a type I don't know how to handle")
	// 	}
	// }

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleLoginCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Made it into the callback")
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	client := &http.Client{}
	url := "http://localhost:8080/auth/realms/demo/protocol/openid-connect/userinfo"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	response, err := client.Do(req)
	fmt.Printf(response.Status)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)
}
