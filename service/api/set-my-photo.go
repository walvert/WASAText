package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/png")

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
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

	// Parse the form data
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error parsing multipart form")
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
	uploadDir := "uploads/user/images/"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	imagePath := fmt.Sprintf("%s%d%s", uploadDir, userId, ext)

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

	err = rt.db.SetMyPhoto(userId, imagePath)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to save image")
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	response := types.SetImageResponse{
		Success:  true,
		Message:  "Profile picture updated success.",
		ImageURL: outFile.Name(),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to return user", http.StatusInternalServerError)
		return
	}
}
