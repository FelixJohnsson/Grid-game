package main

import "math"

type WorldAccessor interface {
	GetPersonInVision(x, y, visionRange int) []PersonInVision
	GetWaterInVision(x, y, visionRange int) []TileInVision
	GetPlantsInVision(x, y, visionRange int) []*Plant
	
	GetPersonByFullName(FullName string) *Person
	GetTileType(x, y int) TileType
	IsAdjacent(x1, y1, x2, y2 int) bool
	CalculateDistance(x1, y1, x2, y2 int) int
	CanWalk(x, y int) bool
}

// NewWorld creates a new world with the given dimensions.
func NewWorld(width, height int) *World {
	world := World{
		Tiles: make([][]Tile, height),
	}

	for i := range world.Tiles {
		world.Tiles[i] = make([]Tile, width)
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
	return w.Tiles[y][x].Type != Water && w.Tiles[y][x].Type != Mountain && w.Tiles[y][x].Building == nil
}

// GetPersonInVision returns the vision of the person at the given location, up to the given range.
func (w *World) GetPersonInVision(x, y, visionRange int) []PersonInVision {
	var persons []PersonInVision

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
				tile := w.Tiles[ty][tx]
				for _, person := range tile.Persons {
					cleanedPerson := PersonInVision{
						FirstName:  person.FirstName,
						FamilyName: person.FamilyName,
						Gender:     person.Gender,
						Age:        person.Age,
						Title:      person.Title,
						Location:   person.Location,
						Body:       person.Body,
					}
					persons = append(persons, cleanedPerson)
				}
				
			}
		}
	}
	
	return persons
}

// GetWaterInVision returns the water in the vision of the person at the given location, up to the given range.
func (w *World) GetWaterInVision(x, y, visionRange int) []TileInVision {
	var water []TileInVision

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
				tile := w.Tiles[ty][tx]
				if tile.Type == Water {
					tileInVision := TileInVision{
						Tile:     tile,
						Location: Location{X: tx, Y: ty},
					}
					water = append(water, tileInVision)
				}
			}
		}
	}

	return water
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
func (w *World) AddBuilding(x, y int, b Building) {
	w.Tiles[y][x].Building = &b
}

// IsAdjacent returns true if the two locations are adjacent to each other.
func (w *World) IsAdjacent(x1, y1, x2, y2 int) bool {
	return (x1 == x2 && (y1 == y2+1 || y1 == y2-1)) || (y1 == y2 && (x1 == x2+1 || x1 == x2-1))
}

// RemoveBuilding removes the building from the tile at the given location.
func (w *World) RemoveBuilding(x, y int) {
	w.Tiles[y][x].Building = nil
}

// GetBuilding returns the building at the given location.
func (w *World) GetBuilding(x, y int) *Building {
	return w.Tiles[y][x].Building
}

// GetAllBuildings returns all the buildings in the world.
func (w *World) GetAllBuildings() []Building {
	var buildings []Building

	for _, row := range w.Tiles {
		for _, tile := range row {
			if tile.Building != nil {
				buildings = append(buildings, *tile.Building)
			}
		}
	}

	return buildings
}

// AddPerson adds a person to the tile at the given location.
func (w *World) AddPerson(x, y int, p *Person) {
	w.Tiles[y][x].Persons = append(w.Tiles[y][x].Persons, p)
}

// GetPersonByFullName returns the person with the given full name in the world.
func (w *World) GetPersonByFullName(FullName string) *Person {
	for _, row := range w.Tiles {
		for _, tile := range row {
			for _, person := range tile.Persons {
				if person.FullName == FullName {
					return person
				}
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
func (w *World) GetPersons(x, y int) []*Person {
	tile := w.Tiles[y][x]

	return tile.Persons
}

// GetAllPersons returns all the persons in the world.
func (w *World) GetAllPersons() []*Person {
	var persons []*Person

	for _, row := range w.Tiles {
		for _, tile := range row {
			persons = append(persons, tile.Persons...)
		}
	}

	return persons
}

// RemovePerson removes the person with the given full name and coordinates from the world.
func (w *World) RemovePerson(FullName string, x, y int) []*Person {
	tile := w.Tiles[y][x]

	// Find the person in the tile and remove it
	everyone := tile.Persons
	for i, p := range everyone {
		if p.FullName == FullName {
			everyone = append(everyone[:i], everyone[i+1:]...)
			break
		}
	}

	// Update the tile with the new list of people
	tile.Persons = everyone

	// Update the world with the updated tile
	w.Tiles[y][x] = tile

	return everyone
}

// MovePerson moves the person with the given full name to the new location.
func (w *World) MovePerson(FullName string, newX, newY int) {
	// Find the person in the world
	var person *Person
	var oldX, oldY int
	for y, row := range w.Tiles {
		for x, tile := range row {
			for _, p := range tile.Persons {
				if p.FullName == FullName {
					person = p
					oldX, oldY = x, y
					break
				}
			}
		}
	}

	// Remove the person from the old location
	w.Tiles[oldY][oldX].Persons = w.RemovePerson(FullName, oldX, oldY)

	// Update the person's location
	person.UpdateLocation(newX, newY)

	// Add the person to the new location
	w.Tiles[newY][newX].Persons = append(w.Tiles[newY][newX].Persons, person)

}

// AddItem adds an item to the tile at the given location.
func (w *World) AddItem(x, y int, item *Item) {
	w.Tiles[y][x].Items = append(w.Tiles[y][x].Items, item)
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
func (w *World) RemovePlant(Plant *Plant, x, y int) Tile {
	w.Tiles[y][x].Plant = nil

	return w.Tiles[y][x]
}
