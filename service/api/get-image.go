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

	validFolders := map[string]bool{
		"user":     true,
		"chats":    true,
		"messages": true,
	}
	if !validFolders[folder] {
		ctx.Logger.Warnf("Invalid folder %s", folder)
		http.Error(w, "Invalid folder", http.StatusBadRequest)
		return
	}

	if strings.Contains(filename, "..") {
		ctx.Logger.Warnf("Invalid filename %s", filename)
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("./uploads/%s/%s", folder, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.Logger.Warnf("File %s does not exist", filePath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
