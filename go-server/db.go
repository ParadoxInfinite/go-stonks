package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Variables with a capital letter are considered publc in go and are available in other files/packages.
// @see https://stackoverflow.com/questions/58365838/how-to-reuse-mongodb-connection-in-go#comment116916095_59119371
var CNX = Connection()
var CTX context.Context;
func Connection() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongodb_uri := os.Getenv("MONGO_DB_URI")
	// Set client options
	clientOptions := options.Client().ApplyURI(mongodb_uri)
	CTX, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Connect to MongoDB
	client, err := mongo.Connect(CTX, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}
