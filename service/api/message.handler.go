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
	var chatId int

	err := json.NewDecoder(r.Body).Decode(&messageRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	idParam := ps.ByName("userId")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if len(messageRequest.Receivers) == 1 {
		user1Id := userId
		user2Id := messageRequest.Receivers[0]

		if user2Id <= user1Id {
			user1Id, user2Id = user2Id, user1Id
		}

		chatId, err = rt.db.GetPrivateChatID(user1Id, user2Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				chatId, err = rt.db.CreateChat(messageRequest.ChatName, false)
				if err != nil {
					http.Error(w, "Error creating chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user1Id, chatId)
				if err != nil {
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user2Id, chatId)
				if err != nil {
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddPrivateChat(user1Id, user2Id, chatId)
				if err != nil {
					http.Error(w, "Error adding private chat", http.StatusInternalServerError)
				}
			} else {
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
			http.Error(w, "Error creating group chat", http.StatusInternalServerError)
			return
		}

		err = rt.db.AddChatToUser(userId, chatID)
		if err != nil {
			http.Error(w, "Error adding chat", http.StatusInternalServerError)
			return
		}

		for _, receiverID := range messageRequest.Receivers {
			err = rt.db.AddChatToUser(receiverID, chatID)
			if err != nil {
				http.Error(w, "Error adding chat", http.StatusInternalServerError)
				return
			}
		}
	}

	err = rt.db.SendMessage(chatId, userId, messageRequest.Text, messageRequest.IsForward)
	if err != nil {
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

	err = rt.db.SendMessage(chatID, userID, messageRequest.Text, messageRequest.IsForward)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to send message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(ps.ByName("messageId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteMessage(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("message deleted"))
	if err != nil {
		return
	}
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	idParam := ps.ByName("userId")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	idParam = ps.ByName("messageId")
	messageId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	err = rt.db.CommentMessage(messageId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("comment added"))
	if err != nil {
		return
	}
}

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	idParam := ps.ByName("userId")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	idParam = ps.ByName("messageId")
	messageId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteComment(messageId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("comment deleted"))
	if err != nil {
		return
	}
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var request types.ForwardRequest
	var chatId int

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	idParam := ps.ByName("messageId")
	messageId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	idParam = ps.ByName("userId")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	text, err := rt.db.GetMessageText(messageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Message not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, recipient := range request.Recipients {
		if recipient.Type == "chat" {
			err = rt.db.SendMessage(recipient.ID, userId, text, true)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if recipient.Type == "user" {
			user1Id := userId
			user2Id := recipient.ID

			if user2Id <= user1Id {
				user1Id, user2Id = user2Id, user1Id
			}

			chatId, err = rt.db.GetPrivateChatID(user1Id, user2Id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					chatId, err = rt.db.CreateChat("", false)
					if err != nil {
						http.Error(w, "Error creating chat", http.StatusInternalServerError)
						return
					}

					err = rt.db.AddChatToUser(user1Id, chatId)
					if err != nil {
						http.Error(w, "Error adding chat", http.StatusInternalServerError)
						return
					}

					err = rt.db.AddChatToUser(user2Id, chatId)
					if err != nil {
						http.Error(w, "Error adding chat", http.StatusInternalServerError)
						return
					}

					err = rt.db.AddPrivateChat(user1Id, user2Id, chatId)
					if err != nil {
						http.Error(w, "Error adding private chat", http.StatusInternalServerError)
					}
				} else {
					http.Error(w, "Error getting chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.SendMessage(chatId, userId, text, true)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
}
