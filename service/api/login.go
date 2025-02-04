package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var username types.LoginToken
	var user types.User
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	if rt.db.UsernameExists(username.Username) {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	user, err = rt.db.CreateUser(username.Username)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to return user", http.StatusInternalServerError)
		return
	}

}
