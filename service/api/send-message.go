package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

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

	username, err := rt.db.GetUsernameById(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("User not found: GetUsernameById")
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("Internal Server Error: GetUsernameById")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse chat ID
	idParam := ps.ByName("chatId")
	chatId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid ChatId")
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	// Check content type to determine how to parse the request
	contentType := r.Header.Get("Content-Type")

	var messageRequest types.MessageRequest
	var uploadedFileName string

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle file upload
		messageRequest, uploadedFileName, err = rt.handleFileUpload(r, ctx)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to handle file upload")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Handle JSON request (text messages)
		err = json.NewDecoder(r.Body).Decode(&messageRequest)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error decoding request body")
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	}

	// Send message to database
	messageId, err := rt.db.SendMessage(
		chatId,
		userId,
		username,
		messageRequest.Type,
		messageRequest.Text,
		uploadedFileName, // Use the uploaded filename, empty for text messages
		messageRequest.IsForward,
		messageRequest.ReplyTo,
	)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error: send message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the complete message response
	response, err := rt.db.GetMessage(messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error: get message response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error("Internal Server Error: encode response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) handleFileUpload(r *http.Request, ctx reqcontext.RequestContext) (types.MessageRequest, string, error) {
	var messageRequest types.MessageRequest

	// Parse multipart form (32MB max memory)
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		return messageRequest, "", fmt.Errorf("failed to parse multipart form: %w", err)
	}

	// Get form values
	messageRequest.Type = r.FormValue("type")
	messageRequest.Text = r.FormValue("text")

	// Parse optional fields
	if isForward := r.FormValue("isForward"); isForward == "true" {
		messageRequest.IsForward = true
	}

	if replyToStr := r.FormValue("replyTo"); replyToStr != "" {
		if replyTo, err := strconv.Atoi(replyToStr); err == nil {
			messageRequest.ReplyTo = replyTo
		}
	}

	// Validate message type
	if messageRequest.Type != "image" && messageRequest.Type != "gif" {
		return messageRequest, "", fmt.Errorf("error: %w. invalid message type for file upload: %s", err, messageRequest.Type)
	}

	// Get the uploaded file
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return messageRequest, "", fmt.Errorf("failed to get uploaded file: %w", err)
	}
	defer file.Close()

	// Validate file
	if err := rt.validateUploadedFile(file, fileHeader, messageRequest.Type); err != nil {
		return messageRequest, "", err
	}

	// Generate unique filename
	filename := rt.generateMessageFilename(fileHeader.Filename)

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/messages/images"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return messageRequest, "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// Save the file
	myFilepath := filepath.Join(uploadsDir, filename)
	if err := rt.saveUploadedFile(file, myFilepath); err != nil {
		return messageRequest, "", fmt.Errorf("failed to save file: %w", err)
	}

	ctx.Logger.Infof("File uploaded successfully: %s", filename)

	return messageRequest, filename, nil
}

func (rt *_router) validateUploadedFile(file multipart.File, fileHeader *multipart.FileHeader, messageType string) error {
	// Check file size (10MB limit)
	const maxFileSize = 10 << 20 // 10 MB
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", fileHeader.Size, maxFileSize)
	}

	// Reset file pointer for content validation
	_, err := file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Reset file pointer again
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(buffer)

	// Validate content type based on message type
	validTypes := map[string][]string{
		"image": {"image/jpeg", "image/jpg", "image/png", "image/webp"},
		"gif":   {"image/gif"},
	}

	if allowedTypes, exists := validTypes[messageType]; exists {
		for _, allowed := range allowedTypes {
			if contentType == allowed {
				return nil // Valid type
			}
		}
	}

	return fmt.Errorf("error: %w. invalid file type: %s for message type: %s", err, contentType, messageType)
}

func (rt *_router) generateMessageFilename(originalFilename string) string {
	// Get file extension
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = ".jpg" // Default extension
	}

	// Generate unique filename with timestamp and random component
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("msg_%d_%d%s", timestamp, timestamp%10000, ext)

	return filename
}

func (rt *_router) saveUploadedFile(file multipart.File, filepath string) error {
	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy uploaded file to destination
	_, err = io.Copy(dst, file)
	return err
}
