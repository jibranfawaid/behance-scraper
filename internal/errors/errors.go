package errors

import "net/http"

type Errors struct {
	Key     string `json:"-"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

const (
	// 500 internal error
	GeneralError      = "general_error"
	ConnectionTimeout = "connection_timeout"
)

var errors = map[string]Errors{
	GeneralError: {
		Key:     GeneralError,
		Status:  http.StatusInternalServerError,
		Message: "Internal error",
	},
	ConnectionTimeout: {
		Key:     ConnectionTimeout,
		Status:  http.StatusInternalServerError,
		Message: "Connection timeout",
	},
}

func GetError(key string) Errors {
	value, ok := errors[key]
	if !ok {
		value = errors[GeneralError]
	}
	return value
}
