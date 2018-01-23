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
})

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleLoginCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
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
	response, err := client.Do(req)
	if err != nil || response.Status == "401 Unauthorized" {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	return
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		response, err := client.Do(req)
		if err != nil || response.Status != "200 OK" {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		} else {
			next.ServeHTTP(w, r)
		}
	})

	return handler
}

//Logout logs the user out
func Logout(w http.ResponseWriter, r *http.Request) {
	URI := server + "/login"
	http.Redirect(w, r, keycloakserver+"/auth/realms/"+realm+"/protocol/openid-connect/logout?redirect_uri="+URI, http.StatusTemporaryRedirect)

}
