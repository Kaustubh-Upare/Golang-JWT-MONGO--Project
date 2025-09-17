package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password" `
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func CreateUser(user User) (primitive.ObjectID, error) {
	Collection := mongoClient.Database(db).Collection("users")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("could not hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	result, err := Collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	fmt.Println("Created a user: ", result.InsertedID)
	uid := result.InsertedID.(primitive.ObjectID)
	return uid, nil
}

func GetUser(email string) (User, error) {
	var result User

	filter := bson.D{{"email", email}}
	Collection := mongoClient.Database(db).Collection("users")

	err := Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("User Not Found")
		}
		return User{}, fmt.Errorf("could not get user: %w", err)
	}

	return result, nil
}

func UpdateUser(email string, user User) error {

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"name": user.Name, "password": user.Password, "createdAt": time.Now(), "email": user.Email}}

	Collection := mongoClient.Database(db).Collection("users")
	result, err := Collection.UpdateOne(context.TODO(), filter, update)

	log.Println(result)
	return err
}

func ValidateUser(email string, password string) (bool, error) {
	var result User
	filter := bson.D{{"email", email}}
	Collection := mongoClient.Database(db).Collection("users")

	err := Collection.FindOne(context.TODO(), filter).Decode(&result)
	log.Println("bc password", password, "yee", result.Password)
	log.Println(1)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(2)
			return false, nil
		}
		log.Println(3)
		return false, fmt.Errorf("database error: %w", err)
	}
	log.Println(4)
	// log.Println("resukt", result)
	np := []byte(result.Password)
	p := []byte(password)
	err = bcrypt.CompareHashAndPassword(np, p)
	log.Println("bcrypt")
	// log.Println(bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)))

	log.Println(5)

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Println("Wrong something bc", err)
		}
		log.Println(6)
		return false, nil //password do not match
	}
	log.Println(7)
	return true, nil //password match
}
