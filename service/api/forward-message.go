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

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	ctx.Logger.Infof("getting id with token: %s", token)
	userId, err := rt.db.GetIdWithToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("Internal Server Error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var request types.ForwardRequest
	var results []types.ForwardResult
	var chatId int

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	idParam := ps.ByName("messageId")
	messageId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid message ID")
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	ctx.Logger.Infof("getting username with id: %d", userId)
	username, err := rt.db.GetUsernameById(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Error("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			rt.baseLogger.WithError(err).Error("Internal Server Error: GetUsernameById")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	ctx.Logger.Infof("getting message info with message id: %d", messageId)
	messageInfo, err := rt.db.GetForwardInfo(messageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("Message not found")
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		} else {
			rt.baseLogger.WithError(err).Error("Internal Server Error: GetForwardInfo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	for _, recipient := range request.Recipients {
		if recipient.Type == "user" {
			user1Id := userId
			user2Id := recipient.ID

			ctx.Logger.Infof("getting chat id with user 2 id: %d", user2Id)
			if user2Id <= user1Id {
				user1Id, user2Id = user2Id, user1Id
			}

			chatId, err = rt.db.GetPrivateChatID(user1Id, user2Id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					ctx.Logger.Infof("creating new chat")
					chatId, err = rt.db.CreateChat("", false)
					if err != nil {
						http.Error(w, "Error creating chat", http.StatusInternalServerError)
						return
					}

					ctx.Logger.Infof("adding chat id: %d to user id: %d", chatId, user1Id)
					err = rt.db.AddChatToUser(user1Id, chatId)
					if err != nil {
						http.Error(w, "Error adding chat", http.StatusInternalServerError)
						return
					}

					ctx.Logger.Infof("adding chat id: %d to user id: %d", chatId, user2Id)
					err = rt.db.AddChatToUser(user2Id, chatId)
					if err != nil {
						http.Error(w, "Error adding chat", http.StatusInternalServerError)
						return
					}

					ctx.Logger.Infof("adding private chat to users: %d, %d", user1Id, user2Id)
					err = rt.db.AddPrivateChat(user1Id, user2Id, chatId)
					if err != nil {
						http.Error(w, "Error adding private chat", http.StatusInternalServerError)
					}

					ctx.Logger.Infof("setting last read to user: %d for chat: %d", user1Id, chatId)
					err = rt.db.SetLastRead(user1Id, chatId, 0)
					if err != nil {
						ctx.Logger.WithError(err).Error("Internal Server Error: SetLastRead 1")
						http.Error(w, "Error setting last read", http.StatusInternalServerError)
						return
					}

					ctx.Logger.Infof("setting last read to user: %d for chat: %d", user2Id, chatId)
					err = rt.db.SetLastRead(user2Id, chatId, 0)
					if err != nil {
						ctx.Logger.WithError(err).Error("Internal Server Error: SetLastRead 2")
						http.Error(w, "Error setting last read", http.StatusInternalServerError)
						return
					}

				} else {
					http.Error(w, "Error getting chat", http.StatusInternalServerError)
					return
				}
			}
		} else {
			chatId = recipient.ID
		}

		ctx.Logger.Infof("sending message to chat: %d", chatId)
		newMessageId, err := rt.db.SendMessage(chatId, userId, username, messageInfo.Type, messageInfo.Text, messageInfo.MediaUrl, true, 0)
		if err != nil {
			ctx.Logger.WithError(err).Error("Internal Server Error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		results = append(results, types.ForwardResult{
			ChatID:    chatId,
			MessageID: newMessageId,
		})
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to return response", http.StatusInternalServerError)
		return
	}
}
