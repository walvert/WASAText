package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	idParam := ps.ByName("userId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chats, err := rt.db.GetUserChats(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			rt.baseLogger.WithError(err).Error("Chats not found")
			http.Error(w, "Chats not found", http.StatusNotFound)
			return
		} else {
			rt.baseLogger.WithError(err).Error("Internal Server Error")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to encode and return chats")
		http.Error(w, "Failed to return chats", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var userId int

	idParam := ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
	}

	err = json.NewDecoder(r.Body).Decode(&userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = rt.db.AddToGroup(chatId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Chat not found", http.StatusNotFound)
		} else if errors.Is(err, errors.New("unauthorized")) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{"chatId":` + strconv.Itoa(userId) + `}`))
	if err != nil {
		return
	}
}

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var username types.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	idParam := ps.ByName("userId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user := types.User{
		ID:       id,
		Username: username.Username,
	}

	err = rt.db.SetMyUsername(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to return user", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/png")

	idParam := ps.ByName("userId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// Parse the form data
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing the image", http.StatusInternalServerError)
		return
	} // Max file size: 10MB

	// Get the uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Ensure upload directory exists
	uploadDir := "uploads/user/"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	imagePath := fmt.Sprintf("%s%d%s", uploadDir, id, ext)

	// Create the new file
	outFile, err := os.Create(imagePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {

		}
	}(outFile)

	// Copy file content
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (rt *_router) LeaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	idParam := ps.ByName("userId")
	userId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	idParam = ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid chat id", http.StatusBadRequest)
		return
	}

	err = rt.db.LeaveGroup(chatId, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
