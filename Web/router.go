package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

//StaticDir is the path to the static files hosted
const StaticDir = "/static/"

//NewRouter creates a new router and maps all the routes defined in routes.go
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.
		PathPrefix(StaticDir).
		Handler(http.StripPrefix(StaticDir, http.FileServer(http.Dir("."+StaticDir))))
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
