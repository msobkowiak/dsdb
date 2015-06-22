package main

import (
	"log"
	"net/http"

	"fmt"
)

var schema DbDescription

func main() {
	schema = LoadSchema("config.yml")

	Bootstrap(schema)

	fmt.Println("Done creating tables...")

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
