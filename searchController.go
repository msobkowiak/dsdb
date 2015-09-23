package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryParams := r.URL.Query()

	table := vars["table"]
	searchType := getValue(queryParams["search_type"])

	var elasticSearchRepo ElasticSeaarchRepository

	switch searchType {
	case "index":

		index := queryParams["index"]
		hashKey := queryParams["hash"]
		rangeOperator := queryParams["range_operator"]
		rangeValue := queryParams["range_value"]

		if index[0] == "primary" {
			primaryKeySearch(rangeOperator, hashKey, rangeValue, table, w)
		} else {
			secondaryIndexSearch(index, rangeOperator, hashKey, rangeValue, table, w)
		}
	case "text":
		field := getValue(queryParams["field"])
		query := getValue(queryParams["query"])

		if field != "" && query != "" {
			data, err := elasticSearchRepo.FullTextSearchQuery(table, field, query, getValue(queryParams["operator"]), getValue(queryParams["precision"]))
			writeResponse(data, err, w)
		} else {
			writeErrorResponse("Missing search parameters", 404, w)
		}
	case "faced":
		field := getValue(queryParams["field"])
		metric := getValue(queryParams["metric"])

		fmt.Println(field, metric)
		data, err := elasticSearchRepo.AggregationSearch(table, field, metric)
		writeResponse(data, err, w)
	case "geo":
		field := getValue(queryParams["field"])
		distance := getValue(queryParams["distance"])
		lat := getValue(queryParams["lat"])
		latValue, _ := strconv.ParseFloat(lat, 64)
		lon := getValue(queryParams["lon"])
		lonValue, _ := strconv.ParseFloat(lon, 64)

		if field != "" && distance != "" {
			data, err := elasticSearchRepo.GeoSearch(table, field, distance, latValue, lonValue)
			fmt.Println(data)
			writeResponse(data, err, w)
		} else {
			writeErrorResponse("Missing search parameters", 404, w)
		}
	}
}

func primaryKeySearch(rangeOperator, hashKey, rangeValue []string, table string, w http.ResponseWriter) {
	if rangeOperator == nil {
		if hashKey != nil {
			getItemsByHash(table, hashKey[0], w)
		} else {
			writeErrorResponse("Missing hash value", 404, w)
		}
	} else if hashKey != nil && rangeValue != nil {
		var repo DynamoBaseRepository
		data, err := repo.GetByOnlyRange(table, hashKey[0], rangeOperator[0], rangeValue)
		writeResponse(data, err, w)
	} else {
		writeErrorResponse("Missing primary key value(s)", 404, w)
	}
}

func secondaryIndexSearch(index, rangeOperator, hashKey, rangeValue []string, table string, w http.ResponseWriter) {
	if rangeOperator == nil {
		if hashKey != nil {
			getItemsByIndexHash(table, index[0], hashKey[0], w)
		} else {
			writeErrorResponse("Missing hash value", 404, w)
		}
	} else if hashKey != nil && rangeValue != nil {
		var repo DynamoBaseRepository
		data, err := repo.GetByIndexRange(table, index[0], hashKey[0], rangeOperator[0], rangeValue)
		writeResponse(data, err, w)
	} else {
		writeErrorResponse("Missing primary key value(s)", 404, w)
	}
}
