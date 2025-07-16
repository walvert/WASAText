package api

import (
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strings"
)

func (rt *_router) getImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "No token found", http.StatusUnauthorized)
		return
	}

	folder := ps.ByName("folder")
	filename := ps.ByName("filename")

	// Debug logging
	fmt.Printf("DEBUG: folder='%s', filename='%s'\n", folder, filename)
	fmt.Printf("DEBUG: Full URL path: %s\n", r.URL.Path)

	validFolders := map[string]bool{
		"user":     true,
		"chats":    true,
		"messages": true,
	}
	if !validFolders[folder] {
		fmt.Printf("DEBUG: Invalid folder '%s'\n", folder)
		http.Error(w, "Invalid folder", http.StatusBadRequest)
		return
	}

	if strings.Contains(filename, "..") {
		fmt.Printf("DEBUG: Invalid filename '%s' (contains ..)\n", filename)
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("./uploads/%s/%s", folder, filename)
	fmt.Printf("DEBUG: Looking for file at: %s\n", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("DEBUG: File not found: %s\n", filePath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fmt.Printf("DEBUG: Serving file: %s\n", filePath)
	http.ServeFile(w, r, filePath)
}
