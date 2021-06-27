package errors

import (
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
)

// ErrorResponse is the response that represents an error.
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error is required by the error interface.
func (e ErrorResponse) Error() string {
	return e.Message
}

// StatusCode is required by routing.HTTPError interface.
func (e ErrorResponse) StatusCode() int {
	return e.Status
}

// InternalServerError creates a new error response representing an internal server error (HTTP 500)
func InternalServerError(msg string) ErrorResponse {
	if msg == "" {
		msg = "We encountered an error while processing your request."
	}
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: msg,
	}
}

// NotFound creates a new error response representing a resource-not-found error (HTTP 404)
func NotFound(msg string) ErrorResponse {
	if msg == "" {
		msg = "The requested resource was not found."
	}
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}

// Unauthorized creates a new error response representing an authentication/authorization failure (HTTP 401)
func Unauthorized(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authenticated to perform the requested action."
	}
	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: msg,
	}
}

// Forbidden creates a new error response representing an authorization failure (HTTP 403)
func Forbidden(msg string) ErrorResponse {
	if msg == "" {
		msg = "You are not authorized to perform the requested action."
	}
	return ErrorResponse{
		Status:  http.StatusForbidden,
		Message: msg,
	}
}

// BadRequest creates a new error response representing a bad request (HTTP 400)
func BadRequest(msg string) ErrorResponse {
	if msg == "" {
		msg = "Your request is in a bad format."
	}
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: msg,
	}
}

// InvalidRequest creates a new error response representing a data validation error (HTTP 400).
func InvalidRequestData(e error, t ut.Translator) ErrorResponse {
	errs := e.(validator.ValidationErrors)
	err := errs.Translate(t)
	var data = make(map[string]string)
	for k, v := range err {
		newKey := strings.Split(k, ".")
		key := newKey[1]
		key = strings.ToLower(key)
		data[key] = v
	}

	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "There is a problem with the data you submitted.",
		Data:    data,
	}
}
