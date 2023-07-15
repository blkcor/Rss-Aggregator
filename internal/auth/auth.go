package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
GetApiKey function could extract ApiKey in the request header
example like this:
Authorization: ApiKey {insert your api key here}
*/
func GetApiKey(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}
