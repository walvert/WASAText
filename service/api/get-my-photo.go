package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	if token == "" {
		ctx.Logger.WithError(errors.New("authorization required"))
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	imageUrl, err := rt.db.GetMyPhoto(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Warn("Unauthorized")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if errors.Is(err, database.ErrEmptyImageURL) {
			ctx.Logger.WithError(err).Info("image url empty")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			ctx.Logger.WithError(err).Error("error getting image url")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	err = json.NewEncoder(w).Encode(map[string]string{"imageUrl": imageUrl})
	if err != nil {
		ctx.Logger.WithError(err).Error("error encoding image url")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
