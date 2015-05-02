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
		"GetAll",
		"GET",
		"/{table}",
		GetAllItems,
	},
	Route{
		"GetByHash",
		"GET",
		"/{table}/{hash}",
		GetByHash,
	},
	Route{
		"GetByHashRange",
		"GET",
		"/{table}/{hash}/{range}",
		GetByHashRange,
	},
	// search route goes here
	Route{
		"DeleteItem",
		"DELETE",
		"/{table}/{hash}",
		DeleteByHash,
	},
	Route{
		"DeleteItems",
		"DELETE",
		"/{table}/{hash}/{range}",
		DeleteByHashRange,
	},
	Route{
		"AddItem",
		"POST",
		"/{table}/{hash}",
		AddItem,
	},
	Route{
		"AddItem",
		"POST",
		"/{table}/{hash}/{range}",
		AddItemHashRange,
	},
	/*Route{
		"AddItem",
		"POST",
		"/{table}",
		AddItem,
	},*/
}
