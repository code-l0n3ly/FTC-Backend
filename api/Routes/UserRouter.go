package Routes

import (
	"net/http"

	Firebase "ftcksu.com/api/v2/Firebase"
	UserHandlers "ftcksu.com/api/v2/api/Controllers"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	SECRET_KEY string
	USERS_PATH string = "/api/users"
)

func BindUserRoutes(router *mux.Router, FB *Firebase.Firebase, AuthToken string) {
	SECRET_KEY = uuid.New().String()
	UserHandler := UserHandlers.New(FB, AuthToken)
	//USERS ENDPOINTS
	router.HandleFunc(USERS_PATH, UserHandler.GetUsers).Methods(http.MethodGet)
	router.HandleFunc(USERS_PATH+"/{id}", UserHandler.GetUser).Methods(http.MethodGet)
	router.HandleFunc(USERS_PATH+"/create", UserHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/login", UserHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/Auth", UserHandler.Auth).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/update/{id}", UserHandler.UpdateUser).Methods(http.MethodPost)
	router.HandleFunc(USERS_PATH+"/delete/{id}", UserHandler.DeleteUser).Methods(http.MethodPost)
}
