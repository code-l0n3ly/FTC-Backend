package Controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"cloud.google.com/go/firestore"
	Auth "ftcksu.com/api/v2/Authentication"
	"ftcksu.com/api/v2/Firebase"
	Types "ftcksu.com/api/v2/types"
	"github.com/gorilla/mux"
)

type UserController struct {
	FB             *Firebase.Firebase
	Authentication *Auth.JWT
}

func New(FB *Firebase.Firebase, AuthToken string) *UserController {
	UserController := &UserController{}
	UserController.FB = FB
	UserController.Authentication = &Auth.JWT{Key: []byte(AuthToken)}
	fmt.Println("Firebase connected to User Handlers")
	return UserController
}
func (UserController *UserController) getUserByID(uid string) (Types.User, error) {
	UserController.FB.GetUsersCollections()
	users := UserController.FB.UsersCollectionsJson()
	intID, _ := strconv.Atoi(uid)
	if intID >= len(users) {
		return Types.User{}, fmt.Errorf("User not found")
	}
	U := users[intID]
	return U, nil
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&Types.Response{
		StatusCode: statusCode,
		Message:    message,
	})
}

func (UserController *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	UserController.FB.GetUsersCollections()
	users := UserController.FB.UsersCollectionsJson()
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Types.Response{
		StatusCode: http.StatusOK,
		Message:    "Users fetched successfully",
		Data:       json.NewEncoder(w).Encode(users),
	})
}

func (UserController *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uid := params["id"]

	user, err := UserController.getUserByID(uid)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&Types.Response{
		StatusCode: http.StatusOK,
		Message:    "User fetched successfully",
		Data:       user,
	})
}

func (UserController *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser Types.User

	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Logic to save newUser in the database
	//_, err = UserController.FB.Users.(newUser.UID).Create(UserController.FB.Ctx, newUser)
	_, _, err = UserController.FB.Users.Add(UserController.FB.Ctx, newUser)
	// ...

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Types.Response{
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data:       newUser,
	})
}

func (UserController *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if UserController.Authentication.ValidateToken(r.Header.Get("Authorization")) != nil {
		var NewInfo Types.User
		err := json.NewDecoder(r.Body).Decode(&NewInfo)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if NewInfo.UID == "" {
			sendErrorResponse(w, http.StatusBadRequest, "UID is required")
			return
		}

		UserToBeUpdated, err := UserController.getUserByID(NewInfo.UID)
		if err != nil {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
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
		_, err = UserController.FB.Client.Collection("Users").Doc(UserToBeUpdated.UID).Update(UserController.FB.Ctx, updates)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		UserToBeUpdated, err = UserController.getUserByID(NewInfo.UID)
		if err != nil {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Types.Response{
			StatusCode: http.StatusOK,
			Message:    "User updated successfully",
			Data:       UserToBeUpdated,
		})
	} else {
		sendErrorResponse(w, http.StatusForbidden, "Permission denied")
		return
	}
}

func (UserController *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if UserController.Authentication.ValidateToken(r.Header.Get("Authorization")) != nil {
		DeletedUser := Types.User{}
		err := json.NewDecoder(r.Body).Decode(&DeletedUser)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		//convert the users object to json and write it to the response
		DeletedUser, err = UserController.getUserByID(DeletedUser.UID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&Types.Response{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
			})
			return
		}
		_, err = UserController.FB.Client.Collection("Users").Doc(DeletedUser.UID).Delete(UserController.FB.Ctx)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Types.Response{
			StatusCode: http.StatusOK,
			Message:    "User deleted successfully",
		})
	} else {
		sendErrorResponse(w, http.StatusForbidden, "Permission denied")
		return
	}
}
func (UserController *UserController) Auth(w http.ResponseWriter, r *http.Request) {
	var LoggedIn Types.User

	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&LoggedIn)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//U, err := UserController.getUserByID(LoggedIn.UID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	LoggedIn, err = UserController.FindUserByUsername(LoggedIn.Username)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Username Or password is incorrect")
		return
	}
	if LoggedIn.Role != "Admin" {

		AuthToken, err := UserController.Authentication.GenerateJWT(LoggedIn.Username + ":" + LoggedIn.Password)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&Types.Response{
			StatusCode: http.StatusOK,
			Message:    "Authenticated",
			Data:       AuthToken,
		})
	} else {
		sendErrorResponse(w, http.StatusBadRequest, "Action not Allowed")
		return
	}
}

func (UserController *UserController) FindUserByUsername(username string) (Types.User, error) {
	UserController.FB.GetUsersCollections()
	users := UserController.FB.UsersCollectionsJson()
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return Types.User{}, fmt.Errorf("User not found")
}

func (UserController *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var LoggedIn Types.User

	// Decode the incoming User json
	err := json.NewDecoder(r.Body).Decode(&LoggedIn)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	User, err := UserController.FindUserByUsername(LoggedIn.Username)
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
	json.NewEncoder(w).Encode(&Types.Response{
		StatusCode: http.StatusOK,
		Message:    "User logged in successfully",
		Data:       User,
	})
}
