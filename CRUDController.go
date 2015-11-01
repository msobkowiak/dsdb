package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const dataSizeLimit = 1048576

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
		elasticRepo ElasticBaseRepository
	)
	if rangeKey != "" {
		ok, err = dynamoRepo.DeleteByHashRange(table, hashKey, rangeKey)
		if err == nil {
			elasticRepo.DeleteByHashRange(table, hashKey, rangeKey)
		}
	} else {
		ok, err = dynamoRepo.DeleteByHash(table, hashKey)
		if err == nil {
			elasticRepo.DeleteByHash(table, hashKey)
		}
	}
	writeBoolResponse(ok, err, w)
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	hashKey := vars["hash"]
	rangeKey := vars["range"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, dataSizeLimit))
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
		t, _ := schema.GetTableDescription(table)
		requested := t.GetAllRequiredAttributes()
		var errs []ValidationError
		for _, r := range requested {
			if !hasAttribute(item, r) && !isPartOfPrimaryKey(r, table) {
				var v ValidationError
				error := v.New("Missing required attribute: " + r)
				errs = append(errs, error)
			}
		}

		if len(errs) > 0 {
			writeErrorsResponse(errs, 400, w)
		} else {
			var dynamoRepo DynamoBaseRepository
			var elasticRepo ElasticBaseRepository
			ok, err := dynamoRepo.Add(table, hashKey, rangeKey, item)
			if err == nil {
				elasticRepo.Add(table, hashKey, rangeKey, item)
			}
			writeBoolResponse(ok, err, w)
		}
	}
}

func createBussinesObject(data map[string]interface{}, tableName string) []Attribute {
	table, _ := schema.GetTableDescription(tableName)
	item := make([]Attribute, len(data))
	count := 0
	for key := range data {
		item[count] = Attribute{
			Description: AttributeDefinition{
				Name: key,
				Type: table.GetTypeOfAttribute(key),
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
	w.WriteHeader(statusCode)
	createTransferObject(w, newError(message, statusCode))
}

func writeErrorsResponse(errors []ValidationError, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	var e ErrorCollection
	createTransferObject(w, e.New(errors, statusCode))
}

func getValue(queryParams []string) string {
	if queryParams != nil {
		return queryParams[0]
	}

	return ""
}
