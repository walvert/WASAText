package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) sendFirstMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var messageRequest types.FirstMessageRequest
	var chatID int

	err := json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	idParam := ps.ByName("userId")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if len(messageRequest.Receivers) == 1 {
		user1ID := userID
		user2ID := messageRequest.Receivers[0]

		if user2ID <= user1ID {
			user1ID, user2ID = user2ID, user1ID
		}

		chatID, err = rt.db.GetPrivateChatID(user1ID, user2ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				chatID, err = rt.db.CreateChat(messageRequest.ChatName, false)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Error creating chat")
					http.Error(w, "Error creating chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user1ID, chatID)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Error adding chat to user 1")
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user2ID, chatID)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Error adding chat to user 2")
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				privateChat := types.PrivateChat{
					User1ID: user1ID,
					User2ID: user2ID,
					ChatID:  chatID,
				}

				err = rt.db.AddPrivateChat(privateChat)
				if err != nil {
					rt.baseLogger.WithError(err).Error("Error adding private chat")
					http.Error(w, "Error adding private chat", http.StatusInternalServerError)
				}
			} else {
				rt.baseLogger.WithError(err).Error("Error getting private chat")
				http.Error(w, "Error getting chat", http.StatusInternalServerError)
				return
			}
		}
	} else if len(messageRequest.Receivers) >= 1 {
		if messageRequest.ChatName == "" {
			http.Error(w, "Must specify a chat name", http.StatusBadRequest)
			return
		}

		chatID, err := rt.db.CreateChat(messageRequest.ChatName, true)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error creating group chat")
			http.Error(w, "Error creating group chat", http.StatusInternalServerError)
			return
		}

		err = rt.db.AddChatToUser(userID, chatID)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error adding chat to user 1")
			http.Error(w, "Error adding chat", http.StatusInternalServerError)
			return
		}

		for _, receiverID := range messageRequest.Receivers {
			err = rt.db.AddChatToUser(receiverID, chatID)
			if err != nil {
				rt.baseLogger.WithError(err).Error("Error adding chat to array user")
				http.Error(w, "Error adding chat", http.StatusInternalServerError)
				return
			}
		}
	}

	err = rt.db.SendMessage(chatID, userID, messageRequest.Text)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to send message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var messageRequest types.MessageRequest

	err := json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	idParam := ps.ByName("userId")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	idParam = ps.ByName("chatId")
	chatID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	err = rt.db.SendMessage(chatID, userID, messageRequest.Text)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to send message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
