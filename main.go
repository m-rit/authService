package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

var gdbclient *mongo.Client

const email = "email"
const password = "password"

type User struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

func main() {
	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)
	err := initDB(ctx)
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/about", helloHandler).Methods("GET")
	r.HandleFunc("/api/v1/books/", resthandler).Methods("GET")
	log.Println("Server started on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))

}
