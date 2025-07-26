package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	idParam := ps.ByName("messageId")
	messageId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid id parameter")
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	comments, err := rt.db.GetComments(messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if comments == nil {
		comments = make([]string, 0)
	}

	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to encode and return chats")
		http.Error(w, "Failed to return chats", http.StatusInternalServerError)
		return
	}
}
