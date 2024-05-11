package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

const secretkeyfile = "w5a8T1V+h6aPtVyzYsH2wZkUaxOEYZyKrVYUzRQxC9E="

func generateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userID,
		"exp":   time.Now().Add(time.Minute * 10).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretkeyfile))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyAuthn(h http.Header) (any, error) {
	authHeader := h.Get("Authorization")
	val := strings.TrimPrefix(authHeader, "Bearer ")

	if val == "" {
		return nil, errors.New("AuthHeader Empty")
	}

	token, err := jwt.Parse(val, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkeyfile), nil
	})

	m, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("InvalidToken")
	}
	if err != nil {
		log.Printf("error in validating jwt for %+v, errpr %+v", m["email"], err)
		return nil, err
	}

	return m["email"], nil
}
