package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	//"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
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

	data, err := RepoGetItemByRange(table, rangeKey)
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
		fmt.Println("index")
		if index[0] == "primary" {
			fmt.Println("primay")
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
			/*} else if index[0] == "secondary" {
				fmt.Println("secondary")
				if rangeOperator == nil {
					if hashKey != nil {
						fmt.Println("no range")
						getItemsByIndexHash(table, hashKey[0], w)
					} else {
						writeErrorResponse("Missing hash value", 404, w)
					}
				} else if hashKey != nil && rangeValue != nil {
					data, err := RepoGetItemsByIndexRangeOp(table, hashKey[0], rangeOperator[0], rangeValue)
					writeCollectionResponse(data, err, w)
				} else {
					writeErrorResponse("Missing primary key value(s)", 404, w)
				}
			} else {
				writeErrorResponse("Invalid search parameters", 404, w)*/
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

	var item []Attribute
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &item); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t, _ := RepoAddItem(table, hashKey, item)
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

	var item []Attribute
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &item); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t, _ := RepoAddItemHashRange(table, hashKey, rangeKey, item)
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

func getItemsByHash(table, hash string, w http.ResponseWriter) {
	if GetSchema(table).HasRange() {
		data, err := RepoGetItemByRange(table, hash)
		writeCollectionResponse(data, err, w)
	} else {
		data, err := RepoGetItemByHash(table, hash)
		writeSingleItemResponse(data, err, w)
	}
}

/*func getItemsByIndexHash(table, hash string, w http.ResponseWriter) {
	if GetSchema(table).HasIndexWRange() {
		data, err := RepoGetItemByIndexRange(table, hash)
		writeCollectionResponse(data, err, w)
	} else {
		data, err := RepoGetItemByIndexHash(table, hash)
		writeCollectionResponse(data, err, w)
		//writeSingleItemResponse(data, err, w)
	}
}*/

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
