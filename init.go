package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

type RestAPI struct {
	Router   *mux.Router
	Firebase *Firebase
}

type Firebase struct {
	ctx    context.Context
	client *firestore.Client
}

func New() *RestAPI {
	ctx := context.Background()
	opt := option.WithCredentialsFile("service_account.json")
	client, err := firestore.NewClient(ctx, "ftc-app-36fad", opt)
	if err != nil {
		log.Fatalf("firestore new error:%s\n", err)
	}
	RestAPI := &RestAPI{
		Firebase: &Firebase{
			ctx:    ctx,
			client: client,
		},
	}

	//
	RestAPI.Router = RestAPI.initializeRouter()

	return RestAPI
}

func main() {
	//Init RestAPI and Router with Firebase client
	RestAPI := New()
	go RestAPI.initLogger()
	//Start server
	log.Fatal(http.ListenAndServe(":8000", RestAPI.Router))
}
