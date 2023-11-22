package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	SECRET_KEY string
	USERS_PATH string = "/api/users"
)

func (RestAPI *RestAPI) initializeRouter() *mux.Router {
	router := mux.NewRouter()
	SECRET_KEY = uuid.New().String()
	//authentication not required Users API endpoints
	router.HandleFunc(USERS_PATH, RestAPI.authenticate(RestAPI.GetUsers)).Methods(http.MethodGet)
	router.HandleFunc(USERS_PATH+"/{id}", RestAPI.authenticate(RestAPI.getUser)).Methods(http.MethodGet)
	router.HandleFunc(USERS_PATH+"/create", RestAPI.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/login", RestAPI.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/Auth", RestAPI.Auth).Methods(http.MethodPost)

	//authentication required Users API endpoints
	router.HandleFunc(USERS_PATH+"/update/{id}", RestAPI.authenticate(RestAPI.updateUser)).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/delete/{id}", RestAPI.authenticate(RestAPI.deleteUser)).Methods(http.MethodPost)

	return router
}
