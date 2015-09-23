package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	var repo DynamoBaseRepository
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

	var repo DynamoBaseRepository
	data, err := repo.GetByHashRange(table, hash, rangeKey)
	writeResponse(data, err, w)
}

func GetByRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rangeKey := vars["range"]
	table := vars["table"]

	var repo DynamoBaseRepository
	data, err := repo.GetByOnlyHash(table, rangeKey)
	writeResponse(data, err, w)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashKey := vars["hash"]
	rangeKey := vars["range"]
	table := vars["table"]

	var (
		ok          bool
		err         error
		dynamoRepo  DynamoBaseRepository
		elasticRepo ElasticBaseepository
	)
	if rangeKey != "" {
		ok, err = dynamoRepo.DeleteByHashRange(table, hashKey, rangeKey)
		elasticRepo.DeleteByHashRange(table, hashKey, rangeKey)
	} else {
		ok, err = dynamoRepo.DeleteByHash(table, hashKey)
		elasticRepo.DeleteByHash(table, hashKey)
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
		item := createBussinesObject(data, table)

		var dynamoRepo DynamoBaseRepository
		var elasticRepo ElasticBaseepository
		ok, err := dynamoRepo.Add(table, hashKey, rangeKey, item)
		elasticRepo.Add(table, hashKey, rangeKey, item)
		writeBoolResponse(ok, err, w)
	}
}

func createBussinesObject(data map[string]interface{}, tableName string) []Attribute {
	table, _ := schema.GetTableDescription(tableName)
	item := make([]Attribute, len(data))
	count := 0
	for key := range data {
		fmt.Println(key)
		item[count] = Attribute{
			Description: AttributeDefinition{
				Name: key,
				Type: table.GetTypeOfAttribute(key),
			},
			Value: data[key],
		}
		count++
	}

	fmt.Println(item)

	return item
}

func createTransferObject(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func getItemsByHash(table, hash string, w http.ResponseWriter) {
	schema, err := schema.GetTableDescription(table)
	if err != nil {
		writeErrorResponse("Table "+table+" not found.", 404, w)
	}

	var repo DynamoBaseRepository
	if schema.HasRange() {
		data, err := repo.GetByOnlyHash(table, hash)
		writeResponse(data, err, w)
	} else {
		data, err := repo.GetByHash(table, hash)
		writeResponse(data, err, w)
	}
}

func getItemsByIndexHash(table, indexName, hash string, w http.ResponseWriter) {
	var repo DynamoBaseRepository
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
