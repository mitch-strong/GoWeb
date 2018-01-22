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
	Route{
		"Index",
		"GET",
		"/",
		authMiddleware(indexHandler),
	},
	Route{
		"PersonList",
		"GET",
		"/people",
		authMiddleware(personListHandler),
	},
	Route{
		"PersonListJSON",
		"GET",
		"/peopleJSON",
		authMiddleware(personListJSONHandler),
	},
	Route{
		"PersonCreate",
		"POST",
		"/people",
		authMiddleware(personCreateHandler),
	},
	Route{
		"GenericJSON",
		"POST",
		"/JSON",
		authMiddleware(genericJSONHandler),
	},
	Route{
		"GenericListJSON",
		"GET",
		"/JSON",
		authMiddleware(genericListJSONHandler),
	},
	Route{
		"handleLogin",
		"GET",
		"/login",
		handleLogin,
	},
	Route{
		"handleLoginCallback",
		"GET",
		"/loginCallback",
		handleLoginCallback,
	},
	// 	Route{
	// 		"logout",
	// 		"GET",
	// 		"/logout",
	// 		Logout,
	// 	},
}
