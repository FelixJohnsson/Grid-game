package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// ----------------- Water -----------------

func (b *Brain) GetWaterInVision() []Tile {
    vision := b.Owner.WorldProvider.GetWaterInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)

    water := make([]Tile, 0)

    for _, tile := range vision {
        if tile.Type == 1 {
            water = append(water, tile)
        }
    }
    return water
}

func (b *Brain) GetWaterSupplyInMemory() Memory {
    if len(b.Memories.LongTermMemory) == 0 && len(b.Memories.ShortTermMemory) == 0 {
        return Memory{}
    }
    for _, memory := range b.Memories.LongTermMemory {
        if memory.Event == "Found water supply" {
            return memory
        }
    }
    for _, memory := range b.Memories.ShortTermMemory {
        if memory.Event == "Found water supply" {
            return memory
        }
    }
    return Memory{}
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
        vision := b.Owner.WorldProvider.GetWaterInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
        if len(vision) == 0 {
            b.GoSearchFor("Water supply")
            return false
        } else {
            b.PhysiologicalNeeds.WayOfGettingWater = true
            return true
        }
}

// ----------------- Food -----------------

func (b *Brain) GetFoodInVision() []*Plant {
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    plants := make([]*Plant, 0)

    for _, plant := range vision {
        if plant.Fruit != nil {
            plants = append(plants, plant)
        }
    }

    return plants
}

func (b *Brain) GetFoodSupplyInMemory() Memory {
    if len(b.Memories.LongTermMemory) == 0 && len(b.Memories.ShortTermMemory) == 0 {
        return Memory{}
    }
    for _, memory := range b.Memories.LongTermMemory {
        if memory.Event == "Found food supply" {
            return memory
        }
    }
    for _, memory := range b.Memories.ShortTermMemory {
        if memory.Event == "Found food supply" {
            return memory
        }
    }
    return Memory{}
}

func (b *Brain) FindFoodSupply() bool {
        vision := b.Owner.WorldProvider.GetFruitingPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
        if len(vision) == 0 {
            b.GoSearchFor("Food supply")
            return false
        } else {
            fmt.Println("Found food supply at", vision[0].Location.X, vision[0].Location.Y, vision[0].Fruit)
            b.PhysiologicalNeeds.WayOfGettingFood = true
            return true
        }
}

func (b *Brain) PredatorFindFood() bool {
    vision := b.Owner.WorldProvider.GetEntitiesInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)

    if len(vision) == 0 {
        return false
    } else {
        closestEntity := b.FindClosestEntity(vision)
        if closestEntity.Species == b.Owner.Species {
            b.GoSearchFor("Food")
            return false
        } else {
            fmt.Println("I found food at: ", closestEntity.Location.X, closestEntity.Location.Y)
            b.AddMotorCortexTask("Hunt", "Walk", Location{closestEntity.Location.X, closestEntity.Location.Y})
            return true
        }
    }
}

func (b *Brain) FindClosestEntity(entities []EntityInVision) EntityInVision {
    fmt.Println(Blue + "I see the following entities:", Reset)
    for _, entity := range entities {
        if entity.FullName != b.Owner.FullName {
            fmt.Println(entity.FirstName, entity.Species)
        }
    }
    closestEntity := entities[0]
    for _, entity := range entities {
        if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, entity.Location) < b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, closestEntity.Location) && entity.Species != b.Owner.Species {
            closestEntity = entity
        }
    }
    return closestEntity
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
        if plant.Name == OakTree {
            plants = append(plants, plant)
        }
    }

    return plants
}

func (b *Brain) FindLumberTrees() bool {
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    if len(vision) == 0 {
        b.GoSearchFor("Lumber tree")
        return false
    } else {
        
    }

    return false
}

// ----------------- Find ---------------------

type Direction struct {
    DX, DY int
}

var directions = []Direction{
    {DX: -1, DY: 0},  // Left
    {DX: 1, DY: 0},   // Right
    {DX: 0, DY: -1},  // Up
    {DX: 0, DY: 1},   // Down
}

func (b *Brain) IsValidLocation(loc Location) bool {
    return loc.X >= 0 && loc.X < SIZE_OF_MAP && loc.Y >= 0 && loc.Y < SIZE_OF_MAP
}

func (b *Brain) DecideLocationToSearch() Location {
    currentLocation := b.Owner.Location
    knownTiles := b.CognitiveMap.KnownTiles

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
        if _, known := b.CognitiveMap.KnownTiles[loc]; !known {
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

    b.AddMotorCortexTask("Searching for " + target, "Walk", targetLocation)
}