package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorBody struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type HTTPErrorResponse struct {
	Status string    `json:"status"`
	Error  ErrorBody `json:"error"`
}

type HTTPError struct {
	Code    int
	Message string
	ErrId   string
	Err     error
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (fn HTTPErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		fmt.Printf("[ERROR] %s\n", err)

		if httpError, ok := err.(*HTTPError); ok {
			w.WriteHeader(httpError.Code)
			json.NewEncoder(w).Encode(HTTPErrorResponse{
				Status: "error",
				Error: ErrorBody{
					Id:      httpError.ErrId,
					Message: httpError.Message,
				},
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HTTPErrorResponse{
				Status: "error",
				Error: ErrorBody{
					Id:      ERR_SERVER_ERROR,
					Message: "Something unexpected happened, please try again later.",
				},
			})
		}
		return
	}
}
