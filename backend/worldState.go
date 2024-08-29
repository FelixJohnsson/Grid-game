package main

// TileType represents different types of terrain.
type TileType int

// Constants representing different types of terrain.
const (
	Grass TileType = iota
	Water
	Mountain
)

// Tile represents a single tile in the world.
type Tile struct {
	Type     TileType  `json:"Type"`
	Building *Building `json:"Building,omitempty"`
	Persons  []*Person   `json:"Persons,omitempty"`
	Items    []*Item    `json:"Items,omitempty"`
	Plants   []*Plant    `json:"Plant,omitempty"`
	NutritionalValue int `json:"NutritionalValue,omitempty"`
}

// World represents a 2D array of tiles.
type World struct {
	Tiles [][]Tile `json:"tiles"`
}
type Vision struct {
	Buildings []BuildingCleaned `json:"buildings"`
	Persons   []PersonCleaned   `json:"persons"`
}
type PersonCleaned struct {
	FullName      string       `json:"FullName"`
	Age 		  int          `json:"Age"`
	Title 		  string       `json:"Title"`
	Location      Location     `json:"Location"`
	IsTalking     bool         `json:"IsTalking"`
	Thinking      string       `json:"Thinking"`
	RightHand     []*Item      `json:"RightHand,omitempty"`
	LeftHand      []*Item      `json:"LeftHand,omitempty"`
	Relationships []Relationship `json:"Relationships"`
}

type PlantCleaned struct {
	Name      string `json:"Name"`
	Age       int    `json:"Age"`
	Health    int    `json:"Health"`
	IsAlive   bool   `json:"IsAlive"`
	ProducesFruit bool   `json:"ProducesFruit"`
	Fruit    []Fruit `json:"Fruit"`
	PlantStage PlantStage    `json:"PlantStage"`
}
type BuildingCleaned struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location Location `json:"location"`
}
type WorldAccessor interface {
    GetVision(x, y, visionRange int) Vision;
	GetPersonByFullName(FullName string) *Person;
	GetTileType(x, y int) TileType;
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
						Age: person.Age,
						Title: person.Title,
                        Location: person.Location,
						IsTalking: person.IsTalking.IsActive,
						Thinking: person.Thinking,
						RightHand: person.RightHand.Items,
						LeftHand: person.LeftHand.Items,
						Relationships: person.Relationships,
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
	w.Tiles[y][x].Plants = append(w.Tiles[y][x].Plants, plant)
}

// GetPlants returns the plants at the given location.
func (w *World) GetPlants(x, y int) []*Plant {
	tile := w.Tiles[y][x]

	return tile.Plants
}

// RemovePlant removes the plant from the tile at the given location.
func (w *World) RemovePlant(Plant *Plant, x, y int) []*Plant {
	tile := w.Tiles[y][x]

	// Find the plant in the tile and remove it
	everything := tile.Plants
	for i, plant := range everything {
		if plant == Plant {
			everything = append(everything[:i], everything[i+1:]...)
			break
		}
	}
	
	// Update the tile with the new list of plants
	tile.Plants = everything

	// Update the world with the updated tile
	w.Tiles[y][x] = tile

	return everything
}
