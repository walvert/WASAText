package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")
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
			ctx.Logger.WithError(err).Error("Internal Server Error")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	chats, err := rt.db.GetMyConversations(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if chats == nil {
		chats = make([]types.Chat, 0)
	}

	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode and return chats")
		http.Error(w, "Failed to return chats", http.StatusInternalServerError)
		return
	}
}
