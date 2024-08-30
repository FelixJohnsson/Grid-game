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
	Message []*Person `json:"message"`
	Status  int      `json:"status"`
}
type BuildingResponse struct {
	Message []Building `json:"message"`
	Status  int        `json:"status"`
}
type MoveRequest struct {
	FullName  string `json:"full_name"`
	Direction string `json:"direction"`
}
type GrabRequest struct {
	ItemName string `json:"ItemName"`
	FullName string `json:"FullName"`
}
type AttackRequest struct {
	FullName       string `json:"FullName"`
	TargetFullName string `json:"TargetFullName"`
}

type WorldResponse struct {
    Message [][]CleanedTile `json:"message"`
    Status  int             `json:"status"`
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
func (w *World) personHandler(writer http.ResponseWriter, r *http.Request) {
	// Ignore favicon requests
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(writer, r)
		return
	}

	logRequest(r)

	response := PersonResponse{
		Message: w.GetAllPersons(),
		Status:  200,
	}

	writeJSONResponse(writer, response)
}

// Hander for /buildings endpoint
func (w *World) buildingHandler(writer http.ResponseWriter, r *http.Request) {
	// Ignore favicon requests
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(writer, r)
		return
	}

	logRequest(r)

	buildings := w.GetAllBuildings()

	response := BuildingResponse{
		Message: buildings,
		Status:  200,
	}

	writeJSONResponse(writer, response)
}

// Handler for /world endpoint
func (w *World) worldHandler(writer http.ResponseWriter, r *http.Request) {
    // Ignore favicon requests
    if r.URL.Path == "/favicon.ico" {
        http.NotFound(writer, r)
        return
    }

    logRequest(r)

    response := WorldResponse{
        Message: w.CleanTiles(),
        Status:  200,
    }

    writeJSONResponse(writer, response)
}

// Handler for /move endpoint
func (w *World) moveHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var moveRequest MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&moveRequest); err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// We need to calculate the new coordinates based on the direction
	startingCoordinates := w.GetPersonByFullName(moveRequest.FullName).Location
	switch moveRequest.Direction {
	case "up":
		startingCoordinates.Y--
	case "down":
		startingCoordinates.Y++
	case "left":
		startingCoordinates.X--
	case "right":
		startingCoordinates.X++
	}

	// Move the person in the world
	w.MovePerson(moveRequest.FullName, startingCoordinates.X, startingCoordinates.Y)

	response := WorldResponse{
		Message: w.CleanTiles(),
		Status:  200,
	}
	writeJSONResponse(writer, response)
}

func (w *World) grabHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var grabRequest GrabRequest
	if err := json.NewDecoder(r.Body).Decode(&grabRequest); err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(grabRequest)

	// Find the person by full name
	person := w.GetPersonByFullName(grabRequest.FullName)
	if person == nil {
		http.Error(writer, "Person not found", http.StatusBadRequest)
		return
	}

	// Get the person's current location
	coordinates := person.Location
	tile := w.Tiles[coordinates.Y][coordinates.X]

	if tile.Items == nil {
		http.Error(writer, "No items to grab", http.StatusBadRequest)
		return
	} else {
		for _, item := range tile.Items {
			if item.Name == grabRequest.ItemName {
				// The item is found
				person.GrabRight(item)

				// Remove the item from the tile
				for i, tileItem := range tile.Items {
					if tileItem.Name == grabRequest.ItemName {
						tile.Items = append(tile.Items[:i], tile.Items[i+1:]...)

						// Update the tile with the new list of items
						w.Tiles[coordinates.Y][coordinates.X].Items = tile.Items

						break // Exit the loop after finding the item
					}
				}

				break // Exit the loop after processing the item
			}
		}
	}
	
	response := WorldResponse{
		Message: w.CleanTiles(),
		Status:  200,
	}
	writeJSONResponse(writer, response)
}

// Handler for /entityAttack endpoint
func (w *World) attackHandler(writer http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var attackRequest AttackRequest
	if err := json.NewDecoder(r.Body).Decode(&attackRequest); err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the person by full name
	person := w.GetPersonByFullName(attackRequest.FullName)
	if person == nil {
		http.Error(writer, "Person not found", http.StatusBadRequest)
		return
	}

	// Check if the target person is one tile away
	targetPerson := w.GetPersonByFullName(attackRequest.TargetFullName)
	if targetPerson == nil {
		http.Error(writer, "Target person not found", http.StatusBadRequest)
		return
	} else if !person.WorldProvider.IsAdjacent(targetPerson.Location.X, targetPerson.Location.Y, person.Location.X, person.Location.Y) {
		http.Error(writer, "Target person is not adjacent", http.StatusBadRequest)
		return
	} else {
		// Attack the target person
		person.AttackWith(targetPerson, "Head", person.Body.RightArm.Hand)
	}

	response := WorldResponse{
		Message: w.CleanTiles(),
		Status:  200,
	}
	writeJSONResponse(writer, response)
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

// Function to write JSON response with error handling
func writeJSONResponse(w http.ResponseWriter, response interface{}) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Attempt to encode the response into JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error
		fmt.Printf("Error encoding response to JSON: %v\n", err)

		// If encoding fails, respond with an internal server error status
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Initialize the world
	world := initializeWorld()

	// Define routes with CORS middleware and pass the world instance
	http.Handle("/world", corsMiddleware(http.HandlerFunc(world.worldHandler)))

	http.Handle("/people", corsMiddleware(http.HandlerFunc(world.personHandler)))
	http.Handle("/buildings", corsMiddleware(http.HandlerFunc(world.buildingHandler)))
	http.Handle("/move", corsMiddleware(http.HandlerFunc(world.moveHandler)))
	http.Handle("/entityGrab", corsMiddleware(http.HandlerFunc(world.grabHandler)))
	http.Handle("/entityAttack", corsMiddleware(http.HandlerFunc(world.attackHandler)))

	// Default handler for the root path or undefined paths
	http.Handle("/", corsMiddleware(http.HandlerFunc(defaultHandler)))

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}