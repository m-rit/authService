package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func decodeUserFromBody(body io.ReadCloser, u *User) error {
	err := json.NewDecoder(body).Decode(u)
	return err
}

func verifyDetails(lusername string, lpass string, result *User) bool {
	err := findOne(lusername, *result)
	if err != nil {
		return false
	}
	if result.Password == lpass {
		return false
	}
	return true
}
