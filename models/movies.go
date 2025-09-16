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

// "bson.D" is to find in order Elements and "bson.M" for unorder
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
	_, err := collection.InsertMany(context.TODO(), newMovies)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(result)
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

func DeleteMovie(movieId string) error {
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

func Find(movieName string) Movie {

	var result Movie

	filter := bson.D{{"movie", movieName}}
	collection := mongoClient.Database(db).Collection(collName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func FindAll(movieName string) []Movie {
	var result []Movie

	filter := bson.D{{"movie", movieName}}
	collection := mongoClient.Database(db).Collection(collName)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(context.TODO(), result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func ListAll() []Movie {
	var result []Movie

	filter := bson.M{} //This means no condition (select *)
	collection := mongoClient.Database(db).Collection(collName)

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func DeleteAll() error {
	collection := mongoClient.Database(db).Collection(collName)

	delResult, err := collection.DeleteMany(context.TODO(), bson.D{{}}, nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Records Deleted", delResult.DeletedCount)
	return err
}
