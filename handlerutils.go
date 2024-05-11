package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// loginHandler is a Handler for login requests. Request Body contains user password and email.
// Handler verifies the password and issues a jwt token
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user, dbObj User
	err := decodeUserFromBody(r.Body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Login request for user %+v", user.Email)
	ok := verifyDetails(user.Email, user.Password, &dbObj)
	if !ok {
		log.Printf("%+v", ok)
		fmt.Fprint(w, "User and password not found. Verify the credentials.")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	token, err := generateJWT(user.Email) // issues jwt token
	fmt.Fprint(w, token)
	return
}

// registerHandler is a Handler for login requests
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := decodeUserFromBody(r.Body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Server error, unable to create your account.", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword
	err = insert(r.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			log.Printf("key already exists %+v", err)
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Account with Email Id %+v already exists , try with different email", user.Email)
			return
		}
		log.Printf("err in registering user %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprint(w, "Registered Successfully")
	w.WriteHeader(http.StatusOK)
	return
}

// helloHandler display about page
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to Auth Server Demo")
	return
}

// resthandler displays rest resource after authorization
func resthandler(w http.ResponseWriter, r *http.Request) {
	useremail, err := verifyAuthn(r.Header) //config file
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "The access token did not work, please issue a new token")
		return
	}
	fmt.Fprintf(w, "authorization successful for %+v", useremail)
	return
}
