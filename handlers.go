package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func (RestAPI *RestAPI) getUserByID(uid string) (User, error) {
	RestAPI.GetUsersCollections()
	users := RestAPI.UsersCollectionsJson()
	intID, err := strconv.Atoi(uid)
	if err != nil {
		return User{}, err
	}
	return users[intID], nil
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&Response{
		StatusCode: statusCode,
		Message:    message,
	})
}

func (RestAPI *RestAPI) GetUsers(w http.ResponseWriter, r *http.Request) {
	RestAPI.GetUsersCollections()
	users := RestAPI.UsersCollectionsJson()
	w.Header().Set("Content-Type", "application/json")

	// Convert the users object to json
	u, err := json.Marshal(users)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "Users fetched successfully",
		Data:       u,
	})
}

func (RestAPI *RestAPI) getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uid := params["id"]

	user, err := RestAPI.getUserByID(uid)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "User fetched successfully",
		Data:       user,
	})
}

func (RestAPI *RestAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User

	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the newUser
	err = newUser.Validate()
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Logic to save newUser in the database
	// ...

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data:       newUser,
	})
}

func (RestAPI *RestAPI) updateUser(w http.ResponseWriter, r *http.Request) {
	// The incoming json should look like this
	// {
	// 	"studentID": 123456,
	// 	"profilePic": "http://example.com/path/to/profile/pic.jpg",
	// 	"bio": "This is a short bio",
	// 	"email": "newEmail@example.com",
	// 	"password": "newPassword",
	// 	"firstName": "John",
	// 	"lastName": "Doe",
	// 	"phone": 1234567890,
	// 	"role": "student",
	// 	"points": 100,
	// 	"github": "https://github.com/username",
	// 	"uid": "uniqueUserID"
	// }
	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	params := mux.Vars(r)
	uid := params["uid"]

	existingUser, err := RestAPI.getUserByID(uid)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//convert the users object to json and write it to the response

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Response{
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		})
		return
	}

	if existingUser.Role != "Administrator" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&Response{
			StatusCode: http.StatusForbidden,
			Message:    "Permission denied",
		})
		return
	}

	// Logic to update the existing user with the updatedUser data
	userValue := reflect.ValueOf(updatedUser)
	userType := reflect.TypeOf(updatedUser)

	var updates []firestore.Update
	for i := 0; i < userValue.NumField(); i++ {
		updates = append(updates, firestore.Update{
			Path:  userType.Field(i).Name,
			Value: userValue.Field(i).Interface(),
		})
	}

	_, err = RestAPI.Firebase.client.Collection("users").Doc(uid).Update(RestAPI.Firebase.ctx, updates)
	// ...

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "User updated successfully",
		Data:       existingUser,
	})
}

func (RestAPI *RestAPI) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uid := params["uid"]

	DeletedUser := User{}
	err := json.NewDecoder(r.Body).Decode(&DeletedUser)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//convert the users object to json and write it to the response
	existingUser, err := RestAPI.getUserByID(uid)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&Response{
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
		})
		return
	}

	if existingUser.Role != "Administrator" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&Response{
			StatusCode: http.StatusForbidden,
			Message:    "Permission denied",
		})
		return
	}

	_, err = RestAPI.Firebase.client.Collection("Users").Doc(DeletedUser.UID).Delete(RestAPI.Firebase.ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "User deleted successfully",
	})
}
