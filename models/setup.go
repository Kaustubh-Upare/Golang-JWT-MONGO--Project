package models

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const db = "movies"
const collName = "movies"

var mongoClient *mongo.Client

func ConnectToDatabase() {
	clientOption := options.Client().ApplyURI(os.Getenv("DBCONNECTION"))

	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		panic(err)
	}
	mongoClient = client
}
