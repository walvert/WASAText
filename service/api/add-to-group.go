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

	lastRead, err := rt.db.GetLastRead(chatId)
	if err != nil {
		rt.baseLogger.WithError(err).Warn("get last read failed")
		http.Error(w, "get last read failed", http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Infof("last read: %d", lastRead)

	err = rt.db.AddToGroup(chatId, userId)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Conflict",
				"message": "User is already a member of this group.",
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
		} else {
			rt.baseLogger.WithError(err).Warn("error adding to group")
			http.Error(w, "Error adding to group", http.StatusInternalServerError)
			return
		}
	}

	err = rt.db.SetLastRead(userId, chatId, lastRead)
	if err != nil {
		rt.baseLogger.WithError(err).Warn("set last read failed")
		http.Error(w, "set last read failed", http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Infof("userId: %d chatId: %d lastRead: %d", userId, chatId, lastRead)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(types.User{
		ID:       userId,
		Username: username.Username,
	})
	if err != nil {
		rt.baseLogger.WithError(err).Error("error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
