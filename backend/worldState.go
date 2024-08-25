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
