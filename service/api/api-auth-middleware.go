package api

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/types"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Validate UUID format
		_, err := uuid.FromString(token)
		if err != nil {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		idParam := ps.ByName("userId")
		id, err := strconv.Atoi(idParam)
		authToken := types.BearerToken{
			Token:  token,
			UserID: id,
		}

		valid, err := rt.db.ValidateToken(authToken)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error validating token")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !valid {
			rt.baseLogger.WithField("token", token).Warn("Invalid or expired token")
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next(w, r, ps)
	}
}
