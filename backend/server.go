package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Struct definitions for API responses
type DefaultResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type PersonResponse struct {
	Message []Person `json:"message"`
	Status  int      `json:"status"`
}

type WorldResponse struct {
	Message World `json:"message"`
	Status  int    `json:"status"`
}

type BuildingResponse struct {
	Message []Building `json:"message"`
	Status  int        `json:"status"`
}

var (
	requestMap sync.Map // sync.Map is safer for concurrent use
)

// CORS middleware to handle CORS headers and preflight requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Handler for /people endpoint
func personHandler(w http.ResponseWriter, r *http.Request) {
	// Ignore favicon requests
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	logRequest(r)

	persons := getPersons()
		
	if len(persons) == 0 {
		createNewPerson()
	}

	response := PersonResponse{
		Message: persons,
		Status:  200,
	}

	writeJSONResponse(w, response)
}

// Hander for /buildings endpoint
func buildingHandler(w http.ResponseWriter, r *http.Request) {
	// Ignore favicon requests
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	logRequest(r)

	buildings := getBuildings()

	if len(buildings) == 0 {
		createNewBuilding(House, "House 1", Location{0, 0})
	}

	response := BuildingResponse{
		Message: buildings,
		Status:  200,
	}

	writeJSONResponse(w, response)
}

// Handler for /world endpoint
func worldHandler(w http.ResponseWriter, r *http.Request) {
	// Ignore favicon requests
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	logRequest(r)

	world := getWorld()
	if world.Tiles == nil {
		world = createNewWorld(10, 10)
	}

	response := WorldResponse{
		Message: world,
		Status:  200,
	}

	writeJSONResponse(w, response)
}

// Default handler for root and undefined paths
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// Ignore favicon requests
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	if isDuplicateRequest(r) {
		http.Error(w, "Duplicate request detected", http.StatusTooManyRequests)
		return
	}

	logRequest(r)

	response := DefaultResponse{
		Message: "Welcome to the API",
		Status:  200,
	}

	writeJSONResponse(w, response)
}

// Function to log the request details
func logRequest(r *http.Request) {
	fmt.Printf("Received request for: %s from %s\n", r.URL.Path, r.RemoteAddr)
}

// Function to check for duplicate requests
func isDuplicateRequest(r *http.Request) bool {
	clientIP := r.RemoteAddr
	requestKey := clientIP + r.URL.Path

	lastRequestTime, exists := requestMap.Load(requestKey)
	if exists && time.Since(lastRequestTime.(time.Time)) < 2*time.Second {
		fmt.Printf("Duplicate request detected: %s\n", r.URL.Path)
		return true
	}

	requestMap.Store(requestKey, time.Now())
	return false
}

// Function to write JSON response
func writeJSONResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Define routes with CORS middleware
	http.Handle("/people", corsMiddleware(http.HandlerFunc(personHandler)))
	http.Handle("/buildings", corsMiddleware(http.HandlerFunc(buildingHandler)))
	http.Handle("/world", corsMiddleware(http.HandlerFunc(worldHandler)))

	// Default handler for the root path or undefined paths
	http.Handle("/", corsMiddleware(http.HandlerFunc(defaultHandler)))

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
