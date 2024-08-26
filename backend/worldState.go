package main

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
type Vision struct {
	Buildings []BuildingCleaned `json:"buildings"`
	Persons   []PersonCleaned   `json:"persons"`
}
type PersonCleaned struct {
	FullName     string       `json:"FullName"`
	Location     Location     `json:"Location"`
}
type BuildingCleaned struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location Location `json:"location"`
}
type WorldAccessor interface {
    GetVision(x, y, visionRange int) Vision
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

func (w *World) GetVision(x, y, visionRange int) Vision {
    var buildings []BuildingCleaned
    var persons []PersonCleaned

    for i := -visionRange; i <= visionRange; i++ {
        for j := -visionRange; j <= visionRange; j++ {
            tx, ty := x+i, y+j

            if tx >= 0 && tx < len(w.Tiles[0]) && ty >= 0 && ty < len(w.Tiles) {
                tile := w.Tiles[ty][tx]

                if tile.Building != nil {
                    cleanedBuilding := BuildingCleaned{
                        Name:     tile.Building.Name,
                        Type:     string(tile.Building.Type),
                        Location: tile.Building.Location,
                    }
                    buildings = append(buildings, cleanedBuilding)
                }

                for _, person := range tile.Persons {
                    cleanedPerson := PersonCleaned{
                        FullName: person.FullName,
                        Location: person.Location,
                    }
                    persons = append(persons, cleanedPerson)
                }
            }
        }
    }

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

	// Add the person to the new location
	w.Tiles[newY][newX].Persons = append(w.Tiles[newY][newX].Persons, person)

	// Update the person's location
	person.Location.X = newX
	person.Location.Y = newY
}

// RequestTask requests a task from the brain for the person by name at the given location.
func (w *World) RequestTaskFrom(fullName string, x, y int) bool {
	tile := w.Tiles[y][x]
	for _, person := range tile.Persons {
		if person.FullName == fullName {
			requestedTask := RequestedAction{action: "Talk", fromPerson: fullName}
			return person.Brain.RequestTask(requestedTask)
		}
	}
	return false
}
