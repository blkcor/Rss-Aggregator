package main

import (
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/auth"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"net/http"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (a *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			responseWithError(w, 403, fmt.Sprintf("Auth err: %v", err))
			return
		}
		user, err := a.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			responseWithError(w, 401, fmt.Sprintf("Error Getting User: %v", err))
			return
		}
		handler(w, r, user)
	}
}
