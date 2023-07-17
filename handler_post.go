package main

import (
	"encoding/json"
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"net/http"
	"strconv"
)

func (a *apiConfig) handlerGetUserPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		Limit string `json:"limit"`
	}
	param := parameter{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error parsing Json Body:%v", err))
		return
	}
	convertedLimit, err := strconv.Atoi(param.Limit)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error parsing Json parameter:%v,please provide the correct parameter", err))
		return
	}
	posts, err := a.DB.GetPostForUser(r.Context(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(convertedLimit),
	})
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error getting posts for user:%v", err))
		return
	}
	postsResponse := make([]Post, 0)
	for _, post := range posts {
		postsResponse = append(postsResponse, dbPostToPost(post))
	}
	responseWithJson(w, 200, postsResponse)
}
