package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		TodoShow,
	},
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		TodoCreate,
	},
	Route{
		"TableCreate",
		"POST",
		"/create",
		TableCreate,
	},
	Route{
		"TablesIndex",
		"GET",
		"/tables",
		TablesIndex,
	},
	Route{
		"Table",
		"GET",
		"/{table}/primary/{hash}",
		TableGetByHash,
	},
}
