package api

import (
	"encoding/json"
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

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/png")

	idParam := ps.ByName("chatId")
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
	uploadDir := "uploads/chats/images/"
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

	err = rt.db.SetGroupPhoto(id, imagePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	response := types.SetImageResponse{
		Success:  true,
		Message:  "Successfully updated group image",
		ImageURL: outFile.Name(),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		rt.baseLogger.WithError(err).Error("Failed to encode response")
		http.Error(w, "Failed to return user", http.StatusInternalServerError)
		return
	}
}
