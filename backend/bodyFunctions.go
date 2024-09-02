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


// DrinkWater - Drink water
func (b *Brain) DrinkWater(action TargetedAction) {
    vision := b.Owner.WorldProvider.GetWaterInVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
    if len(vision) == 0 {
        fmt.Println("I can't see any water. Do I remember where I saw water last time?")
        return
    }

    closestWater := b.Owner.FindClosestWaterSupply(vision)

    path := b.DecidePathTo(closestWater.Location.X, closestWater.Location.Y)
    if path == nil {
        fmt.Println("I can't find a path to the water.")
    } else {
        b.WalkOverPath(closestWater.Location.X, closestWater.Location.Y)
        water := Liquid{"Water"}
        b.Owner.Drink(water)
        b.PhysiologicalNeeds.WayOfGettingWater = true
        if b.PhysiologicalNeeds.Thirst < 20 {
            b.RemoveActionFromActionList(action)
        }
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
    hasFoodStorage := b.FindInOwnedItems("Food Box")

    if hasFoodStorage == nil {
        fmt.Println("I need a food storage to store food.")
        woodLog := b.GetWood()

        foodBox := b.Craft("Food Box")
        b.Owner.WorldProvider.DestroyItem(woodLog)
        b.Owner.GrabRight(foodBox)
        b.Owner.OwnedItems = append(b.Owner.OwnedItems, foodBox)
        b.RemoveActionFromActionList(action)
        b.PhysiologicalNeeds.FoodSupply = true
        fmt.Println("I got the wood and crafted a food box.")
    }
}

// Craft - Craft an item
func (b *Brain) Craft(item string) *Item {
    switch item {
    case "Stone Axe":
            stoneAxe := Item{"Stone Axe", 1, 10, 5, []Material{materials[0]}, make([]Residue, 0), Location{b.Owner.Location.X, b.Owner.Location.Y}}
        return &stoneAxe
    case "Food Box":
        foodBox := Item{"Food Box", 1, 1, 1, []Material{materials[0]}, make([]Residue, 0), Location{b.Owner.Location.X, b.Owner.Location.Y}}
        return &foodBox
    }
    return nil
}

// GetWood - Get wood
func (b *Brain) GetWood() *Item {
    nearestTree := b.FindWood()
    if nearestTree != nil {
        fmt.Println("I found wood. Now I need to get it.")
        b.WalkOverPath(nearestTree.Location.X, nearestTree.Location.Y)
        woodLog := b.ChopDownTree()
        return woodLog
    } else {
        fmt.Println("I can't see any wood. Do I remember where I saw wood last time?")
        return nil
    }
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

// ChopDownTree - Chop down a tree
func (b *Brain) ChopDownTree() *Item {
    if b.HasItemEquippedInRight("Stone Axe") {
        fmt.Println("I have a stone axe in my right hand. I can chop down the tree.")
        b.Owner.DropRight("Stone Axe")
        wood := CreateNewItem("Wood log")
        b.Owner.GrabRight(wood)
        fmt.Println("I chopped down the tree.")
        return wood
    } else if b.HasItemEquippedInLeft("Stone Axe") {
        fmt.Println("I have a stone axe in my left hand. I can chop down the tree.")
        b.Owner.DropLeft("Stone Axe")
        wood := CreateNewItem("Wood log")
        b.Owner.GrabLeft(wood)
        fmt.Println("I chopped down the tree.")
        return wood
    } else {
        fmt.Println("I need a stone axe to chop down the tree.")
        return nil
    }
}

// ConstructShelter - Construct a shelter
func (b *Brain) ConstructShelter() bool {
    newShelter := NewShelter(b.Owner.Location.X, b.Owner.Location.Y, b.Owner)
    hasWoodLog := b.FindInOwnedItems("Wood log")
    if hasWoodLog != nil {
        b.Owner.WorldProvider.AddShelter(b.Owner.Location.X, b.Owner.Location.Y, newShelter)
        return true
    } else {
        return false
    }
}

// MakeShelter - Make a shelter
func (b *Brain) MakeShelter(action TargetedAction) {
    hasWoodLog := b.FindInOwnedItems("Wood log")

    if hasWoodLog == nil {
        fmt.Println("I need wood to make a shelter.")
        b.GetWood()
        b.ConstructShelter()
    } else {
        fmt.Println("I have wood. Now I can make a shelter.")
        success := b.ConstructShelter()
        if success {
            b.RemoveActionFromActionList(action)
            b.PhysiologicalNeeds.HasShelter = true
            b.AddMemoryToLongTerm("Made a shelter", "Shelter", b.Owner.Location)
            fmt.Println("I made a shelter.")
        } else {
            fmt.Println("I couldnt make a shelter here.")
        }
    }
}

// ImproveDefense - Improve defense
func (b *Brain) ImproveDefense(action TargetedAction) {
    fmt.Println("I'm training to improve my defense.")

    b.Owner.CombatSkill += 1
    fmt.Println("My combat skill is now", b.Owner.CombatSkill)
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
            if TileInVision.Type == 1 {
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
    case "Shelter":
        return false
    }
    return false
}
