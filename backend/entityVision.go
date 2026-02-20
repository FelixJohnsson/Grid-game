package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const (
	cognitiveMapWaterMaxAge  = 8 * time.Minute
	cognitiveMapFoodMaxAge   = 2 * time.Minute
	cognitiveMapLumberMaxAge = 10 * time.Minute
	cognitiveMapGrassMaxAge  = 10 * time.Minute
)

// ----------------- Water -----------------

func (b *Brain) GetWaterInVision() []Tile {
	vision := b.Owner.WorldProvider.GetWaterInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)

	water := make([]Tile, 0)

	for _, tile := range vision {
		if tile.Type == 1 && !b.isLocationBlockedByOtherEntity(tile.Location) {
			water = append(water, tile)
		}
	}
	return water
}

func (b *Brain) GetWaterSupplyInMemory() Memory {
	if location, ok := b.GetClosestKnownWaterLocationFromCognitiveMap(); ok {
		return Memory{
			Event:    "Found water supply",
			Details:  "Water",
			Location: location,
		}
	}

	return b.GetClosestValidMemoryByEvent("Found water supply", b.hasWaterAtLocation)
}

func (b *Brain) FindClosestWaterSupply(water []Tile) Tile {
	closestWater := water[0]
	for _, tile := range water {
		if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, tile.Location) < b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, closestWater.Location) {
			closestWater = tile
		}
	}

	return closestWater
}

func (b *Brain) FindWaterSupply() bool {
	vision := b.GetWaterInVision()
	if len(vision) == 0 {
		b.GoSearchFor("Water supply")
		return false
	} else {
		closestWater := b.FindClosestWaterSupply(vision)
		b.AddMemoryToLongTerm("Found water supply", "Water", closestWater.Location)
		b.PhysiologicalNeeds.WayOfGettingWater = true
		return true
	}
}

// ----------------- Food -----------------

func (b *Brain) GetFoodInVision() []*Plant {
	vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
	plants := make([]*Plant, 0)

	for _, plant := range vision {
		if plant != nil &&
			plant.ProducesFruit &&
			HasRipeFruit(plant) &&
			!b.isLocationBlockedByOtherEntity(plant.Location) {
			plants = append(plants, plant)
		}
	}

	return plants
}

func (b *Brain) isCognitiveMapTileFresh(tile CognitiveMapTile, maxAge time.Duration) bool {
	if maxAge <= 0 {
		return true
	}
	if tile.LastSeenUnixMs <= 0 {
		// Backward compatibility with older map entries.
		return true
	}
	ageMs := time.Now().UnixMilli() - tile.LastSeenUnixMs
	return ageMs <= maxAge.Milliseconds()
}

func (b *Brain) GetClosestKnownLocationFromCognitiveMap(
	matches func(Location, CognitiveMapTile) bool,
	maxAge time.Duration,
) (Location, bool) {
	knownTiles := b.GetKnownTilesSnapshot()

	found := false
	bestDistance := 0
	bestLocation := Location{}

	for location, tile := range knownTiles {
		if !b.isCognitiveMapTileFresh(tile, maxAge) {
			continue
		}
		if !matches(location, tile) {
			continue
		}

		distance := b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, location)
		if !found || distance < bestDistance {
			found = true
			bestDistance = distance
			bestLocation = location
		}
	}

	return bestLocation, found
}

func (b *Brain) GetClosestKnownWaterLocationFromCognitiveMap() (Location, bool) {
	return b.GetClosestKnownLocationFromCognitiveMap(func(location Location, tile CognitiveMapTile) bool {
		return tile.TileType == Water && b.hasWaterAtLocation(location)
	}, cognitiveMapWaterMaxAge)
}

func (b *Brain) GetClosestKnownFoodLocationFromCognitiveMap() (Location, bool) {
	return b.GetClosestKnownLocationFromCognitiveMap(func(location Location, tile CognitiveMapTile) bool {
		hasFoodObservation := tile.Plant.IsAlive && tile.Plant.ProducesFruit && tile.Plant.HasRipeFruit && tile.Plant.FruitCount > 0
		return hasFoodObservation && b.hasFoodAtLocation(location)
	}, cognitiveMapFoodMaxAge)
}

func (b *Brain) hasWaterAtLocation(location Location) bool {
	return b.Owner.WorldProvider.GetTileType(location.X, location.Y) == Water &&
		!b.isLocationBlockedByOtherEntity(location)
}

func (b *Brain) hasFoodAtLocation(location Location) bool {
	if b.isLocationBlockedByOtherEntity(location) {
		return false
	}
	food := b.Owner.WorldProvider.GetFruitingPlantsInVision(location.X, location.Y, 0)
	return len(food) > 0
}

// GetClosestValidMemoryByEvent returns the closest memory matching an event whose location is still valid.
// Invalid locations are pruned from memory to avoid repeatedly targeting stale resources.
func (b *Brain) GetClosestValidMemoryByEvent(event string, isValid func(Location) bool) Memory {
	if len(b.Memories.LongTermMemory) == 0 && len(b.Memories.ShortTermMemory) == 0 {
		return Memory{}
	}

	currentLocation := b.Owner.Location
	found := false
	bestDistance := 0
	bestMemory := Memory{}
	staleLocations := make(map[Location]bool)

	checkMemory := func(memory Memory) {
		if memory.Event != event {
			return
		}
		if !isValid(memory.Location) {
			staleLocations[memory.Location] = true
			return
		}

		distance := b.Owner.WorldProvider.CalculateDistance(currentLocation, memory.Location)
		if !found || distance < bestDistance {
			found = true
			bestDistance = distance
			bestMemory = memory
		}
	}

	for _, memory := range b.Memories.ShortTermMemory {
		checkMemory(memory)
	}
	for _, memory := range b.Memories.LongTermMemory {
		checkMemory(memory)
	}

	for location := range staleLocations {
		b.RemoveMemoriesByEventAtLocation(event, location)
	}

	if found {
		return bestMemory
	}
	return Memory{}
}

func (b *Brain) GetFoodSupplyInMemory() Memory {
	if location, ok := b.GetClosestKnownFoodLocationFromCognitiveMap(); ok {
		return Memory{
			Event:    "Found food supply",
			Details:  "Food",
			Location: location,
		}
	}

	return b.GetClosestValidMemoryByEvent("Found food supply", b.hasFoodAtLocation)
}

func (b *Brain) FindFoodSupply() bool {
	vision := b.GetFoodInVision()
	if len(vision) == 0 {
		b.GoSearchFor("Food supply")
		return false
	} else {
		LogInfo(fmt.Sprintf("Found food supply at %d, %d", vision[0].Location.X, vision[0].Location.Y), b.Owner.FullName)
		closestPlant := b.FindClosestPlant(vision)
		b.AddMemoryToLongTerm("Found food supply", "Food", closestPlant.Location)
		b.PhysiologicalNeeds.WayOfGettingFood = true
		return true
	}
}

func (b *Brain) FindClosestPlant(plants []*Plant) *Plant {
	closestFood := plants[0]
	for _, plant := range plants {
		if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, plant.Location) < b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, closestFood.Location) {
			closestFood = plant
		}
	}

	return closestFood
}

// ----------------- Lumber -----------------

func (b *Brain) GetLumberInVision() []*Plant {
	vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)

	plants := make([]*Plant, 0)

	for _, plant := range vision {
		if plant != nil && plant.Name == OakTree && !b.isLocationBlockedByOtherEntity(plant.Location) {
			plants = append(plants, plant)
		}
	}

	return plants
}

func (b *Brain) hasLumberAtLocation(location Location) bool {
	if b.isLocationBlockedByOtherEntity(location) {
		return false
	}
	plants := b.Owner.WorldProvider.GetPlantsInVision(location.X, location.Y, 0)
	for _, plant := range plants {
		if plant != nil && plant.Name == OakTree {
			return true
		}
	}
	return false
}

func (b *Brain) GetClosestKnownLumberLocationFromCognitiveMap() (Location, bool) {
	return b.GetClosestKnownLocationFromCognitiveMap(func(location Location, tile CognitiveMapTile) bool {
		return tile.Plant.IsAlive && tile.Plant.Name == OakTree && b.hasLumberAtLocation(location)
	}, cognitiveMapLumberMaxAge)
}

func (b *Brain) GetLumberSupplyInMemory() Memory {
	if location, ok := b.GetClosestKnownLumberLocationFromCognitiveMap(); ok {
		return Memory{
			Event:    "Found lumber tree",
			Details:  "Lumber",
			Location: location,
		}
	}

	return b.GetClosestValidMemoryByEvent("Found lumber tree", b.hasLumberAtLocation)
}

func (b *Brain) FindLumberTrees() bool {
	trees := b.GetLumberInVision()
	if len(trees) == 0 {
		return false
	}

	closestTree := b.FindClosestPlant(trees)
	b.AddMemoryToLongTerm("Found lumber tree", "Lumber", closestTree.Location)
	return true
}

func (b *Brain) FindSticks() bool {
	vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
	if len(vision) == 0 {
		b.GoSearchFor("Sticks")
		return false
	} else {

	}
	return false
}

func (b *Brain) GetHighGrassInVision() []*Plant {
	vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
	plants := make([]*Plant, 0)

	for _, plant := range vision {
		if plant != nil &&
			plant.IsAlive &&
			plant.Name == HighGrass &&
			!b.isLocationBlockedByOtherEntity(plant.Location) {
			plants = append(plants, plant)
		}
	}

	return plants
}

func (b *Brain) hasGrassAtLocation(location Location) bool {
	if b.isLocationBlockedByOtherEntity(location) {
		return false
	}
	plants := b.Owner.WorldProvider.GetPlantsInVision(location.X, location.Y, 0)
	for _, plant := range plants {
		if plant != nil && plant.IsAlive && plant.Name == HighGrass {
			return true
		}
	}
	return false
}

func (b *Brain) GetClosestKnownGrassLocationFromCognitiveMap() (Location, bool) {
	return b.GetClosestKnownLocationFromCognitiveMap(func(location Location, tile CognitiveMapTile) bool {
		return tile.Plant.IsAlive && tile.Plant.Name == HighGrass && b.hasGrassAtLocation(location)
	}, cognitiveMapGrassMaxAge)
}

func (b *Brain) GetGrassSupplyInMemory() Memory {
	if location, ok := b.GetClosestKnownGrassLocationFromCognitiveMap(); ok {
		return Memory{
			Event:    "Found grass supply",
			Details:  "Grass",
			Location: location,
		}
	}

	return b.GetClosestValidMemoryByEvent("Found grass supply", b.hasGrassAtLocation)
}

func (b *Brain) FindGrassSupply() bool {
	grassPlants := b.GetHighGrassInVision()
	if len(grassPlants) == 0 {
		b.GoSearchFor("High grass")
		return false
	}

	closestGrass := b.FindClosestPlant(grassPlants)
	b.AddMemoryToLongTerm("Found grass supply", "Grass", closestGrass.Location)
	return true
}

// ----------------- Stone -----------------

func (b *Brain) GetStoneInVision() []Tile {
	vision := b.Owner.WorldProvider.GetVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
	stoneTiles := make([]Tile, 0)

	for _, tile := range vision {
		if b.isOwnBaseLocation(tile.Location) {
			continue
		}
		if b.isLocationBlockedByOtherEntity(tile.Location) {
			continue
		}
		for _, item := range tile.Items {
			if item != nil && item.Name == stoneItemName {
				stoneTiles = append(stoneTiles, tile)
				break
			}
		}
	}

	return stoneTiles
}

func (b *Brain) hasStoneAtLocation(location Location) bool {
	if b.isOwnBaseLocation(location) {
		return false
	}
	if b.isLocationBlockedByOtherEntity(location) {
		return false
	}

	tile := b.Owner.WorldProvider.GetTile(location.X, location.Y)
	for _, item := range tile.Items {
		if item != nil && item.Name == stoneItemName {
			return true
		}
	}
	return false
}

func (b *Brain) GetStoneSupplyInMemory() Memory {
	return b.GetClosestValidMemoryByEvent("Found stone supply", b.hasStoneAtLocation)
}

func (b *Brain) FindStoneSupply() bool {
	stoneTiles := b.GetStoneInVision()
	if len(stoneTiles) == 0 {
		b.GoSearchFor("Stone")
		return false
	}

	closestStone := stoneTiles[0]
	for _, tile := range stoneTiles {
		if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, tile.Location) <
			b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, closestStone.Location) {
			closestStone = tile
		}
	}

	b.AddMemoryToLongTerm("Found stone supply", "Stone", closestStone.Location)
	return true
}

func (b *Brain) isOwnBaseLocation(location Location) bool {
	return b.Base.Claimed && b.Base.Location == location
}

func (b *Brain) isLocationBlockedByOtherEntity(location Location) bool {
	tile := b.Owner.WorldProvider.GetTile(location.X, location.Y)
	return tile.Entity != nil && tile.Entity.FullName != b.Owner.FullName
}

// ----------------- Find ---------------------

type Direction struct {
	DX, DY int
}

var directions = []Direction{
	{DX: -1, DY: 0}, // Left
	{DX: 1, DY: 0},  // Right
	{DX: 0, DY: -1}, // Up
	{DX: 0, DY: 1},  // Down
}

func (b *Brain) IsValidLocation(loc Location) bool {
	return loc.X >= 0 && loc.X < SIZE_OF_MAP && loc.Y >= 0 && loc.Y < SIZE_OF_MAP
}

func (b *Brain) DecideLocationToSearch() Location {
	currentLocation := b.Owner.Location
	knownTiles := b.GetKnownTilesSnapshot()

	// Map to keep track of potential locations to explore
	frontier := make(map[Location]bool)

	// For each known tile, check its adjacent tiles
	for loc := range knownTiles {
		for _, dir := range directions {
			adjacentLoc := Location{X: loc.X + dir.DX, Y: loc.Y + dir.DY}
			if b.IsValidLocation(adjacentLoc) {
				// If the adjacent location is not known, add it to the frontier
				if _, known := knownTiles[adjacentLoc]; !known {
					frontier[adjacentLoc] = true
				}
			}
		}
	}

	// If there are frontier locations, choose one
	if len(frontier) > 0 {
		// Convert frontier map keys to a slice
		frontierLocations := make([]Location, 0, len(frontier))
		for loc := range frontier {
			frontierLocations = append(frontierLocations, loc)
		}

		// Sort frontier locations by distance to current location
		sort.Slice(frontierLocations, func(i, j int) bool {
			return b.Distance(frontierLocations[i], currentLocation) < b.Distance(frontierLocations[j], currentLocation)
		})

		// Return the closest unexplored location
		return frontierLocations[0]
	}

	// If no frontier locations, fallback to random movement or another strategy
	return b.RandomUnvisitedLocation()
}

// Helper function to calculate Manhattan distance between two locations
func (b *Brain) Distance(i, j Location) int {
	return abs(i.X-j.X) + abs(i.Y-j.Y)
}

// Helper function to get a random unvisited location in the world
func (b *Brain) RandomUnvisitedLocation() Location {
	for {
		randX := rand.Intn(SIZE_OF_MAP)
		randY := rand.Intn(SIZE_OF_MAP)
		loc := Location{X: randX, Y: randY}
		if !b.IsTileKnown(loc) {
			return loc
		}
	}
}

// Utility function
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// GoSearchFor - Go search for something - This assumes that the target isnt in memory or vision
func (b *Brain) GoSearchFor(target string) {
	targetLocation := b.DecideLocationToSearch()

	b.MotorCortexCurrentTask = MotorCortexAction{"Searching for " + target, "Walk", targetLocation, false, false}
}
