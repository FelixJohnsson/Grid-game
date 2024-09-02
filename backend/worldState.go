package main

import (
	"math"
)

type WorldAccessor interface {
	GetPersonInVision(x, y, visionRange int) []PersonInVision
	GetWaterInVision(x, y, visionRange int) []Tile
	// Returns a Tile slice of grass tiles in the vision of the person at the given location, up to the given range.
	GetGrassInVision(x, y, visionRange int) []Tile
	GetPlantsInVision(x, y, visionRange int) []*Plant
	
	GetPersonByFullName(FullName string) *Person
	GetTileType(x, y int) TileType
	GetTile(x, y int) Tile
	IsAdjacent(x1, y1, x2, y2 int) bool
	CalculateDistance(x1, y1, x2, y2 int) int
	CanWalk(x, y int) bool

	MovePerson(person *Person, newX, newY int)

	AddItem(x, y int, item *Item)
	DestroyItem(item *Item)
	RemovePlant(Plant *Plant) Tile

	AddShelter(x, y int, shelter *Shelter)

	DisplayMap()
}

// NewTile creates a new tile with the given type and updates it's location.
func NewTile(t TileType, x, y int) Tile {
	return Tile{
		Type:     t,
		Location: Location{X: x, Y: y},
	}
}

// NewWorld creates a new world with the given dimensions.
func NewWorld(width, height int) *World {
	world := World{
		Tiles: make([][]Tile, height),
		Width: width,
		Height: height,
	}

	for i := range world.Tiles {
		world.Tiles[i] = make([]Tile, width)
		for j := range world.Tiles[i] {
			world.Tiles[i][j] = NewTile(Grass, j, i)
		}
	}

	return &world
}

// SetTile sets the tile at the given location to the given type.
func (w *World) SetTileType(x, y int, t TileType) {
	w.Tiles[y][x].Type = t
}

// GetTile returns the tile at the given location.
func (w *World) GetTile(x, y int) Tile {
	return w.Tiles[y][x]
}

// GetTiles returns all the tiles in the world.
func (w *World) GetTiles() [][]Tile {
	return w.Tiles
}

// CanWalk returns true if the person can walk on the tile at the given location.
func (w *World) CanWalk(x, y int) bool {
	return w.Tiles[y][x].Type != Mountain
}

// GetPersonInVision returns the vision of the person at the given location, up to the given range.
func (w *World) GetPersonInVision(x, y, visionRange int) []PersonInVision {
	var persons []PersonInVision

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
				tile := w.Tiles[ty][tx]
					cleanedPerson := PersonInVision{
						FirstName:  tile.Person.FirstName,
						FamilyName: tile.Person.FamilyName,
						Gender:     tile.Person.Gender,
						Age:        tile.Person.Age,
						Title:      tile.Person.Title,
						Location:   tile.Person.Location,
						Body:       tile.Person.Body,
					}

					persons = append(persons, cleanedPerson)
				
				
			}
		}
	}
	
	return persons
}

// GetWaterInVision returns the water in the vision of the person at the given location, up to the given range.
func (w *World) GetWaterInVision(x, y, visionRange int) []Tile {
	var water []Tile

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
				tile := w.Tiles[ty][tx]
				if tile.Type == Water {
					tileInVision := tile
					water = append(water, tileInVision)
				}
			}
		}
	}

	return water
}

// GetGrassInVision returns the grass in the vision of the person at the given location, up to the given range.
func (w *World) GetGrassInVision(x, y, visionRange int) []Tile {
	var grass []Tile

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
				tile := w.Tiles[ty][tx]
				if tile.Type == Grass {
					tileInVision := tile
					grass = append(grass, tileInVision)
				}
			}
		}
	}

	return grass
}


// GetPlantsInVision returns the plants in the vision of the person at the given location, up to the given range.
func (w *World) GetPlantsInVision(x, y, visionRange int) []*Plant {
	var plants []*Plant

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
				tile := w.Tiles[ty][tx]
				if tile.Plant != nil {
					plants = append(plants, tile.Plant)
				}
			}
		}
	}

	return plants
}


// CalculateDistance calculates the distance between two locations.
func (w *World) CalculateDistance(x1, y1, x2, y2 int) int {
	return int(math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2)))
}

// AddBuilding adds a building to the tile at the given location.
func (w *World) AddShelter(x, y int, shelter *Shelter) {
	w.Tiles[y][x].Shelter = shelter
}

// IsAdjacent returns true if the two locations are adjacent to each other.
func (w *World) IsAdjacent(x1, y1, x2, y2 int) bool {
	return (x1 == x2 && (y1 == y2+1 || y1 == y2-1)) || (y1 == y2 && (x1 == x2+1 || x1 == x2-1))
}

// AddPerson adds a person to the tile at the given location.
func (w *World) AddPerson(x, y int, person *Person) {
	w.Tiles[y][x].Person = person
}

// GetPersonByFullName returns the person with the given full name in the world.
func (w *World) GetPersonByFullName(FullName string) *Person {
	for _, row := range w.Tiles {
		for _, tile := range row {
			if tile.Person.FullName == FullName {
				return tile.Person
			}
		}
	}
	return nil
}

// Get tile type at a given location
func (w *World) GetTileType(x, y int) TileType {
	return w.Tiles[y][x].Type
}

// GetPersons returns the persons at the given location.
func (w *World) GetPersons(x, y int) *Person {
	tile := w.Tiles[y][x]

	return tile.Person
}

// GetAllPersons returns all the persons in the world.
func (w *World) GetAllPersons() []*Person {
	var persons []*Person

	for _, row := range w.Tiles {
		for _, tile := range row {
			persons = append(persons, tile.Person)
		}
	}

	return persons
}

// RemovePerson removes the person with the given full name and coordinates from the world.
func (w *World) RemovePerson(person *Person, x, y int) bool {
	tile := w.Tiles[y][x]

	// Remove the person from the tile
	if tile.Person == person {
		tile.Person = nil
		w.Tiles[y][x] = tile
		return true
	} else {
		return false
	}
}

// MovePerson moves the person with the given full name to the new location.
func (w *World) MovePerson(person *Person, newX, newY int) {
	oldX, oldY := person.Location.X, person.Location.Y

	w.RemovePerson(person, oldX, oldY)
	w.AddPerson(newX, newY, person)
	
	person.UpdateLocation(newX, newY)
}

// AddItem adds an item to the tile at the given location.
func (w *World) AddItem(x, y int, item *Item) {
	item.Location.X = x
	item.Location.Y = y
	w.Tiles[y][x].Items = append(w.Tiles[y][x].Items, item)
}

// DestroyItem removes the memory allocation of the pointer to the item
func (w *World) DestroyItem(item *Item) {
	w.RemoveItem(item, item.Location.X, item.Location.Y)
	item = nil
}

// GetItems returns the items at the given location.
func (w *World) GetItems(x, y int) []*Item {
	tile := w.Tiles[y][x]

	return tile.Items
}

// RemoveItem removes the item from the tile at the given location.
func (w *World) RemoveItem(Item *Item, x, y int) []*Item {
	tile := w.Tiles[y][x]

	// Find the item in the tile and remove it
	everything := tile.Items
	for i, item := range everything {
		if item == Item {
			everything = append(everything[:i], everything[i+1:]...)
			break
		}
	}

	// Update the tile with the new list of items
	tile.Items = everything

	// Update the world with the updated tile
	w.Tiles[y][x] = tile

	return everything
}

// AddPlant adds a plant to the tile at the given location.
func (w *World) AddPlant(x, y int, plant *Plant) {
	w.Tiles[y][x].Plant = plant
}

// GetPlants returns the plants at the given location.
func (w *World) GetPlants(x, y int) *Plant {
	tile := w.Tiles[y][x]

	return tile.Plant
}

// RemovePlant removes the plant from the tile at the given location.
func (w *World) RemovePlant(Plant *Plant) Tile {
	w.Tiles[Plant.Location.Y][Plant.Location.X].Plant = nil

	return w.Tiles[Plant.Location.Y][Plant.Location.X]
}
