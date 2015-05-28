package main

import (
	"log"
	"net/http"
)

var schema DbDescription

func main() {
	schema = LoadSchema("config.yml")

	Bootstrap(schema)
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
