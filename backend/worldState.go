package main

import (
	"fmt"
)

// TileType represents different types of terrain.
type TileType int

// Tile represents a single tile in the world.
type Tile struct {
	Type     TileType  `json:"type"`
	Building *Building `json:"building,omitempty"`
	Persons  []*Person   `json:"persons,omitempty"`
}

// Constants representing different types of terrain.
const (
	Grass TileType = iota
	Water
	Mountain
)

// World represents a 2D array of tiles.
type World struct {
	Tiles [][]Tile `json:"tiles"`
}

// NewWorld creates a new world with the given dimensions.
func NewWorld(width, height int) World {
	world := World{
		Tiles: make([][]Tile, height),
	}

	for i := range world.Tiles {
		world.Tiles[i] = make([]Tile, width)
	}

	return world
}

// SetTile sets the tile at the given location to the given type.
func (w *World) SetTile(x, y int, t TileType) {
	w.Tiles[y][x].Type = t
}

// GetTile returns the tile at the given location.
func (w *World) GetTile(x, y int) Tile {
	return w.Tiles[y][x]
}
type Vision struct {
	Buildings []BuildingCleaned       `json:"buildings"`
	Persons   []PersonCleaned  `json:"persons"`
}

type PersonCleaned struct {
	FullName     string   `json:"name"`
	Location Location `json:"location"`
}

type BuildingCleaned struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location Location `json:"location"`
}

// GetVision returns the tiles in the vision range of the person at the given location.
func (w *World) GetVision(x, y, visionRange int) Vision {
	var buildings []BuildingCleaned
	var persons []PersonCleaned

	// Loop through the vision range
	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			// Calculate the coordinates in the world
			tx, ty := x+i, y+j

			// Ensure the coordinates are within the world boundaries
			if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
				tile := w.Tiles[ty][tx]

				// Add building if it exists
				if tile.Building != nil {
					cleanedBuilding := BuildingCleaned{
						Name:     tile.Building.Name,
						Type:     string(tile.Building.Type),
						Location: tile.Building.Location,
					}
					buildings = append(buildings, cleanedBuilding)
				}

				// Add people if they exist
				for _, person := range tile.Persons {
					cleanedPerson := PersonCleaned{
						FullName:     person.FullName,
						Location: person.Location,
					}
					persons = append(persons, cleanedPerson)
				}
			}
		}
	}

	// Prepare the vision result
	vision := Vision{
		Buildings: buildings,
		Persons:   persons,
	}

	return vision
}

// AddBuilding adds a building to the tile at the given location.
func (w *World) AddBuilding(x, y int, b Building) {
	w.Tiles[y][x].Building = &b
}

// RemoveBuilding removes the building from the tile at the given location.
func (w *World) RemoveBuilding(x, y int) {
	w.Tiles[y][x].Building = nil
}

// GetBuilding returns the building at the given location.
func (w *World) GetBuilding(x, y int) *Building {
	return w.Tiles[y][x].Building
}

// UpdateState updates the state of all buildings in the world.
func (w *World) UpdateState() {
	for y, row := range w.Tiles {
		for x := range row {
			if w.Tiles[y][x].Building != nil {
				w.Tiles[y][x].Building.UpdateState()
			}
		}
	}
}

// getWorld returns the current world state, creating a new one if necessary.
func getWorld() World {
	world, err := loadWorldFromFile()
	if err != nil {
		fmt.Println("Creating a new world due to:", err)
		return createNewWorld(10, 10)
	}
	return world
}

// createNewWorld creates a new world with the specified dimensions and populates it.
func createNewWorld(width, height int) World {
	world := NewWorld(width, height)
	populateWorld(&world)
	saveWorldToFile(world)
	return world
}

// populateWorld populates the world with buildings and people.
func populateWorld(world *World) {
	// Add buildings to the world
	buildings := getBuildings()
	for _, b := range buildings {
		world.AddBuilding(b.Location.X, b.Location.Y, b)
	}

	// Add people to the world
	persons := getPersons()
	for _, p := range persons {
		tile := &world.Tiles[p.Location.Y][p.Location.X]
		tile.Persons = append(tile.Persons, &p) // Append the person to the Persons slice
	}

	// Save the world to a file after populating it
	saveWorldToFile(*world)
}
