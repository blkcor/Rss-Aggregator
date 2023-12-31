package main

import (
	"encoding/json"
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (a *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}
	feed, err := a.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error creating feed:%v", err))
		return
	}
	responseWithJson(w, 200, dbFeedToFeed(feed))
}

func (a *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.DB.GetFeeds(r.Context())
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error getting feeds:%v", err))
		return
	}
	feedsResponse := make([]Feed, 0)
	for _, feed := range feeds {
		feedsResponse = append(feedsResponse, dbFeedToFeed(feed))
	}
	responseWithJson(w, 200, feedsResponse)
}
