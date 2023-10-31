package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (RestAPI *RestAPI) GetUsers(w http.ResponseWriter, r *http.Request) {
	RestAPI.GetUsersCollections()
	users := RestAPI.UsersCollectionsJson()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//convert the users object to json and write it to the response
	u, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(fmt.Sprintf("%s", u)))
}

func (RestAPI *RestAPI) getUser(w http.ResponseWriter, r *http.Request) {
	RestAPI.GetUsersCollections()
	users := RestAPI.UsersCollectionsJson()
	params := mux.Vars(r)
	intID, err := strconv.Atoi(params["id"])
	if err != nil {
		// handle error
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//convert the users object to json and write it to the response
	u, err := json.Marshal(users[intID])
	if err != nil {
		panic(err)
	}
	w.Write([]byte(fmt.Sprintf("%s", u)))
}

func (RestAPI *RestAPI) createUser(w http.ResponseWriter, r *http.Request) {
	// Logic to create a new user in the database
	// ...
	// Return the response
}

func (RestAPI *RestAPI) updateUser(w http.ResponseWriter, r *http.Request) {
	// Logic to update a specific user in the database
	// ...
	// Return the response
}

func (RestAPI *RestAPI) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Logic to delete a specific user from the database
	// ...
	// Return the response
}
