package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/julienschmidt/httprouter"
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

	if users == nil {
		users = make([]types.User, 0)
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error encoding users")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
