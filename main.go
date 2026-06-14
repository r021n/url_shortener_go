package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("GET /demo/greet", greetHandler)
	mux.HandleFunc("POST /demo/greet", greetPostHandler)
	mux.HandleFunc("GET /demo/users/{name}", userHandler)

	mux.HandleFunc("/", notFoundHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// --- handlers ---

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"message": "URL Shortener API is running",
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Hello! Send a POST request to this endpoint with your name.",
	})
}

func greetPostHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "Invalid JSON. send a body like: {\"name\": \"YourName\"}",
		})
		return
	}
	if body.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "The 'name' field is required",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Hello " + body.Name + "! Welcome to the URL Shortener project.",
	})
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	writeJSON(w, http.StatusOK, map[string]any{
		"user": name,
		"info": "This endpoint extracts the path parameter using r.PathValue().",
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotFound, map[string]any{
		"error": "Route not found.",
	})
}

// --- helper ---
func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
