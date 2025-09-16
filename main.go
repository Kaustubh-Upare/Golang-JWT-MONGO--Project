package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaustubh-upare/jwtWithMongo/controllers"
	"github.com/kaustubh-upare/jwtWithMongo/middleware"
	"github.com/kaustubh-upare/jwtWithMongo/models"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Bhai")
	val := os.Getenv("some")
	log.Printf("SOME=%q\n", val)
	fmt.Println(val)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database
	models.ConnectToDatabase()

	fmt.Println("hello 10")
	handlers := controllers.NewMovieHandler()

	router := http.NewServeMux()

	router.HandleFunc("GET /pring", ping)

	// Create
	router.HandleFunc("POST /movies", handlers.Create)
	router.HandleFunc("POST /movies/bulk", handlers.CreateMany)
	router.HandleFunc("GET /movies/test", handlers.TestBulk)

	// Read
	router.Handle("GET /movies", middleware.TestMiddleware(http.HandlerFunc(handlers.ListAll)))
	router.HandleFunc("GET /movies/search", handlers.GetByName) //?name=foo

	//Update By ID
	router.HandleFunc("PATCH /movies/{id}", handlers.Update) //user/{id}

	// DELETION
	router.HandleFunc("DELETE /movies/{id}", handlers.DeleteOne)
	router.HandleFunc("DELETE /movies", handlers.DeleteAll)

	log.Println("Listening on port 8080")
	errL := http.ListenAndServe(":8080", router)
	if errL != nil {
		log.Fatal(errL)
	}
}
