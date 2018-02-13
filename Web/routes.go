package main

import "net/http"

//Route object creates to keep track of routes for router
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes is an array of Route objects
type Routes []Route

var routes = Routes{
	//Home page
	///Authenticated
	Route{
		"Index",
		"GET",
		"/",
		authMiddleware(indexHandler),
	},
	//HTML list of people
	///Authenticated
	Route{
		"PersonList",
		"GET",
		"/people",
		authMiddleware(personListHandler),
	},
	//JSON list of people
	///Authenticated
	Route{
		"PersonListJSON",
		"GET",
		"/peopleJSON",
		authMiddleware(personListJSONHandler),
	},
	//Create new peoson POST
	///Authenticated
	Route{
		"PersonCreate",
		"POST",
		"/people",
		authMiddleware(personCreateHandler),
	},
	//Create new generic JSON object
	///Authenticated
	Route{
		"GenericJSON",
		"POST",
		"/JSON",
		authMiddleware(genericJSONHandler),
	},
	//List of JSON objects
	///Authenticated
	Route{
		"GenericListJSON",
		"GET",
		"/JSON",
		authMiddleware(genericListJSONHandler),
	},
	//Login page
	///Unauthenticated
	Route{
		"handleLogin",
		"GET",
		"/login",
		handleLogin,
	},
	//Login helper
	//Authenticated
	Route{
		"handleLoginCallback",
		"GET",
		"/loginCallback",
		handleLoginCallback,
	},
	//Logout, redirects to login
	///Unauthenticatec
	Route{
		"logout",
		"GET",
		"/logout",
		logout,
	},
}
