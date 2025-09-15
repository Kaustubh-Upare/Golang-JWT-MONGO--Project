package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

	router := http.NewServeMux()

	router.HandleFunc("GET /pring", ping)

	log.Println("Listening on port 8080")
	errL := http.ListenAndServe(":8080", router)
	if errL != nil {
		log.Fatal(errL)
	}
}
