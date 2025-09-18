package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kaustubh-upare/jwtWithMongo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateToken(userId primitive.ObjectID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId.Hex(),
		"ttl":    time.Now().Add(time.Hour * 24 * 10).Unix(), //10days
	})

	// Sign and get the complete encoded token as a string using the secret
	log.Println("Token being created with userId string: ", userId)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	log.Printf("Token being created with userId string2 Token String: ", tokenString)
	return tokenString, err
}

func CookieBoiler(auth string, token string) *http.Cookie {
	return &http.Cookie{
		Name:     auth,
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 10), //10days
		HttpOnly: true,                                // Prevents JavaScript from accessing the cookie
		Secure:   true,                                // Ensures cookie is only sent over HTTPS
		SameSite: http.SameSiteLaxMode,                // Prevents cross-site request forgery
		Path:     "/",
	}
}

func ValidateCookie(tokenString string) (primitive.ObjectID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	log.Println("Parsed Value  ", token)

	if err != nil {
		log.Println("2nd step method", err)
		return primitive.NilObjectID, err
	}

	if !token.Valid {

		log.Println("3rd step method", err)
		return primitive.NilObjectID, fmt.Errorf("Invalid Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {

		log.Println("4th step method", err)
		return primitive.NilObjectID, fmt.Errorf("Could not parse claims")
	}

	if ttl, ok := claims["ttl"].(float64); ok {
		if time.Now().Unix() > int64(ttl) {

			log.Println("ttl step method", err)
			return primitive.NilObjectID, fmt.Errorf("Token Has Expired")
		}
	} else {
		return primitive.NilObjectID, fmt.Errorf("ttl claim is not a number")
	}

	// Check For user is there or not
	userIdString := claims["userId"].(string)
	log.Printf("Validating token with userId string: %s", userIdString)
	objId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		log.Println("objectId", err)
		return primitive.NilObjectID, fmt.Errorf("String is not valid")
	}

	userExist, err := models.CheckForUser(objId)
	if err != nil {
		log.Println("userExist", err)

		return primitive.NilObjectID, fmt.Errorf("Internal Server Error")
	}

	if !userExist {
		return primitive.NilObjectID, fmt.Errorf("User Not Found")
	}

	return objId, nil

}
