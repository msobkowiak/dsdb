package main

import (
	"log"
	"net/http"
)

var schema DbDescription

// func LoadSchema() DbDescription {
// 	return dbDescription
// }

func main() {
	schema = LoadSchema()
	Bootstrap()
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
