package main

type ErrorCollection struct {
	StatusCode int               `json:"status_code"`
	Message    []ValidationError `json:"message"`
}

type ValidationError struct {
	ErrorMessage string `json:"validation_error"`
}

func (v ValidationError) New(message string) ValidationError {
	return ValidationError{
		ErrorMessage: message,
	}
}

func (e ErrorCollection) New(message []ValidationError, statusCode int) ErrorCollection {
	return ErrorCollection{
		StatusCode: statusCode,
		Message:    message,
	}
}
