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
	"strings"
	"time"
)

func (rt *_router) createChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Authorization check
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	// Get user info
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

	// Check content type to determine how to parse the request
	contentType := r.Header.Get("Content-Type")

	var messageRequest types.FirstMessageRequest
	var uploadedFileName string

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle file upload
		messageRequest, uploadedFileName, err = rt.handleCreateChatFileUpload(r, ctx)
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

	var chatId int

	if len(messageRequest.Receivers) == 1 {
		user1Id := userId
		user2Id, err := rt.db.GetUserByUsername(messageRequest.Receivers[0])
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.Logger.WithError(err).Error("User not found")
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
		}

		if user2Id < user1Id {
			user1Id, user2Id = user2Id, user1Id
		}

		chatId, err = rt.db.GetPrivateChatID(user1Id, user2Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				chatId, err = rt.db.CreateChat(messageRequest.ChatName, false)
				if err != nil {
					ctx.Logger.WithError(err).Error("Internal Server Error: CreateChat")
					http.Error(w, "Error creating chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user1Id, chatId)
				if err != nil {
					ctx.Logger.WithError(err).Error("Internal Server Error: AddChatToUser 1")
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddChatToUser(user2Id, chatId)
				if err != nil {
					ctx.Logger.WithError(err).Error("Internal Server Error: AddChatToUser 2")
					http.Error(w, "Error adding chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.AddPrivateChat(user1Id, user2Id, chatId)
				if err != nil {
					ctx.Logger.WithError(err).Error("Internal Server Error: AddPrivateChat")
					http.Error(w, "Error adding private chat", http.StatusInternalServerError)
					return
				}

				err = rt.db.SetLastRead(user1Id, chatId, 0)
				if err != nil {
					ctx.Logger.WithError(err).Error("Internal Server Error: SetLastRead 1")
					http.Error(w, "Error setting last read", http.StatusInternalServerError)
					return
				}

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
	} else if len(messageRequest.Receivers) >= 1 {
		if messageRequest.ChatName == "" {
			ctx.Logger.WithError(err).Error("ChatName required")
			http.Error(w, "Must specify a chat name", http.StatusBadRequest)
			return
		}

		chatId, err = rt.db.CreateChat(messageRequest.ChatName, true)
		if err != nil {
			ctx.Logger.WithError(err).Error("Internal Server Error: CreateGroupChat")
			http.Error(w, "Error creating group chat", http.StatusInternalServerError)
			return
		}

		err = rt.db.AddChatToUser(userId, chatId)
		if err != nil {
			ctx.Logger.WithError(err).Error("Internal Server Error: AddGroupChatToUser 1")
			http.Error(w, "Error adding chat", http.StatusInternalServerError)
			return
		}

		for _, receiverUsername := range messageRequest.Receivers {
			receiverID, err := rt.db.GetUserByUsername(receiverUsername)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					ctx.Logger.WithError(err).Errorf("User not found: %s", receiverUsername)
					http.Error(w, "User not found", http.StatusNotFound)
					return
				} else {
					ctx.Logger.WithError(err).Error("Internal Server Error: GetUsernameById")
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			err = rt.db.AddChatToUser(receiverID, chatId)
			if err != nil {
				ctx.Logger.WithError(err).Error("Internal Server Error: AddGroupChatToUsers")
				http.Error(w, "Error adding chat", http.StatusInternalServerError)
				return
			}

			err = rt.db.SetLastRead(receiverID, chatId, 0)
			if err != nil {
				ctx.Logger.WithError(err).Error("Internal Server Error: SetLastRead receivers")
				http.Error(w, "Error setting last read", http.StatusInternalServerError)
				return
			}
		}
	}

	// Send initial message if provided
	if messageRequest.Type != "" {
		_, err = rt.db.SendMessage(
			chatId,
			userId,
			username,
			messageRequest.Type,
			messageRequest.Text,
			uploadedFileName, // Use the uploaded filename
			messageRequest.IsForward,
			0,
		)
		if err != nil {
			ctx.Logger.WithError(err).Error("Internal Server Error: SendMessage")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	chatInfo, err := rt.db.GetChatInfo(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("Chat not found")
			http.Error(w, "Chat not found", http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("Internal Server Error: GetChatInfo")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(chatInfo)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error encoding chatInfo")
		http.Error(w, "Error encoding chatInfo", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) handleCreateChatFileUpload(r *http.Request, ctx reqcontext.RequestContext) (types.FirstMessageRequest, string, error) {
	var messageRequest types.FirstMessageRequest

	// Parse multipart form (32MB max memory)
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		return messageRequest, "", fmt.Errorf("failed to parse multipart form: %v", err)
	}

	// Get form values
	messageRequest.Type = r.FormValue("type")
	messageRequest.Text = r.FormValue("text")
	messageRequest.ChatName = r.FormValue("chatName")

	// Parse receivers (should be JSON array string)
	receiversStr := r.FormValue("receivers")
	if receiversStr != "" {
		err := json.Unmarshal([]byte(receiversStr), &messageRequest.Receivers)
		if err != nil {
			return messageRequest, "", fmt.Errorf("failed to parse receivers: %v", err)
		}
	}

	// Parse optional fields
	if isForward := r.FormValue("isForward"); isForward == "true" {
		messageRequest.IsForward = true
	}

	// Validate message type
	if messageRequest.Type != "image" && messageRequest.Type != "gif" {
		return messageRequest, "", fmt.Errorf("invalid message type for file upload: %s", messageRequest.Type)
	}

	// Get the uploaded file
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return messageRequest, "", fmt.Errorf("failed to get uploaded file: %v", err)
	}
	defer file.Close()

	// Validate file
	if err := rt.validateCreateChatUploadedFile(file, fileHeader, messageRequest.Type); err != nil {
		return messageRequest, "", err
	}

	// Generate unique filename
	filename := rt.generateCreateChatMessageFilename(fileHeader.Filename)

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/messages"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return messageRequest, "", fmt.Errorf("failed to create uploads directory: %v", err)
	}

	// Save the file
	filepath := filepath.Join(uploadsDir, filename)
	if err := rt.saveCreateChatUploadedFile(file, filepath); err != nil {
		return messageRequest, "", fmt.Errorf("failed to save file: %v", err)
	}

	ctx.Logger.Infof("File uploaded successfully for new chat: %s", filename)

	return messageRequest, filename, nil
}

func (rt *_router) validateCreateChatUploadedFile(file multipart.File, fileHeader *multipart.FileHeader, messageType string) error {
	// Check file size (10MB limit)
	const maxFileSize = 10 << 20 // 10 MB
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", fileHeader.Size, maxFileSize)
	}

	// Reset file pointer for content validation
	file.Seek(0, 0)

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	// Reset file pointer again
	file.Seek(0, 0)

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

	return fmt.Errorf("invalid file type: %s for message type: %s", contentType, messageType)
}

func (rt *_router) generateCreateChatMessageFilename(originalFilename string) string {
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

func (rt *_router) saveCreateChatUploadedFile(file multipart.File, filepath string) error {
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
