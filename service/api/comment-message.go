package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

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
