package Firebase

import (
	"encoding/json"

	Types "ftcksu.com/api/v2/types"
)

func (Firebase *Firebase) GetUsersCollections() {
	Firebase.Users = Firebase.Client.Collection("Users")
}

func (Firebase *Firebase) GetEventsCollections() {
	Firebase.Events = Firebase.Client.Collection("Events")
}

func (Firebase *Firebase) UsersCollectionsJson() (UserJson []Types.User) {
	//implement a function that serlieses the users collection to json by using the type 'User' from models.go and return json object
	col := Firebase.Users
	docs, err := col.Documents(Firebase.Ctx).GetAll()
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {
		var user Types.User
		data := doc.Data()
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			// handle error
		}
		err = json.Unmarshal(jsonBytes, &user)
		if err != nil {
			// handle error
		}
		UserJson = append(UserJson, user)
	}
	return UserJson
}

func (Firebase *Firebase) EventsCollectionsJson() (EventJson []Types.Event) {
	//implement a function that serlieses the events collection to json by using the type 'Event' from models.go and return json object
	col := Firebase.Events
	docs, err := col.Documents(Firebase.Ctx).GetAll()
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {
		var event Types.Event
		doc.DataTo(&event)
		EventJson = append(EventJson, event)
	}
	return EventJson
}
