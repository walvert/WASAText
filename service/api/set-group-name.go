package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	var request types.ChatNameRequest

	idParam := ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error converting string to int")
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error decoding request body")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = rt.db.SetGroupName(chatId, request.ChatName)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error setting group name")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chat, err := rt.db.GetChatInfo(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting chat info")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(chat)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
