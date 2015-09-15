package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	var repo DynamoCRUDRepository
	data, err := repo.GetAll(table)
	writeResponse(data, err, w)
}

func GetByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	table := vars["table"]

	getItemsByHash(table, hash, w)
}

func GetByHashRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	rangeKey := vars["range"]
	table := vars["table"]

	var repo DynamoCRUDRepository
	data, err := repo.GetByHashRange(table, hash, rangeKey)
	writeResponse(data, err, w)
}

func GetByRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rangeKey := vars["range"]
	table := vars["table"]

	var repo DynamoCRUDRepository
	data, err := repo.GetByOnlyHash(table, rangeKey)
	writeResponse(data, err, w)
}

func Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryParams := r.URL.Query()

	table := vars["table"]
	searchType := getValue(queryParams["search_type"])

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
			data, err := FullTextSearchQuery(table, field, query, getValue(queryParams["operator"]), getValue(queryParams["precision"]))
			writeResponse(data, err, w)
		} else {
			writeErrorResponse("Missing search parameters", 404, w)
		}
	case "faced":
		field := getValue(queryParams["field"])
		metric := getValue(queryParams["metric"])

		fmt.Println(field, metric)
		data, err := AggregationSearch(table, field, metric)
		writeResponse(data, err, w)
	case "geo":
		field := getValue(queryParams["field"])
		distance := getValue(queryParams["distance"])
		lat := getValue(queryParams["lat"])
		latValue, _ := strconv.ParseFloat(lat, 64)
		lon := getValue(queryParams["lon"])
		lonValue, _ := strconv.ParseFloat(lon, 64)

		if field != "" && distance != "" {
			data, err := GeoSearch(table, field, distance, latValue, lonValue)
			writeResponse(data, err, w)
		} else {
			writeErrorResponse("Missing search parameters", 404, w)
		}
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashKey := vars["hash"]
	rangeKey := vars["range"]
	table := vars["table"]

	var ok bool
	var err error
	var dynamoRepo DynamoCRUDRepository
	if rangeKey != "" {
		ok, err = dynamoRepo.DeleteByHashRange(table, hashKey, rangeKey)
		// TODO delete from elastisearch
	} else {
		ok, err = dynamoRepo.DeleteByHash(table, hashKey)
		// TODO delete from elastisearch
	}
	writeBoolResponse(ok, err, w)
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	hashKey := vars["hash"]
	rangeKey := vars["range"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		createTransferObject(w, err)
	} else {
		item := createBussinesObject(data)

		var dynamoRepo DynamoCRUDRepository
		ok, err := dynamoRepo.Add(table, hashKey, rangeKey, item)
		// TODO correct elastisearch
		AddToElasticSearch(table, table, hashKey, rangeKey, item)
		writeBoolResponse(ok, err, w)
	}
}

func createBussinesObject(data map[string]interface{}) []Attribute {
	item := make([]Attribute, len(data))
	count := 0
	for key := range data {
		item[count] = Attribute{
			Description: AttributeDefinition{
				Name: key,
			},
			Value: data[key],
		}
		count++
	}

	return item
}

func createTransferObject(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func primaryKeySearch(rangeOperator, hashKey, rangeValue []string, table string, w http.ResponseWriter) {
	if rangeOperator == nil {
		if hashKey != nil {
			getItemsByHash(table, hashKey[0], w)
		} else {
			writeErrorResponse("Missing hash value", 404, w)
		}
	} else if hashKey != nil && rangeValue != nil {
		var repo DynamoCRUDRepository
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
		var repo DynamoCRUDRepository
		data, err := repo.GetByIndexRange(table, index[0], hashKey[0], rangeOperator[0], rangeValue)
		writeResponse(data, err, w)
	} else {
		writeErrorResponse("Missing primary key value(s)", 404, w)
	}
}

func getItemsByHash(table, hash string, w http.ResponseWriter) {
	schema, err := GetTableDescription(table, schema.Tables)
	if err != nil {
		writeErrorResponse("Table "+table+" not found.", 404, w)
	}

	var repo DynamoCRUDRepository
	if schema.HasRange() {
		data, err := repo.GetByOnlyHash(table, hash)
		writeResponse(data, err, w)
	} else {
		data, err := repo.GetByHash(table, hash)
		writeResponse(data, err, w)
	}
}

func getItemsByIndexHash(table, indexName, hash string, w http.ResponseWriter) {
	var repo DynamoCRUDRepository
	data, err := repo.GetByIndexHash(table, indexName, hash)
	writeResponse(data, err, w)
}

func writeResponse(data interface{}, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		createTransferObject(w, GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		createTransferObject(w, data)
	}
}

func writeBoolResponse(ok bool, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
}

func writeErrorResponse(message string, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	createTransferObject(w, newError(message, statusCode))
}

func getValue(queryParams []string) string {
	if queryParams != nil {
		return queryParams[0]
	}

	return ""
}
