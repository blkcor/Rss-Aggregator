package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX err: %v", msg)
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	responseWithJson(w, code, ErrorResponse{
		Error: msg,
	})
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	//get the json string
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to Marshal json response: %v", payload)
		w.WriteHeader(500)
		return
	}
	//write the response data
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("An error happen when responsing the data: %v", err)
		return
	}
}
