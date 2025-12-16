package handlers

import (
	"encoding/json"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) error {
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Welcome",
	})

	return nil
}
