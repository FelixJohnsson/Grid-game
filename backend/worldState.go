package main

import (
	"errors"
	"fmt"
	"math"
)

var SIZE_OF_MAP = 100

var (
	ErrOutOfBounds          = errors.New("location out of bounds")
	ErrNilEntity            = errors.New("entity is nil")
	ErrNilItem              = errors.New("item is nil")
	ErrNilPlant             = errors.New("plant is nil")
	ErrNilShelter           = errors.New("shelter is nil")
	ErrTileOccupied         = errors.New("tile is occupied")
	ErrUnwalkableTile       = errors.New("tile is not walkable")
	ErrEntityNotFoundOnTile = errors.New("entity not found on tile")
	ErrItemNotFoundOnTile   = errors.New("item not found on tile")
	ErrPlantNotFoundOnTile  = errors.New("plant not found on tile")
	ErrNoRipeFruitOnPlant   = errors.New("plant has no ripe fruit")
)

type WorldAccessor interface {
	GetVision(x, y, visionRange int) []Tile
	GetEntityInVision(x, y, visionRange int) []EntityInVision
	GetWaterInVision(x, y, visionRange int) []Tile
	GetGrassInVision(x, y, visionRange int) []Tile
	GetPlantsInVision(x, y, visionRange int) []*Plant
	GetFruitingPlantsInVision(x, y, visionRange int) []*Plant

	GetPersonByFullName(FullName string) *Entity
	GetTileType(x, y int) TileType
	GetTile(x, y int) Tile
	IsTileEmpty(x, y int) bool
	IsAdjacent(x1, y1, x2, y2 int) bool
	CalculateDistance(Location1, Location2 Location) int
	CanWalk(x, y int) bool

	MoveEntity(entity *Entity, newX, newY int) error

	AddItem(x, y int, item *Item) error
	DestroyItem(item *Item) error
	RemovePlant(Plant *Plant) (Tile, error)
	ConsumeRipeFruitAt(x, y int) (Fruit, error)

	AddShelter(x, y int, shelter *Shelter) error
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
		Tiles:  make([][]Tile, height),
		Width:  width,
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

func (w *World) InBounds(x, y int) bool {
	return x >= 0 && x < w.Width && y >= 0 && y < w.Height
}

func (w *World) validateBounds(x, y int) error {
	if !w.InBounds(x, y) {
		return fmt.Errorf("%w: (%d,%d)", ErrOutOfBounds, x, y)
	}
	return nil
}

// SetTile sets the tile at the given location to the given type.
func (w *World) SetTileType(x, y int, t TileType) {
	if err := w.validateBounds(x, y); err != nil {
		return
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Tiles[y][x].Type = t
}

// GetTile returns the tile at the given location.
func (w *World) GetTile(x, y int) Tile {
	if !w.InBounds(x, y) {
		// Return an unwalkable tile for out-of-bounds reads.
		return Tile{
			Type:     Mountain,
			Location: Location{X: x, Y: y},
		}
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Tiles[y][x]
}

// GetTiles returns all the tiles in the world.
func (w *World) GetTiles() [][]Tile {
	w.mu.RLock()
	defer w.mu.RUnlock()

	tileCopy := make([][]Tile, len(w.Tiles))
	for y, row := range w.Tiles {
		tileCopy[y] = make([]Tile, len(row))
		copy(tileCopy[y], row)
	}

	return tileCopy
}

// IsTileWater - Check if the person is standing on water
func (w *World) IsTileWater(x, y int) bool {
	if err := w.validateBounds(x, y); err != nil {
		return false
	}

	tile := w.GetTile(x, y)
	return tile.Type == Water
}

// IsTileEmpty - Check if a tile is empty
func (w *World) IsTileEmpty(x, y int) bool {
	if err := w.validateBounds(x, y); err != nil {
		return false
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	tile := w.Tiles[y][x]
	return tile.Shelter == nil && tile.Plant == nil && tile.Entity == nil
}

// CanWalk returns true if the person can walk on the tile at the given location.
func (w *World) CanWalk(x, y int) bool {
	if err := w.validateBounds(x, y); err != nil {
		return false
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.Tiles[y][x].Type != Mountain
}

// GetEntityInVision returns the vision of the person at the given location, up to the given range.
func (w *World) GetEntityInVision(x, y, visionRange int) []EntityInVision {
	var persons []EntityInVision
	w.mu.RLock()
	defer w.mu.RUnlock()

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if w.InBounds(tx, ty) {
				tile := w.Tiles[ty][tx]
				if tile.Entity == nil {
					continue
				}

				cleanedPerson := EntityInVision{
					FirstName:  tile.Entity.FirstName,
					FamilyName: tile.Entity.FamilyName,
					Gender:     tile.Entity.Gender,
					Age:        tile.Entity.Age,
					Title:      tile.Entity.Title,
					Location:   tile.Entity.Location,
					Body:       tile.Entity.Body,
				}

				persons = append(persons, cleanedPerson)
			}
		}
	}

	return persons
}

func (w *World) GetVision(x, y, visionRange int) []Tile {
	var vision []Tile
	w.mu.RLock()
	defer w.mu.RUnlock()

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if w.InBounds(tx, ty) {
				tile := w.Tiles[ty][tx]
				vision = append(vision, tile)
			}
		}
	}

	return vision
}

// GetWaterInVision returns the water in the vision of the person at the given location, up to the given range.
func (w *World) GetWaterInVision(x, y, visionRange int) []Tile {
	var water []Tile
	w.mu.RLock()
	defer w.mu.RUnlock()

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if w.InBounds(tx, ty) {
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
	w.mu.RLock()
	defer w.mu.RUnlock()

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j

			if w.InBounds(tx, ty) {
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
	w.mu.RLock()
	defer w.mu.RUnlock()

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j
			if w.InBounds(tx, ty) {
				tile := w.Tiles[ty][tx]
				if tile.Plant != nil {
					plants = append(plants, tile.Plant)
				}
			}
		}
	}

	return plants
}

func (w *World) GetFruitingPlantsInVision(x, y, visionRange int) []*Plant {
	var plants []*Plant
	w.mu.RLock()
	defer w.mu.RUnlock()

	for i := -visionRange; i <= visionRange; i++ {
		for j := -visionRange; j <= visionRange; j++ {
			tx, ty := x+i, y+j
			if w.InBounds(tx, ty) {
				tile := w.Tiles[ty][tx]
				if tile.Plant != nil && tile.Plant.ProducesFruit && HasRipeFruit(tile.Plant) {
					plants = append(plants, tile.Plant)
				}
			}
		}
	}
	return plants
}

// CalculateDistance calculates the distance between two locations.
func (w *World) CalculateDistance(Location1, Location2 Location) int {
	return int(math.Abs(float64(Location1.X-Location2.X)) + math.Abs(float64(Location1.Y-Location2.Y)))
}

// AddBuilding adds a building to the tile at the given location.
func (w *World) AddShelter(x, y int, shelter *Shelter) error {
	if w.intentEngineIsEnabled() {
		result := w.submitWorldIntent(worldIntentRequest{
			intentType: worldIntentAddShelter,
			shelter:    shelter,
			x:          x,
			y:          y,
		})
		return result.err
	}
	return w.addShelterDirect(x, y, shelter)
}

func (w *World) addShelterDirect(x, y int, shelter *Shelter) error {
	if shelter == nil {
		return ErrNilShelter
	}
	if err := w.validateBounds(x, y); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	w.Tiles[y][x].Shelter = shelter
	return nil
}

// IsAdjacent returns true if the two locations are adjacent to each other.
func (w *World) IsAdjacent(x1, y1, x2, y2 int) bool {
	return (x1 == x2 && (y1 == y2+1 || y1 == y2-1)) || (y1 == y2 && (x1 == x2+1 || x1 == x2-1))
}

// AddEntity adds a person to the tile at the given location.
func (w *World) AddEntity(x, y int, entity *Entity) error {
	if entity == nil {
		return ErrNilEntity
	}
	if err := w.validateBounds(x, y); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Tiles[y][x].Entity != nil && w.Tiles[y][x].Entity != entity {
		return fmt.Errorf("%w: (%d,%d)", ErrTileOccupied, x, y)
	}

	w.Tiles[y][x].Entity = entity
	entity.UpdateLocation(x, y)
	return nil
}

// GetPersonByFullName returns the person with the given full name in the world.
func (w *World) GetPersonByFullName(FullName string) *Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, row := range w.Tiles {
		for _, tile := range row {
			if tile.Entity != nil && tile.Entity.FullName == FullName {
				return tile.Entity
			}
		}
	}
	return nil
}

// Get tile type at a given location
func (w *World) GetTileType(x, y int) TileType {
	if err := w.validateBounds(x, y); err != nil {
		return Mountain
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.Tiles[y][x].Type
}

// GetPersons returns the persons at the given location.
func (w *World) GetPersons(x, y int) *Entity {
	if err := w.validateBounds(x, y); err != nil {
		return nil
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	tile := w.Tiles[y][x]

	return tile.Entity
}

// GetAllPersons returns all the persons in the world.
func (w *World) GetAllPersons() []*Entity {
	var persons []*Entity
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, row := range w.Tiles {
		for _, tile := range row {
			if tile.Entity != nil {
				persons = append(persons, tile.Entity)
			}
		}
	}

	return persons
}

// RemoveEntity removes the person with the given full name and coordinates from the world.
func (w *World) RemoveEntity(entity *Entity, x, y int) error {
	if entity == nil {
		return ErrNilEntity
	}
	if err := w.validateBounds(x, y); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// Remove the person from the tile
	if w.Tiles[y][x].Entity != entity {
		return fmt.Errorf("%w: (%d,%d)", ErrEntityNotFoundOnTile, x, y)
	}
	w.Tiles[y][x].Entity = nil
	return nil
}

// MoveEntity moves the person with the given full name to the new location.
func (w *World) MoveEntity(entity *Entity, newX, newY int) error {
	if w.intentEngineIsEnabled() {
		result := w.submitWorldIntent(worldIntentRequest{
			intentType: worldIntentMoveEntity,
			entity:     entity,
			x:          newX,
			y:          newY,
		})
		return result.err
	}
	return w.moveEntityDirect(entity, newX, newY)
}

func (w *World) moveEntityDirect(entity *Entity, newX, newY int) error {
	if entity == nil {
		return ErrNilEntity
	}
	if err := w.validateBounds(newX, newY); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	oldX, oldY := entity.Location.X, entity.Location.Y
	if !w.InBounds(oldX, oldY) {
		return fmt.Errorf("%w: (%d,%d)", ErrOutOfBounds, oldX, oldY)
	}

	if w.Tiles[oldY][oldX].Entity != entity {
		return fmt.Errorf("%w: (%d,%d)", ErrEntityNotFoundOnTile, oldX, oldY)
	}
	if w.Tiles[newY][newX].Type == Mountain {
		return fmt.Errorf("%w: (%d,%d)", ErrUnwalkableTile, newX, newY)
	}
	if w.Tiles[newY][newX].Entity != nil && w.Tiles[newY][newX].Entity != entity {
		return fmt.Errorf("%w: (%d,%d)", ErrTileOccupied, newX, newY)
	}

	w.Tiles[oldY][oldX].Entity = nil
	w.Tiles[newY][newX].Entity = entity
	entity.UpdateLocation(newX, newY)

	return nil
}

// AddItem adds an item to the tile at the given location.
func (w *World) AddItem(x, y int, item *Item) error {
	if w.intentEngineIsEnabled() {
		result := w.submitWorldIntent(worldIntentRequest{
			intentType: worldIntentAddItem,
			item:       item,
			x:          x,
			y:          y,
		})
		return result.err
	}
	return w.addItemDirect(x, y, item)
}

func (w *World) addItemDirect(x, y int, item *Item) error {
	if item == nil {
		return ErrNilItem
	}
	if err := w.validateBounds(x, y); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	item.Location.X = x
	item.Location.Y = y
	w.Tiles[y][x].Items = append(w.Tiles[y][x].Items, item)
	return nil
}

// DestroyItem removes the memory allocation of the pointer to the item
func (w *World) DestroyItem(item *Item) error {
	if w.intentEngineIsEnabled() {
		result := w.submitWorldIntent(worldIntentRequest{
			intentType: worldIntentDestroyItem,
			item:       item,
		})
		return result.err
	}
	return w.destroyItemDirect(item)
}

func (w *World) destroyItemDirect(item *Item) error {
	if item == nil {
		return ErrNilItem
	}
	_, err := w.RemoveItem(item, item.Location.X, item.Location.Y)
	return err
}

// GetItems returns the items at the given location.
func (w *World) GetItems(x, y int) []*Item {
	if err := w.validateBounds(x, y); err != nil {
		return nil
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	tile := w.Tiles[y][x]
	items := make([]*Item, len(tile.Items))
	copy(items, tile.Items)

	return items
}

// RemoveItem removes the item from the tile at the given location.
func (w *World) RemoveItem(Item *Item, x, y int) ([]*Item, error) {
	if Item == nil {
		return nil, ErrNilItem
	}
	if err := w.validateBounds(x, y); err != nil {
		return nil, err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	tile := w.Tiles[y][x]

	// Find the item in the tile and remove it
	everything := tile.Items
	found := false
	for i, item := range everything {
		if item == Item {
			everything = append(everything[:i], everything[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return everything, fmt.Errorf("%w: (%d,%d)", ErrItemNotFoundOnTile, x, y)
	}

	// Update the tile with the new list of items
	tile.Items = everything

	// Update the world with the updated tile
	w.Tiles[y][x] = tile

	return everything, nil
}

// AddPlant adds a plant to the tile at the given location.
func (w *World) AddPlant(x, y int, plant *Plant) error {
	if plant == nil {
		return ErrNilPlant
	}
	if err := w.validateBounds(x, y); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	plant.Location = Location{X: x, Y: y}
	w.Tiles[y][x].Plant = plant
	return nil
}

// GetPlants returns the plants at the given location.
func (w *World) GetPlants(x, y int) *Plant {
	if err := w.validateBounds(x, y); err != nil {
		return nil
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	tile := w.Tiles[y][x]

	return tile.Plant
}

// RemovePlant removes the plant from the tile at the given location.
func (w *World) RemovePlant(Plant *Plant) (Tile, error) {
	if Plant == nil {
		return Tile{}, ErrNilPlant
	}
	if err := w.validateBounds(Plant.Location.X, Plant.Location.Y); err != nil {
		return Tile{}, err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	tile := w.Tiles[Plant.Location.Y][Plant.Location.X]
	if tile.Plant != Plant {
		return tile, fmt.Errorf("%w: (%d,%d)", ErrPlantNotFoundOnTile, Plant.Location.X, Plant.Location.Y)
	}
	tile.Plant = nil
	w.Tiles[Plant.Location.Y][Plant.Location.X] = tile

	return tile, nil
}

// ConsumeRipeFruitAt removes one ripe fruit from a plant at a location.
func (w *World) ConsumeRipeFruitAt(x, y int) (Fruit, error) {
	if w.intentEngineIsEnabled() {
		result := w.submitWorldIntent(worldIntentRequest{
			intentType: worldIntentConsumeRipeFruit,
			x:          x,
			y:          y,
		})
		return result.fruit, result.err
	}
	return w.consumeRipeFruitAtDirect(x, y)
}

func (w *World) consumeRipeFruitAtDirect(x, y int) (Fruit, error) {
	if err := w.validateBounds(x, y); err != nil {
		return Fruit{}, err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	tile := &w.Tiles[y][x]
	if tile.Plant == nil {
		return Fruit{}, fmt.Errorf("%w: (%d,%d)", ErrPlantNotFoundOnTile, x, y)
	}

	for i, fruit := range tile.Plant.Fruit {
		if fruit.IsRipe {
			tile.Plant.Fruit = append(tile.Plant.Fruit[:i], tile.Plant.Fruit[i+1:]...)
			return fruit, nil
		}
	}

	return Fruit{}, fmt.Errorf("%w: (%d,%d)", ErrNoRipeFruitOnPlant, x, y)
}
