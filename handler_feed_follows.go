package main

import (
	"encoding/json"
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (a *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	param := parameter{}
	err := decoder.Decode(&param)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}
	feedFollows, err := a.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    param.FeedID,
	})
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Fail to create FeedsFollows: %v", err))
		return
	}
	responseWithJson(w, 200, dbFeedFollowToFeedFollow(feedFollows))
}

func (a *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := a.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Fail to get FeedsFollows: %v", err))
		return
	}
	feedFollowResponse := make([]database.FeedFollow, 0)
	for _, feedFollow := range feedFollows {
		feedFollowResponse = append(feedFollowResponse, feedFollow)
	}
	responseWithJson(w, 200, feedFollowResponse)
}
