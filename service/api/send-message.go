package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}
	var messageRequest types.MessageRequest
	var messageId int

	err := json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error decoding request body")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
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

	idParam := ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid ChatId")
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	username, err := rt.db.GetUsernameById(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("User not found: GetUsernameById")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("Internal Server Error: GetUsernameById")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if messageRequest.Type == "text" {
		messageId, err = rt.db.SendMessage(chatId, userId, username, messageRequest.Type, messageRequest.Text, messageRequest.MediaURL, messageRequest.IsForward, messageRequest.ReplyTo)
		if err != nil {
			ctx.Logger.WithError(err).Error("Internal Server Error: send text message")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		messageId, err = rt.db.SendMessage(chatId, userId, username, messageRequest.Type, messageRequest.Text, messageRequest.MediaURL, messageRequest.IsForward, messageRequest.ReplyTo)
		if err != nil {
			ctx.Logger.WithError(err).Error("Internal Server Error: send image message")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Rename the file using message_id
		ext := filepath.Ext(messageRequest.MediaURL) // Get file extension
		newFilePath := fmt.Sprintf("uploads/messages/msg_%d%s", messageId, ext)

		err = os.Rename(messageRequest.MediaURL, newFilePath) // Rename file
		if err != nil {
			http.Error(w, "Failed to rename file", http.StatusInternalServerError)
			return
		}
	}

	response, err := rt.db.GetMessage(messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error: get message response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error: encode response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
