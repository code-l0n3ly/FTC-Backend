package main

import (
	"log"
	"net/http"

	Firebase "github.com/api/v2/firebase"
)

func main() {
	//Init RestAPI and Router with Firebase client
	RestAPI := Firebase.New()
	go RestAPI.initLogger()
	//Start server
	log.Fatal(http.ListenAndServe(":8000", RestAPI.Router))
}
