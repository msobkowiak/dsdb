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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(data)
	}
}

func GetByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	table := vars["table"]

	if GetSchema(table).HasRange() {
		getItemWithRange(table, hash, w)
	} else {
		getItem(table, hash, w)
	}
}

func GetByHashRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	rangeKey := vars["range"]
	table := vars["table"]

	data, err := RepoGetItemByHashRange(table, hash, rangeKey)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(data)
	}
}

func GetByRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rangeKey := vars["range"]
	table := vars["table"]

	data, err := RepoGetItemByRange(table, rangeKey)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(data)
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queryParams := r.URL.Query()

	fmt.Println(queryParams)
	searchType := queryParams["search_type"]
	index := queryParams["index"]
	hashKey := queryParams["hash"]
	rangeOperator := queryParams["range_operator"]
	rangeValue := queryParams["range_value"]

	table := vars["table"]

	if searchType[0] == "index" && (index[0] == "primary" || index[0] == "secendary") {
		if rangeOperator == nil {
			if hashKey != nil {
				if GetSchema(table).HasRange() {
					fmt.Println("1")
					getItemWithRange(table, hashKey[0], w)
				} else {
					fmt.Println("2")
					getItem(table, hashKey[0], w)
				}
			} else {
				err := Error{
					404,
					"Missing hash value",
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(err)
			}
		} else if rangeOperator != nil {
			if hashKey != nil && rangeValue != nil {
				data, err := RepoGetItemsByRangeOp(table, hashKey[0], rangeOperator[0], rangeValue)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
				} else {
					w.WriteHeader(http.StatusFound)
					json.NewEncoder(w).Encode(data)
				}
			}
		}
	} else {
		err := Error{
			404,
			"Invalid search parameters",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	}
}

func DeleteByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	table := vars["table"]

	ok, err := RepoDeleteItem(table, hash)
	if ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
}

func DeleteByHashRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashKey := vars["hash"]
	rangeKey := vars["range"]
	table := vars["table"]

	ok, err := RepoDeleteItemWithRange(table, hashKey, rangeKey)
	if ok {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}
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

func getItem(table, hash string, w http.ResponseWriter) {
	data, err := RepoGetItemByHash(table, hash)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(data)
	}
}

func getItemWithRange(table, hash string, w http.ResponseWriter) {
	data, err := RepoGetItemByRange(table, hash)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetErrorMsg(err, 404))
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(data)
	}
}
