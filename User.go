package main

import "encoding/json"

func (RestAPI *RestAPI) GetUsersCollections() {
	RestAPI.Firebase.Users = RestAPI.Firebase.client.Collection("Users")
}

func (RestAPI *RestAPI) GetEventsCollections() {
	RestAPI.Firebase.Events = RestAPI.Firebase.client.Collection("Events")
}

func (RestAPI *RestAPI) UsersCollectionsJson() (UserJson []User) {
	//implement a function that serlieses the users collection to json by using the type 'User' from models.go and return json object
	col := RestAPI.Firebase.Users
	docs, err := col.Documents(RestAPI.Firebase.ctx).GetAll()
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {
		var user User
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

func (RestAPI *RestAPI) EventsCollectionsJson() (EventJson []Event) {
	//implement a function that serlieses the events collection to json by using the type 'Event' from models.go and return json object
	col := RestAPI.Firebase.Events
	docs, err := col.Documents(RestAPI.Firebase.ctx).GetAll()
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {
		var event Event
		doc.DataTo(&event)
		EventJson = append(EventJson, event)
	}
	return EventJson
}
