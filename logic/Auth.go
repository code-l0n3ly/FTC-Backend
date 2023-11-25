package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (RestAPI *RestAPI) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		uid := params["id"]
		if len(uid) == 0 {
			sendErrorResponse(w, http.StatusBadRequest, "Missing required Data")
			return
		}
		RestAPI.GetUsersCollections()
		users := RestAPI.UsersCollectionsJson()
		intID, _ := strconv.Atoi(uid)
		if intID >= len(users) {
			sendErrorResponse(w, http.StatusBadRequest, "Wrong required fields")
		}
		LoginUser := users[intID]
		fmt.Printf("Login attempt : %v\n", LoginUser)
		authHeader := r.Header.Get("Authorization")
		// Check if the Authorization header is empty or doesn't start with "Bearer".
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get the token from the Authorization header.
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token. This is a placeholder and should be replaced with your actual token validation logic.
		if RestAPI.CheckBearerToken(LoginUser.Username, LoginUser.Password) != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If the token is valid, call the next handler.
		next.ServeHTTP(w, r)
	}
}

func (RestAPI *RestAPI) CheckBearerToken(username, password string) string {
	// Concatenate the username, password, and role with a separator
	data := username + ":" + password + ":" + "Administrator" + ":" + SECRET_KEY

	// Generate a SHA-256 hash of the data
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	token := fmt.Sprintf("%x", hash)

	return token
}

func (RestAPI *RestAPI) generateBearerToken(username, password, role string) string {
	// Concatenate the username, password, and role with a separator
	data := username + ":" + password + ":" + role + ":" + SECRET_KEY

	// Generate a SHA-256 hash of the data
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	token := fmt.Sprintf("%x", hash)

	return token
}
