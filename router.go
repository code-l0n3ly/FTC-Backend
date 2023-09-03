package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() *mux.Router {
	router := mux.NewRouter()

	// Register API endpoints
	router.HandleFunc("/users", GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/users", createUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", updateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", deleteUser).Methods(http.MethodDelete)

	return router
}
