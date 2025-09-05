package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		ctx.Logger.WithError(err).Error("Error converting message ID to int")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chatDeleted, err := rt.db.DeleteMessage(id)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error deleting message")
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
