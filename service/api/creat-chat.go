package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"path/filepath"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var messageRequest types.FirstMessageRequest
	var chatId int

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
	}

	err := json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userId, err := rt.db.GetIdWithToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Error("User not found")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else {
			rt.baseLogger.WithError(err).Error("Internal Server Error")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

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

	if len(messageRequest.Receivers) == 1 {
		user1Id := userId
		user2Id := messageRequest.Receivers[0]

		if user2Id < user1Id {
			user1Id, user2Id = user2Id, user1Id
		}

		chatId, err = rt.db.GetPrivateChatID(user1Id, user2Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				chatId, err = rt.db.CreateChat(messageRequest.ChatName, false)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Internal Server Error: CreateChat")
					http.Error(w, "Error creating chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user1Id, chatId)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Internal Server Error: AddChatToUser 1")
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user2Id, chatId)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Internal Server Error: AddChatToUser 2")
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddPrivateChat(user1Id, user2Id, chatId)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Internal Server Error: AddPrivateChat")
					http.Error(w, "Error adding private chat", http.StatusInternalServerError)
				}

				err = rt.db.SetLastRead(user1Id, chatId)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Internal Server Error: SetLastRead 1")
					http.Error(w, "Error setting last read", http.StatusInternalServerError)
					return
				}

				err = rt.db.SetLastRead(user2Id, chatId)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Internal Server Error: SetLastRead 2")
					http.Error(w, "Error setting last read", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Error getting chat", http.StatusInternalServerError)
				return
			}
		}
	} else if len(messageRequest.Receivers) >= 1 {
		if messageRequest.ChatName == "" {
			rt.baseLogger.WithError(err).Error("ChatName required")
			http.Error(w, "Must specify a chat name", http.StatusBadRequest)
			return
		}

		chatId, err = rt.db.CreateChat(messageRequest.ChatName, true)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Internal Server Error: CreateGroupChat")
			http.Error(w, "Error creating group chat", http.StatusInternalServerError)
			return
		}

		err = rt.db.AddChatToUser(userId, chatId)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Internal Server Error: AddGroupChatToUser 1")
			http.Error(w, "Error adding chat", http.StatusInternalServerError)
			return
		}

		for _, receiverID := range messageRequest.Receivers {
			err = rt.db.AddChatToUser(receiverID, chatId)
			if err != nil {
				rt.baseLogger.WithError(err).Error("Internal Server Error: AddGroupChatToUsers")
				http.Error(w, "Error adding chat", http.StatusInternalServerError)
				return
			}

			err = rt.db.SetLastRead(receiverID, chatId)
			if err != nil {
				rt.baseLogger.WithError(err).Error("Internal Server Error: SetLastRead receivers")
				http.Error(w, "Error setting last read", http.StatusInternalServerError)
				return
			}
		}
	}

	rt.baseLogger.Info("Sending first message from user ", userId)
	rt.baseLogger.Info("Sending first message to chat ", chatId)

	if messageRequest.Type == "text" {
		_, err = rt.db.SendMessage(chatId, userId, username, messageRequest.Type, messageRequest.Text, messageRequest.MediaURL, messageRequest.IsForward, 0)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Internal Server Error: SendTextMessage")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		messageId, err := rt.db.SendMessage(chatId, userId, username, messageRequest.Type, messageRequest.Text, messageRequest.MediaURL, messageRequest.IsForward, 0)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error sending media message")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 2: Rename the file using message_id
		ext := filepath.Ext(messageRequest.MediaURL) // Get file extension
		newFilePath := fmt.Sprintf("uploads/media/msg_%d%s", messageId, ext)

		err = os.Rename(messageRequest.MediaURL, newFilePath) // Rename file
		if err != nil {
			rt.baseLogger.WithError(err).Error("Internal Server Error: Rename")
			http.Error(w, "Failed to rename file", http.StatusInternalServerError)
			return
		}
	}

	chatInfo, err := rt.db.GetChatInfo(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Error("Chat not found")
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		} else {
			rt.baseLogger.WithError(err).Error("Internal Server Error: GetChatInfo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	err = json.NewEncoder(w).Encode(chatInfo)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding chatInfo")
		http.Error(w, "Error encoding chatInfo", http.StatusInternalServerError)
		return
	}
}
