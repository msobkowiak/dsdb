package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Data struct {
	Name  string
	Value string
	Type  string
}

func GetAllItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	data, err := RepoGetAllItems(table)
	writeCollectionResponse(data, err, w)
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

	data, err := RepoGetItemByHashRange(table, hash, rangeKey)
	writeSingleItemResponse(data, err, w)
}

func GetByRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rangeKey := vars["range"]
	table := vars["table"]

	data, err := RepoGetItemsByHash(table, rangeKey)
	writeCollectionResponse(data, err, w)
}

func Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryParams := r.URL.Query()

	searchType := queryParams["search_type"]
	index := queryParams["index"]
	hashKey := queryParams["hash"]
	rangeOperator := queryParams["range_operator"]
	rangeValue := queryParams["range_value"]

	table := vars["table"]

	if searchType[0] == "index" {
		if index[0] == "primary" {
			primaryKeySearch(rangeOperator, hashKey, rangeValue, table, w)
		} else {
			secondaryIndexSearch(index, rangeOperator, hashKey, rangeValue, table, w)
		}
	}
}

func DeleteByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	table := vars["table"]

	ok, err := RepoDeleteItem(table, hash)
	writeDeleteResponse(ok, err, w)
}

func DeleteByHashRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashKey := vars["hash"]
	rangeKey := vars["range"]
	table := vars["table"]

	ok, err := RepoDeleteItemWithRange(table, hashKey, rangeKey)
	writeDeleteResponse(ok, err, w)
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	hashKey := vars["hash"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	var data []Data
	if err := json.Unmarshal(body, &data); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	item := convertDataToAttrubute(data)

	t, _ := RepoAddItem(table, hashKey, item)
	AddToElasticSearch(dbDescription.Name, table, hashKey, "", item)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

func AddItemHashRange(w http.ResponseWriter, r *http.Request) {
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
	var data []Data
	if err := json.Unmarshal(body, &data); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	item := convertDataToAttrubute(data)

	t, _ := RepoAddItemHashRange(table, hashKey, rangeKey, item)
	AddToElasticSearch(dbDescription.Name, table, hashKey, rangeKey, item)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}

/*func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}*/

func convertDataToAttrubute(data []Data) []Attribute {
	item := make([]Attribute, len(data))
	for i := range data {
		item[i] = Attribute{
			Description: AttributeDefinition{
				Name: data[i].Name,
				Type: data[i].Type,
			},
			Value: data[i].Value,
		}
	}

	return item
}

func primaryKeySearch(rangeOperator, hashKey, rangeValue []string, table string, w http.ResponseWriter) {
	if rangeOperator == nil {
		if hashKey != nil {
			getItemsByHash(table, hashKey[0], w)
		} else {
			writeErrorResponse("Missing hash value", 404, w)
		}
	} else if hashKey != nil && rangeValue != nil {
		data, err := RepoGetItemsByRangeOp(table, hashKey[0], rangeOperator[0], rangeValue)
		writeCollectionResponse(data, err, w)
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
		data, err := RepoGetItemsByIndexRangeOp(table, index[0], hashKey[0], rangeOperator[0], rangeValue)
		writeCollectionResponse(data, err, w)
	} else {
		writeErrorResponse("Missing primary key value(s)", 404, w)
	}
}

func getItemsByHash(table, hash string, w http.ResponseWriter) {
	schema, err := GetTableDescription(table, schema.Tables)
	if err != nil {
		writeErrorResponse("Table "+table+" not found.", 404, w)
	}

	if schema.HasRange() {
		data, err := RepoGetItemsByHash(table, hash)
		writeCollectionResponse(data, err, w)
	} else {
		data, err := RepoGetItemByHash(table, hash)
		writeSingleItemResponse(data, err, w)
	}
}

func getItemsByIndexHash(table, indexName, hash string, w http.ResponseWriter) {
	data, err := RepoGetItemByIndexHash(table, indexName, hash)
	writeCollectionResponse(data, err, w)
}

func writeSingleItemResponse(data map[string]string, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(data)
	}
}

func writeCollectionResponse(data []map[string]string, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(data)
	}
}

func writeDeleteResponse(ok bool, err error, w http.ResponseWriter) {
	if ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
}

func writeErrorResponse(message string, statusCode int, w http.ResponseWriter) {
	err := Error{
		statusCode,
		message,
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(err)
}
