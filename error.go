package main

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func newError(message string, statusCode int) Error {
	return Error{
		statusCode,
		message,
	}
}

func GetErrorMsg(err error, statusCode int) Error {
	return Error{
		StatusCode: statusCode,
		Message:    err.Error(),
	}
}
