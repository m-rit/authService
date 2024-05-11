package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dbname = "authdb"
const usercollname = "users"

func initDB(ctx context.Context) error {

	uri := "mongodb://localhost:27017"
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri))
	if err != nil {
		return err
	}

	gdbclient = client
	coll := getCollection(dbname, usercollname)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic(err)
	}

	/*defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()*/
	return nil
}

func getCollection(dbName, collName string) *mongo.Collection {
	database := gdbclient.Database(dbname)
	coll := database.Collection(collName)
	return coll
}

func findOne(key string, result interface{}) error {
	coll := getCollection(dbname, usercollname)
	err := coll.FindOne(context.TODO(), bson.D{{"email", key}}).
		Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the email %s\n", key)
		return err
	}
	if err != nil {
		panic(err)
	}
	return nil
}

func insert(ctx context.Context, u User) error {
	coll := getCollection(dbname, usercollname)
	log.Printf("inserting %+v", u)
	_, err := coll.InsertOne(ctx, u)
	if err != nil {
		fmt.Println("Something went wrong trying to insert the new documents:")
		return err
	}
	return nil
}
