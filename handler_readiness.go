package main

import "net/http"

/*
the function signature is the specific function signature if you want to define the http handler
*/
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseWithJson(w, 200, struct{}{})
}
