package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

)

// GetUsersById handles GET /users/:id
func GetUsersById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	_ = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"handler": "GetUsersById",
		"method":  "GET",
		"path":    "/users/:id",
		"status":  "TODO: implement",
	})
}
