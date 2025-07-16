package api

import (
	"encoding/json"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "No Authorization Header found", http.StatusUnauthorized)
		return
	}

	userId, err := rt.db.GetIdWithToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	users, err := rt.db.GetUsers(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting users")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error encoding users")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
