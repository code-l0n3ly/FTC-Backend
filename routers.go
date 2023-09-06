package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (RestAPI *RestAPI) initializeRouter() *mux.Router {
	router := mux.NewRouter()
	// Users API endpoints - TODO: Add authentication middleware
	router.HandleFunc("/users", RestAPI.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", RestAPI.getUser).Methods(http.MethodGet)
	router.HandleFunc("/users", RestAPI.createUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", RestAPI.updateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", RestAPI.deleteUser).Methods(http.MethodDelete)
	return router
}
