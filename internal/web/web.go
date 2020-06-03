package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	headerContentType   = "Content-Type"
	mimeApplicationJSON = "application/json"
)

// Respond create json response and outputs json representation of the passed data into response body.
func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set(headerContentType, mimeApplicationJSON)
	w.WriteHeader(status)
	if data != nil {
		EncodeBody(w, r, data)
	}
}

// RespondError create json error response and outputs passed error into response body.
func RespondError(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {
	Respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...)},
	})
}

// EncodeBody encodes passed date to json format and writes it into Response body.
func EncodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// DecodeBody decode json from request body into passed pointer struct.
func DecodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
