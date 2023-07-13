package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	//get the json string
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to Marshal json response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
