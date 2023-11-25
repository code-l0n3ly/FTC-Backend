package main

import (
	"log"
	"net/http"

	Firebase "ftcksu.com/ftcksu/ftcapp/firebase"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	SECRET_KEY string
	USERS_PATH string = "/api/users"
)

func initializeRouter(RestAPI *Firebase.RestAPI) *mux.Router {
	router := mux.NewRouter()
	SECRET_KEY = uuid.New().String()
	//USERS ENDPOINTS
	router.HandleFunc(USERS_PATH, RestAPI.authenticate(RestAPI.GetUsers)).Methods(http.MethodGet)
	router.HandleFunc(USERS_PATH+"/{id}", RestAPI.authenticate(RestAPI.getUser)).Methods(http.MethodGet)
	router.HandleFunc(USERS_PATH+"/create", RestAPI.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/login", RestAPI.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/Auth", RestAPI.Auth).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/update/{id}", RestAPI.authenticate(RestAPI.updateUser)).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/delete/{id}", RestAPI.authenticate(RestAPI.deleteUser)).Methods(http.MethodPost)

	return router
}

func initLogger() {
	// file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(file)
	// log.Println("Log file initialized")
}

func StartServer(RestAPI *Firebase.RestAPI) {
	//Init RestAPI and Router with Firebase client
	go RestAPI.initLogger()
	//Start server
	log.Fatal(http.ListenAndServe(":8000", RestAPI.Router))
}
