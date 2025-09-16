package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaustubh-upare/jwtWithMongo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct{}

func NewMovieHandler() *MovieHandler {
	return &MovieHandler{}
}

func (h *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	var m models.Movie
	// fmt.Println("hello 1")
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	// fmt.Println("hello 2")
	if err := models.InsertMovie(m); err != nil {
		http.Error(w, "insert failed", http.StatusInternalServerError)
		return
	}

	// fmt.Println("hello 3")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (h *MovieHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var ms []models.Movie

	if err := json.NewDecoder(r.Body).Decode(&ms); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := models.InsertMany(ms); err != nil {
		http.Error(w, "Insertion Failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "created": "Succesfully"})
}

func (h *MovieHandler) GetByName(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "missing ?name", http.StatusBadRequest)
		return
	}

	mv := models.Find(name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "movie": mv})
}

func (h *MovieHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	mvs := models.ListAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mvs)
}
func (h *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var m models.Movie

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := models.UpdateMovie(id, m); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "Update": "succesfully"})
}

func (h *MovieHandler) TestBulk(w http.ResponseWriter, r *http.Request) {
	var ms []models.Movie

	for i := 0; i < 100000; i++ {
		movie := models.Movie{
			Movie: fmt.Sprintf("Movie Title yes%d", i),
			Actors: []string{
				fmt.Sprintf("Actor A%d", i),
				fmt.Sprintf("Actor B%d", i),
			},
		}
		ms = append(ms, movie)
	}

	err := models.InsertMany(ms)
	if err != nil {
		http.Error(w, "Bulk Insertion failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"success": true, "created": "Bulk Created Succesfully"})
}

func (h *MovieHandler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := models.DeleteMovie(id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true, "success": "deleted Successfully"})
}

func (h *MovieHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := models.DeleteAll(); err != nil {
		http.Error(w, "Error While Deleting", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true, "success": "deleted Successfully"})
}
