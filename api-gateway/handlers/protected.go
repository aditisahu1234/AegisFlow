package handlers		//writes response to client

import (
	"encoding/json"
	"net/http"
)

func ProtectedHandler(
	w http.ResponseWriter,		//create a handler
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{			//response
		"message": "protected resource",
	}

	json.NewEncoder(w).Encode(response)		//sends json to client
}