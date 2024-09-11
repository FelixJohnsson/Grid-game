package main

import "fmt"

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
        }
        fmt.Println("Found water supply.")
        closestWater := b.FindClosestWaterSupply(vision)
        b.AddMemoryToLongTerm("Found water supply", "Water", closestWater.Location)
        b.PhysiologicalNeeds.WayOfGettingWater = true
        b.MotorCortexCurrentTask.Finished = true
        b.MotorCortexCurrentTask.IsActive = false
        return true
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
    closestPlant := b.Owner.FindTheClosestPlant(vision)
    if closestPlant != nil {
        b.AddMemoryToLongTerm("Found food supply", "Food", closestPlant.Location)
        b.PhysiologicalNeeds.WayOfGettingFood = true
        return true
    }
    return false
}

// FindClosestFoodSupply - Find the closest food supply
func (b *Brain) FindClosestFoodSupply(food []Tile) Tile {
	closestFood := food[0]
	for _, tile := range food {
		if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location.X, b.Owner.Location.Y, tile.Location.X, tile.Location.Y) < b.Owner.WorldProvider.CalculateDistance(b.Owner.Location.X, b.Owner.Location.Y, closestFood.Location.X, closestFood.Location.Y) {
			closestFood = tile
		}
	}

	return closestFood
}


// ----------------- Find ---------------------

// GoSearchFor - Go search for something
func (b *Brain) GoSearchFor(target string) {
    switch target {
    case "Water supply":
        distanceToTravel := b.Owner.Curiosity
        currentLocation := b.Owner.Location
        
        // The grid is 100x100, with the center at (50, 50)
        xDistance := currentLocation.X - 50
        yDistance := currentLocation.Y - 50
        
        newLocation := currentLocation
        
        if xDistance > 0 && yDistance > 0 {
            //"Southwest"
            newLocation.X -= distanceToTravel
            newLocation.Y += distanceToTravel
        } else if xDistance > 0 && yDistance < 0 {
            //"Northwest"
            newLocation.X -= distanceToTravel
            newLocation.Y -= distanceToTravel
        } else if xDistance < 0 && yDistance > 0 {
            //"Southeast"
            newLocation.X += distanceToTravel
            newLocation.Y += distanceToTravel
        } else if xDistance < 0 && yDistance < 0 {
            //"Northeast"
            newLocation.X += distanceToTravel
            newLocation.Y -= distanceToTravel
        } else if xDistance > 0 {
            //"West" 
            newLocation.X -= distanceToTravel
        } else if xDistance < 0 {
            //"East" 
            newLocation.X += distanceToTravel
        } else if yDistance > 0 {
            //"South"
            newLocation.Y += distanceToTravel
        } else if yDistance < 0 {
            //"North"
            newLocation.Y -= distanceToTravel
        } else {
            //"At the center" 
        }
        
        // Ensure the new location stays within the grid 
        if newLocation.X < 1 {
            newLocation.X = 1
        } else if newLocation.X > SIZE_OF_MAP {
            newLocation.X = SIZE_OF_MAP
        }
        
        if newLocation.Y < 1 {
            newLocation.Y = 1
        } else if newLocation.Y > SIZE_OF_MAP {
            newLocation.Y = SIZE_OF_MAP
        }

        b.MotorCortexCurrentTask = MotorCortexAction{"Searching for a water supply", "Walk", Location{newLocation.X, newLocation.Y}, false, false}
    default:
        // Handle other targets
    }
}