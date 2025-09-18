package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kaustubh-upare/jwtWithMongo/models"
	"github.com/kaustubh-upare/jwtWithMongo/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct{}

// Constructor for the Handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	uid, err := models.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	// Jwt Authentication
	token, err := utils.CreateToken(uid.String())
	if err != nil {
		log.Println("token failed", err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	//set token in cookie
	cookie := utils.CookieBoiler("auth", token)

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email query parameter is required", http.StatusBadRequest)
		return
	}

	reqUser, err := models.GetUser(email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqUser)
}

func (h *UserHandler) ValidateUser(w http.ResponseWriter, r *http.Request) {

	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("valid", loginData)
	isValid, err, userId := models.ValidateUser(loginData.Email, loginData.Password)
	if err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if isValid {
		token, err := utils.CreateToken(userId.String())
		if err != nil {
			http.Error(w, "Error in Token Creation", http.StatusInternalServerError)
			return
		}

		cookie := utils.CookieBoiler("auth", token)

		http.SetCookie(w, cookie)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	} else {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	}
}
