package main

import (
	"math/rand"
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
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    if len(vision) == 0 {
        b.GoSearchFor("Food supply")
        return false
    } else {
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

func (b *Brain) DecideDirectionToSearch() Location {
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

    b.MotorCortexCurrentTask = MotorCortexAction{"Searching for " + target, "Walk", targetLocation, false, false}
}