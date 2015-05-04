package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

		/*if route.Queries != nil {
			addQueryStringMatch(router, route.Queries)
		}*/

	}
	return router
}

/*func addQueryStringMatch(router *mux.Router, queries Queries) {
	for i := range queries {
		router.Queries(queries[i].Key, queries[i].Value)
	}
}*/
