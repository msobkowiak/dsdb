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
		"GetByHash",
		"DELETE",
		"/{table}/{hash}",
		DeleteByHash,
	},
	Route{
		"GetByHashRange",
		"GET",
		"/{table}/{hash}/{range}",
		GetByHashRange,
	},
	Route{
		"FilterByRange",
		"GET",
		"/{table}/primary/{range}/{op:EQ|LE|LT|GE|GT|BEGINS_WITH}/{value}",
		GetByRange,
	},
	Route{
		"FilterByRange",
		"GET",
		"/{table}/primary/{range}/BETWEEN/{value1}/{value2}",
		GetByHashRangeBetween,
	},
	/*Route{
		"AddItem",
		"POST",
		"/{table}",
		AddItem,
	},*/
}
