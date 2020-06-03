package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/remisb/jwt-api/internal/token"
	"github.com/remisb/jwt-api/internal/web"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	port = 8000
	secret = "default remis jwt secret"
	id = "id-value"
)

var tgen = token.NewTokenHmacSha(secret)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	//r.Route("/api/v1/", func(r chi.Router) {
	//	r.Mount("/tags", tag.InitRouter())
	//})

	address := fmt.Sprintf(":%d", port)
	log.Printf("Server started at %s", address)
	//http.ListenAndServe(address, r)

	r.Get("/ping", pong)
	r.Get("/token", generateToken)
	r.Get("/protected", validateToken)

	log.Fatal(http.ListenAndServe(address, r))
}

func mid(w http.ResponseWriter, r *http.Request) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		fmt.Println("Malformed token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Malformed Token"))
	}
}

func generateToken(w http.ResponseWriter, r *http.Request) {
	// Generate token with id
	token := tgen.Generate(id)
	response := make(map[string]string)
	response[id] = token

	id2 := "idstring2"
	token2 := tgen.Generate(id2)
	response[id2] = token2
	web.Respond(w, r, http.StatusOK, response)
}

func validateToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err := errors.New("authorization data is not provided")
		web.RespondError(w, r, http.StatusUnauthorized, err)
		return
	}

	auth := strings.Split(authHeader, "Bearer ")
	if len(auth) != 2 {
		err := errors.New("malformed authorization token")
		web.RespondError(w, r, http.StatusBadGateway, err)
	}

	jwtToken := auth[1]
	// Check if the token is valid for a given duration
	valid, id, issueTime := tgen.Valid(jwtToken, 10 * time.Minute)

	response := make(map[string]string)
	response["valid"] = strconv.FormatBool(valid)
	response["id"] = id
	response["issueTime"] = issueTime.String()
	web.Respond(w, r, http.StatusOK, response)

}

func pong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
