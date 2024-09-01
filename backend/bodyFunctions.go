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

// GetFoodForStorage - Get food for storage
func (b *Brain) GetFoodForStorage(action TargetedAction) {
    fmt.Println("Do I have a storage where I can store food?")

    storage := ""
    basket := false

    if b.Owner.OwnedItems != nil {
        for _, item := range b.Owner.OwnedItems {
            if item.Name == "Wooden Crate" || item.Name == "Wooden Box" {
                storage = item.Name
                fmt.Println("I have a" + storage + "where I can store food.")
            } else {
                fmt.Println("I don't have a storage where I can store food. So I should make one.")
                fmt.Println("I should find materials to make a storage.")

                success := b.FindWood()
                if success != nil {
                    fmt.Println("I found wood. Now I need to get it.")
                    b.CurrentTask = TargetedAction{"Walk", "Wood", true, []BodyPartType{"RightLeg", "LeftLeg"}, 100}
                                                                                                                                          
                    go b.WalkOverPath(success.Location.X, success.Location.Y)
                }

            }
            if item.Name == "Woven Grass Basket" {
                basket = true
            }
        }
    }

    fmt.Println("Do I have a basket or something to carry food?")
    if basket {
        fmt.Println("I have a basket. Then I should find food.")

        success := b.Find("Food supply")
        if success {
            fmt.Println("I found food. Now I need to get it.")
        } else {
            fmt.Println("I can't see any food. Do I remember where I saw food last time?")

            if len(b.Memories.LongTermMemory) > 0 {
                for _, memory := range b.Memories.LongTermMemory {
                    if memory.Event == "Found food supply" {
                        fmt.Println("I remember where I saw food last time.")
                        fmt.Println("I should go there and get food.")

                        b.CurrentTask = TargetedAction{"Walk", memory.Details, true, []BodyPartType{"RightLeg", "LeftLeg"}, 100}
                        // Make a go routine to walk to the location
                        go b.WalkOverPath(memory.Location.X, memory.Location.Y)
                    }
                }
            }
    }
    }
}

// GetWood - Get wood
func (b *Brain) GetWood(action TargetedAction) {
    
}


// FindWood - Find wood
func (b *Brain) FindWood() *Plant {
    fmt.Println("Check vision for wood")
    vision := b.Owner.WorldProvider.GetPlantsInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    lumberTrees := []*Plant{}

    for _, PlantInVision := range vision {
        if PlantInVision.Name == "Oak Tree" {
            lumberTrees = append(lumberTrees, PlantInVision)
        }
    }
    closestLumberTree := b.Owner.FindTheClosestPlant(lumberTrees)
    if closestLumberTree != nil {
        fmt.Println("Found wood at", closestLumberTree.Location.X, closestLumberTree.Location.Y)
        return closestLumberTree
    }
    return nil
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

    }
    return false
}