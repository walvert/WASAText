package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

	messageType, err := rt.db.GetMessageType(messageId)
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
		if recipient.Type == "user" {
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
			}
		}

		if messageType == "text" {
			_, err = rt.db.SendMessage(chatId, userId, username, text, messageType, "", true, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			newMessageId, err := rt.db.SendMessage(chatId, userId, username, text, messageType, "", true, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			files, err := os.ReadDir("uploads/media")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var ext string
			for _, file := range files {
				name := file.Name()
				if strings.HasPrefix(name, "msg_"+strconv.Itoa(messageId)+".") {
					ext = filepath.Ext(name) // Returns extension with dot (e.g., ".jpg")
				}
			}

			originalPath := fmt.Sprintf("uploads/media/msg_%d%s", messageId, ext)
			originalFile, err := os.Open(originalPath)
			if err != nil {
				http.Error(w, "Error opening file", http.StatusInternalServerError)
				return
			}
			err = originalFile.Close()
			if err != nil {
				http.Error(w, "Error closing file", http.StatusInternalServerError)
				return
			}

			newPath := fmt.Sprintf("uploads/media/msg_%d%s", newMessageId, ext)
			newFile, err := os.Create(newPath)
			if err != nil {
				http.Error(w, "Error creating file", http.StatusInternalServerError)
				return
			}
			err = newFile.Close()
			if err != nil {
				http.Error(w, "Error closing file", http.StatusInternalServerError)
				return
			}

			_, err = io.Copy(newFile, originalFile)
			if err != nil {
				http.Error(w, "Error copying file", http.StatusInternalServerError)
				return
			}
		}
	}
}
