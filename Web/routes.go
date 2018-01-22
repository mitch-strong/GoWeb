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
		Index,
	},
	Route{
		"PersonList",
		"GET",
		"/people",
		PersonList,
	},
	Route{
		"PersonListJSON",
		"GET",
		"/peopleJSON",
		PersonListJSON,
	},
	Route{
		"PersonCreate",
		"POST",
		"/people",
		PersonCreate,
	},
	Route{
		"GenericJSON",
		"POST",
		"/JSON",
		GenericJSON,
	},
	Route{
		"GenericListJSON",
		"GET",
		"/JSON",
		GenericListJSON,
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
}
