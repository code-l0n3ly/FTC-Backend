package Firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Firebase struct {
	Ctx    context.Context
	Client *firestore.Client
	Users  *firestore.CollectionRef
	Events *firestore.CollectionRef
}

func New(Credntials string, AppName string) *Firebase {
	ctx := context.Background()
	opt := option.WithCredentialsJSON([]byte(Credntials))
	client, err := firestore.NewClient(ctx, AppName, opt)
	if err != nil {
		log.Fatalf("firestore new error:%s\n", err)
	}
	Firebase := &Firebase{
		Ctx:    ctx,
		Client: client,
	}
	return Firebase
}
