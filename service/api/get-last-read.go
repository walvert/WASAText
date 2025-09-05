package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getLastRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	idParam := ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Warn("invalid chatId")
		http.Error(w, "Invalid chatId", http.StatusBadRequest)
		return
	}

	lastReadId, err := rt.db.GetLastRead(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Warn("message id not found")
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Warn("get last read failed")
			http.Error(w, "get last read failed", http.StatusInternalServerError)
			return
		}
	}

	response := map[string]interface{}{
		"lastReadId": lastReadId,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Warn("encode response failed")
		http.Error(w, "encode response failed", http.StatusInternalServerError)
		return
	}
}
