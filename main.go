package main

import (
	"fmt"
	"log"
	"net/http"

	Firebase "ftcksu.com/api/v2/Firebase"
	Routes "ftcksu.com/api/v2/api/Routes"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type RestAPI struct {
	Router   *mux.Router
	Firebase *Firebase.Firebase
}

// Start server
func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func main() {
	AuthToken := viper.Get("AUTH_TOKEN")
	FirebaseCredentials := viper.Get("FIREBASE_CREDENTIALS")
	FirebaseAppName := viper.Get("FIREBASE_APP_NAME")
	fmt.Println(AuthToken.(string))
	//Init RestAPI and Router with Firebase client
	FB := Firebase.New(FirebaseCredentials.(string), FirebaseAppName.(string))
	RestAPI := &RestAPI{
		Router:   mux.NewRouter(),
		Firebase: FB,
	}
	//Bind routes
	Routes.BindUserRoutes(RestAPI.Router, RestAPI.Firebase, AuthToken.(string))
	//Init logger :: TODO
	//go RestAPI.InitLogger()
	//Start server
	log.Fatal(http.ListenAndServe(":8000", RestAPI.Router))
}
