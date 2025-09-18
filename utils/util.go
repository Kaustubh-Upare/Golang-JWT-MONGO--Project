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

func CreateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"ttl":    time.Now().Add(time.Hour * 24 * 10).Unix(), //10days
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

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
		// Path:     "/",
	}
}

func ValidateCookie(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("Invalid Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return fmt.Errorf("Could not parse claims")
	}

	if ttl, ok := claims["ttl"].(float64); ok {
		if time.Now().Unix() > int64(ttl) {
			return fmt.Errorf("Token Has Expired")
		}
	} else {
		return fmt.Errorf("ttl claim is not a number")
	}

	// Check For user is there or not
	userIdString := claims["userId"].(string)

	objId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		log.Println("objectId", err)
		return fmt.Errorf("String is not valid")
	}

	userExist, err := models.CheckForUser(objId)
	if err != nil {
		log.Println("userExist", err)

		return fmt.Errorf("Internal Server Error")
	}

	if !userExist {
		return fmt.Errorf("User Not Found")
	}

	return nil

}
