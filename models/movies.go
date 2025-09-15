package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `
	Movie  string             `json:"movie"`
	Actors []string           `json:"actors"`
}

func InsertMovie(movie Movie) error {
	collection := mongoClient.Database(db).Collection(collName)
	inserted, err := collection.InsertOne(context.TODO(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a record with id: ", inserted.InsertedID)
	return err
}
func InsertMany(movies []Movie) error {

	newMovies := make([]interface{}, len(movies))
	for i, movie := range movies {
		newMovies[i] = movie
	}
	collection := mongoClient.Database(db).Collection(collName)
	result, err := collection.InsertMany(context.TODO(), newMovies)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
	return err
}

func UpdateMovie(movieId string, movie Movie) error {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"movie": movie.Movie, "actors": movie.Actors}}
	collection := mongoClient.Database(db).Collection(collName)
	result, err := collection.UpdateOne(context.TODO(), filter, update)

	log.Println("New Record ", result)
	return err
}

func deleteMovie(movieId string) error {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	collection := mongoClient.Database(db).Collection(collName)
	result, err := collection.DeleteOne(context.TODO(), filter)
	log.Println("Deleted Succesfully ", result)
	return err
}
