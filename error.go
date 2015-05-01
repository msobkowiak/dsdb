package main

type Error struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func GetErrorMsg(err error, statusCode int) Error {
	return Error{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}
