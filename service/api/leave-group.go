package api

import (
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	userId, err := rt.db.GetIdWithToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Error("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			rt.baseLogger.WithError(err).Error("Internal Server Error: GetIdWithToken")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	idParam := ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Invalid chat id")
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	err = rt.db.LeaveGroup(chatId, userId)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Internal Server Error: Leave group")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
