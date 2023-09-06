package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func (RestAPI *RestAPI) initLogger() {
	col := RestAPI.Firebase.client.Collection("Users")
	iter := col.Snapshots(context.Background())
	defer iter.Stop()
	for {
		//Next() call blocks, until changes are received
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			//return err
			log.Printf("next iter error %s\n", err)
		}
		for _, change := range doc.Changes {
			// access the change.Doc returns the Document,
			// which contains Data() and DataTo(&p) methods.
			switch change.Kind {
			case firestore.DocumentAdded:
				log.Printf("User added, %s\n", change.Doc.Ref.Path)
				// on added it returns the existing ones.
			case firestore.DocumentModified:
				log.Printf("User modified, %s\n", change.Doc.Ref.Path)
			case firestore.DocumentRemoved:
				log.Printf("User removed, %s\n", change.Doc.Ref.Path)
			}
		}
	}
}
