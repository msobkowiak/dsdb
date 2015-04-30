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
		"Table",
		"GET",
		"/{table}",
		GetAllItems,
	},
	Route{
		"Table",
		"GET",
		"/{table}/primary/{hash}",
		GetByHash,
	},
	Route{
		"Table",
		"GET",
		"/{table}/primary/{hash}/{range}",
		GetByHashRange,
	},
	Route{
		"Table",
		"GET",
		"/{table}/primary/{hash}/{range}/{op:EQ|LE|LT|GE|GT|BEGINS_WITH}/{value}",
		GetByHashRangeOp,
	},
	Route{
		"Table",
		"GET",
		"/{table}/primary/{hash}/{range}/BETWEEN/{value1}/{value2}",
		GetByHashRangeBetween,
	},
}
