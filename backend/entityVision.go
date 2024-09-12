package main

import (
	"fmt"
	"math/rand"
)

// IsWaterInVision - Find a water supply in vision
func (b *Brain) IsWaterInVision() bool {
        vision := b.Owner.WorldProvider.GetWaterInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
        return len(vision) != 0
}

// GetWaterInVision - Find a water supply in vision
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

//IsWaterInMemory - Find a water supply in memory
func (b *Brain) IsWaterInMemory() bool {
    if len(b.Memories.LongTermMemory) == 0 && len(b.Memories.ShortTermMemory) == 0 {
        return false
    }
    for _, memory := range b.Memories.LongTermMemory {
        if memory.Event == "Found water supply" {
            return true
        }
    }
    for _, memory := range b.Memories.ShortTermMemory {
        if memory.Event == "Found water supply" {
            return true
        }
    }
    return false
}

// FindAndNoteWaterSupply - Find a water supply and add it to the memory
func (b *Brain) FindAndNoteWaterSupply() bool {
        vision := b.Owner.WorldProvider.GetWaterInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
        if len(vision) == 0 {
            b.GoSearchFor("Water supply")
            return false
        } else {
            fmt.Println(Red + "Found water supply." + Reset)
            closestWater := b.FindClosestWaterSupply(vision)
            b.AddMemoryToLongTerm("Found water supply", "Water", closestWater.Location)
            b.PhysiologicalNeeds.WayOfGettingWater = true
            b.MotorCortexCurrentTask.Finished = true
            b.MotorCortexCurrentTask.IsActive = false
            return true
        }
}

//FindClosestWater - Find the closest water from a list of water
func (b *Brain) FindClosestWaterSupply(water []Tile) Tile {
	closestWater := water[0]
	for _, tile := range water {
		if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location.X, b.Owner.Location.Y, tile.Location.X, tile.Location.Y) < b.Owner.WorldProvider.CalculateDistance(b.Owner.Location.X, b.Owner.Location.Y, closestWater.Location.X, closestWater.Location.Y) {
			closestWater = tile
		}
	}

	return closestWater
}

// IsFoodInVision - Find a food supply in vision
func (b *Brain) IsFoodInVision() bool {
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    return len(vision) != 0
}

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

// IsFoodInMemory - Find a food supply in memory
func (b *Brain) IsFoodInMemory() bool {
    if len(b.Memories.LongTermMemory) == 0 && len(b.Memories.ShortTermMemory) == 0 {
        return false
    }
    for _, memory := range b.Memories.LongTermMemory {
        if memory.Event == "Found food supply" {
            return true
        }
    }
    for _, memory := range b.Memories.ShortTermMemory {
        if memory.Event == "Found food supply" {
            return true
        }
    }
    return false
}

// FindFoodSupply - Find a food supply
func (b *Brain) FindAndNoteFoodSupply() bool {
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    if len(vision) == 0 {
        return false
    }
    fmt.Println(Red + "Found food supply." + Reset)
    closestPlant := b.Owner.FindTheClosestPlant(vision)
    if closestPlant != nil {
        b.AddMemoryToLongTerm("Found food supply", "Food", closestPlant.Location)
        b.PhysiologicalNeeds.WayOfGettingFood = true
        return true
    }
    return false
}

func (b *Brain) GetLumberInVision() []*Plant {
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)

    plants := make([]*Plant, 0)

    for _, plant := range vision {
        if plant.Name == "Oak Tree" {
            plants = append(plants, plant)
        }
    }

    return plants
}

func (b *Brain) FindAndNoteLumberTrees() bool {
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    lumberTrees := []*Plant{}

    for _, PlantInVision := range vision {
        if PlantInVision.Name == "Oak Tree" {
            lumberTrees = append(lumberTrees, PlantInVision)
        }
    }
    closestLumberTree := b.Owner.FindTheClosestPlant(lumberTrees)
    if closestLumberTree != nil {
        b.AddMemoryToLongTerm("Found lumber tree", "Lumber tree", closestLumberTree.Location)
        return true
    }
    return false
}



// FindClosestFoodSupply - Find the closest food supply
func (b *Brain) FindClosestPlant(plants []*Plant) *Plant {
	closestFood := plants[0]
	for _, plant := range plants {
		if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location.X, b.Owner.Location.Y, plant.Location.X, plant.Location.Y) < b.Owner.WorldProvider.CalculateDistance(b.Owner.Location.X, b.Owner.Location.Y, closestFood.Location.X, closestFood.Location.Y) {
			closestFood = plant
		}
	}

	return closestFood
}


// ----------------- Find ---------------------

func (b *Brain) DecideDirectionToSearch() Location {
    // Size of the map (constant)
    const SIZE_OF_MAP = 100

    // Current location and curiosity factor
    distanceToTravel := b.Owner.Curiosity
    currentLocation := b.Owner.Location

    // Coordinates relative to the center (50, 50)
    relativeXPos := float64(currentLocation.X - SIZE_OF_MAP/2) / float64(SIZE_OF_MAP/2)
    relativeYPos := float64(currentLocation.Y - SIZE_OF_MAP/2) / float64(SIZE_OF_MAP/2)

    // Initialize probabilities for each direction (North, South, East, West)
    probabilityNorth := 0.0
    probabilitySouth := 0.0
    probabilityEast := 0.0
    probabilityWest := 0.0

    // Calculate probabilities based on the current position
    if relativeYPos < 0 { // Entity is in the northern part
        probabilitySouth = -relativeYPos // More likely to go South
    } else { // Entity is in the southern part
        probabilityNorth = relativeYPos // More likely to go North
    }

    if relativeXPos < 0 { // Entity is in the western part
        probabilityEast = -relativeXPos // More likely to go East
    } else { // Entity is in the eastern part
        probabilityWest = relativeXPos // More likely to go West
    }

    // Normalize the probabilities so they sum to 1
    totalProbability := probabilityNorth + probabilitySouth + probabilityEast + probabilityWest

    if totalProbability > 0 {
        probabilityNorth /= totalProbability
        probabilitySouth /= totalProbability
        probabilityEast /= totalProbability
        probabilityWest /= totalProbability
    }

    // Pick a direction based on the weighted probabilities
    randValue := rand.Float64()
    var direction string

    if randValue < probabilityNorth {
        direction = "North"
    } else if randValue < probabilityNorth + probabilitySouth {
        direction = "South"
    } else if randValue < probabilityNorth + probabilitySouth + probabilityEast {
        direction = "East"
    } else {
        direction = "West"
    }

    // Calculate the new location based on the chosen direction and distance to travel
    newLocation := currentLocation
    switch direction {
    case "North":
        newLocation.Y = max(0, currentLocation.Y-distanceToTravel)
    case "South":
        newLocation.Y = min(SIZE_OF_MAP, currentLocation.Y+distanceToTravel)
    case "East":
        newLocation.X = min(SIZE_OF_MAP, currentLocation.X+distanceToTravel)
    case "West":
        newLocation.X = max(0, currentLocation.X-distanceToTravel)
    }

    return newLocation
}

// GoSearchFor - Go search for something
func (b *Brain) GoSearchFor(target string) {
    targetLocation := b.DecideDirectionToSearch()   

    switch target {
    case "Water supply":
        b.MotorCortexCurrentTask = MotorCortexAction{"Searching for a water supply", "Walk", targetLocation, false, false}
    case "Food supply":
        b.MotorCortexCurrentTask = MotorCortexAction{"Searching for a food supply", "Walk", targetLocation, false, false}
    case "Lumber tree":
        b.MotorCortexCurrentTask = MotorCortexAction{"Searching for a lumber tree", "Walk", targetLocation, false, false}
    default:
        // Handle other targets
    }
}