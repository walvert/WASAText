package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "No Authorization Header found", http.StatusUnauthorized)
		return
	}

	userId, err := rt.db.GetIdWithToken(token)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting user id")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idParam := ps.ByName("messageId")
	messageId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error converting message id")
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	err = rt.db.CommentMessage(messageId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error commenting message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := rt.db.GetComments(messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting comments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error encoding comments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "No Authorization Header found", http.StatusUnauthorized)
		return
	}

	userId, err := rt.db.GetIdWithToken(token)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting user id")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idParam := ps.ByName("messageId")
	messageId, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error converting message id")
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteComment(messageId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error deleting comment")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := rt.db.GetComments(messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting comments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error encoding comments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
