package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var username types.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if username is missing
	if username.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get or create user
	userId, err := rt.db.GetUserByUsername(username.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userId, err = rt.db.CreateUser(username.Username)
			if err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}
	}

	// Get bearer token
	token := ctx.ReqUUID.String()

	userToken := types.BearerToken{
		Token:  ctx.ReqUUID.String(),
		UserID: userId,
	}

	err = rt.db.UpsertToken(userToken)
	if err != nil {
		http.Error(w, "Failed to set token", http.StatusInternalServerError)
	}

	// Create response
	response := types.LoginResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		rt.baseLogger.WithError(err).Error("Failed to encode response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Infof("response: %v", response)
}
