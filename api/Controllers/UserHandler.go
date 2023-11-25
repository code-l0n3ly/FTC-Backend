package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
)

func (RestAPI *RestAPI) getUserByID(uid string) (User, error) {
	RestAPI.GetUsersCollections()
	users := RestAPI.UsersCollectionsJson()
	intID, _ := strconv.Atoi(uid)
	if intID >= len(users) {
		return User{}, fmt.Errorf("User not found")
	}
	U := users[intID]
	return U, nil
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "Users fetched successfully",
		Data:       json.NewEncoder(w).Encode(users),
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

	// Logic to save newUser in the database

	_, err = RestAPI.Firebase.client.Collection("Users").Doc(newUser.UID).Create(RestAPI.Firebase.ctx, newUser)
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
	var NewInfo User
	err := json.NewDecoder(r.Body).Decode(&NewInfo)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if NewInfo.UID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "UID is required")
		return
	}

	params := mux.Vars(r)
	uid := params["id"]

	existingUser, err := RestAPI.getUserByID(uid)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	UserToBeUpdated, err := RestAPI.getUserByID(NewInfo.UID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	if existingUser.Role != "Administrator" {
		sendErrorResponse(w, http.StatusForbidden, "Permission denied")
		return
	}

	userValue := reflect.ValueOf(NewInfo)
	userType := reflect.TypeOf(NewInfo)

	var updates []firestore.Update
	for i := 0; i < userValue.NumField(); i++ {
		// Get the field value and the field type.
		fieldValue := userValue.Field(i)
		fieldType := userType.Field(i)

		// Check if the field value is not the zero value for its type.
		if !fieldValue.IsZero() {
			// If it's not, add an update for this field.
			updates = append(updates, firestore.Update{
				Path:  fieldType.Name,
				Value: fieldValue.Interface(),
			})
		}
	}

	_, err = RestAPI.Firebase.client.Collection("Users").Doc(UserToBeUpdated.UID).Update(RestAPI.Firebase.ctx, updates)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	UserToBeUpdated, err = RestAPI.getUserByID(NewInfo.UID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "User updated successfully",
		Data:       UserToBeUpdated,
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

func (RestAPI *RestAPI) Auth(w http.ResponseWriter, r *http.Request) {
	var LoggedIn User

	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&LoggedIn)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	U, err := RestAPI.getUserByID(LoggedIn.UID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "User logged in successfully",
		Data:       RestAPI.generateBearerToken(LoggedIn.Username, LoggedIn.Password, U.Role),
	})
}

func (RestAPI *RestAPI) FindUserByUsername(username string) (User, error) {
	RestAPI.GetUsersCollections()
	users := RestAPI.UsersCollectionsJson()
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("User not found")
}

func (RestAPI *RestAPI) Login(w http.ResponseWriter, r *http.Request) {
	var LoggedIn User

	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&LoggedIn)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	User, err := RestAPI.FindUserByUsername(LoggedIn.Username)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Username Or password is incorrect")
		return
	}

	if User.Password != LoggedIn.Password {
		sendErrorResponse(w, http.StatusBadRequest, "Username Or password is incorrect")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Response{
		StatusCode: http.StatusOK,
		Message:    "User logged in successfully",
		Data:       User,
	})
}
