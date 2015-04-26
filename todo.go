package dsbs

import "time"

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}
type Todos []Todo

type Table struct {
	Name          string `json:"name"`
	Id            string `json:"id"`
	IdType        string `json:"idType"`
	ReadThrouput  int64  `json:"readThrouput"`
	WriteThrouput int64  `json:"writeThrouput"`
}
type Tables []Table
