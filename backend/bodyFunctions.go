package main

import (
	"fmt"
	"math/rand"
)

// ---------------- Actions ----------------

// ClearAirway - Clear the airway of the person - Nose or Mouth
func (b *Brain) ClearAirway(action TargetedAction) {
    randomNumber := rand.Intn(100)
	fmt.Println("Trying to clear airway of the mouth")

    if action.Target == "Mouth" && randomNumber < 20 {
        b.Owner.Body.Head.Mouth.IsObstructed = false
        b.RemoveActionFromActionList(action)
        fmt.Println(b.Owner.FullName + " cleared the airway of the mouth.")
        return
    }
    if action.Target == "Nose" && randomNumber < 20 {
        b.Owner.Body.Head.Nose.IsObstructed = false
        b.RemoveActionFromActionList(action)
        fmt.Println(b.Owner.FullName + " cleared the airway of the nose.")
        return
    }
}

// FixNose - Fix the nose of the person
func (b *Brain) FixBrokenNose(action TargetedAction) {
    randomNumber := rand.Intn(100)
	fmt.Println("Trying to fix the nose")

    if randomNumber < 20 {
        b.Owner.Body.Head.Nose.IsBroken = false
        b.RemoveActionFromActionList(action)
        fmt.Println(b.Owner.FullName + " fixed the nose.")
        b.ApplyPain(101)
    }
}

// FindWaterSupply - Find a water supply
func (b *Brain) FindWaterSupply(action TargetedAction) {
    fmt.Println("Looking for water supply")
    success := b.Find("Clean water")
    if success {
        b.RemoveActionFromActionList(action)
    }

}

// FindFoodSupply - Find a food supply
func (b *Brain) FindFoodSupply(action TargetedAction) {
    fmt.Println("Looking for food supply")

    success := b.Find("Food supply")
    if success {
        b.RemoveActionFromActionList(action)
    }
}

// HasFruitsThatAreEdible checks the plant produces fruit and theyre ripe and edible
func (b *Brain) HasFruitsThatAreEdible(plant *Plant) bool {
    if plant != nil && plant.ProducesFruit && len(plant.Fruit) > 0 {
        for _, fruit := range plant.Fruit {
            if fruit.IsRipe {
                return true
            }
        }
    } else {
        return false
    }
    return false
}

// Find - Find a target
func (b *Brain) Find(target string) bool {
    // Find whatever the target is
    // For example, find water, find food, find shelter, find a person, etc.

    switch target {
    case "Clean water":
        fmt.Println("Check vision for water supply")
        vision := b.Owner.WorldProvider.GetWaterInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
        for _, TileInVision := range vision {
            if TileInVision.Tile.Type == 1 {
                fmt.Println("Found water supply at", TileInVision.Location.X, TileInVision.Location.Y)
                b.AddMemoryToLongTerm("Found water supply", "Water", TileInVision.Location)
                b.PhysiologicalNeeds.WayOfGettingWater = true
                return true
            }
        }
    case "Food supply":
        fmt.Println("Check vision for food supply")
        vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
        for _, PlantInVision := range vision {
            if PlantInVision.ProducesFruit {
                if b.HasFruitsThatAreEdible(PlantInVision) {
                    fmt.Println("Found food supply at", PlantInVision.Location.X, PlantInVision.Location.Y)
                    b.AddMemoryToLongTerm("Found food supply", "Food", PlantInVision.Location)
                    b.PhysiologicalNeeds.WayOfGettingFood = true
                    return true
                }
            }
        }
    case "Have food for storage":
        fmt.Println("Check vision for food supply")
    }

    return false
}