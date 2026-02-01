package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTP Error ID
const (
	ERR_UNAUTHORIZED       = "unauthorized"
	ERR_INVALID_QUERY      = "invalid_query"
	ERR_INVALID_BODY       = "invalid_body"
	ERR_ACCESS_DENIED      = "access_denied"
	ERR_RESOURCE_NOT_FOUND = "resource_not_found"
	ERR_SERVER_ERROR       = "server_error"
)

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
