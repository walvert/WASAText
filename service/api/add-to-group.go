package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var userId int
	var username types.UsernameRequest

	idParam := ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
	}

	err = json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	userId, err = rt.db.GetUserByUsername(username.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Warn("user not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			rt.baseLogger.WithError(err).Warn("error getting user")
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}
	}

	err = rt.db.AddToGroup(chatId, userId)
	if err != nil {
		if errors.Is(err, database.AlreadyExists) {
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "User is already a member.",
			})
			if err != nil {
				rt.baseLogger.WithError(err).Warn("error encoding response")
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
				return
			}
			return
		} else if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Warn("user not found")
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		}
	}

	response := types.User{
		ID:       userId,
		Username: username.Username,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		rt.baseLogger.WithError(err).Error("error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
