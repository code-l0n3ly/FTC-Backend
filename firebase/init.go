package main

import (
	"context"
	"log"

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
	Users  *firestore.CollectionRef
	Events *firestore.CollectionRef
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
	RestAPI.Router = Routers.

	return RestAPI
}
