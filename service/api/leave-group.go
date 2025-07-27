package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	userId, err := rt.db.GetIdWithToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("User not found")
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
		ctx.Logger.WithError(err).Error("Invalid chat id")
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	chatDeleted, err := rt.db.LeaveGroup(chatId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error: Leave group")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(types.ChatDeleted{
		ChatDeleted: chatDeleted,
	})
	if err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
