package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	keycloak "github.com/mitch-strong/keycloakgo"
)

//Global vairable definitions
var err error

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
	keycloak.HandleLogin(w, r)
	return
}

//handleLoginCallback is a fuction that verifies login success and forwards to index
func handleLoginCallback(w http.ResponseWriter, r *http.Request) {
	keycloak.HandleLoginCallback(w, r)
	return
}

//authMiddleware is a middlefuntion that verifies authentication before each redirect
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return keycloak.AuthMiddleware(next)
}

//logout logs the user out
func logout(w http.ResponseWriter, r *http.Request) {
	keycloak.Logout(w, r)
}
