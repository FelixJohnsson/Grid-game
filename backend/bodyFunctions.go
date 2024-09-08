package main

import (
	"fmt"
)

// ---------------- Actions ----------------

// MotorCortexFindDrinkWater - Drink water
func (e *Entity) MotorCortexFindDrinkWater() {
    vision := e.Owner.WorldProvider.GetWaterInVision(e.Owner.Location.X, e.Owner.Location.Y, e.Owner.VisionRange)
    if len(vision) == 0 {
        fmt.Println("I can't see any water. Do I remember where I saw water last time?")
        return
    }

    closestWater := e.Owner.FindClosestWaterSupply(vision)

    path := e.DecidePathTo(closestWater.Location.X, closestWater.Location.Y)
    if path == nil {
        fmt.Println("I can't find a path to the water.")
    } else {
        e.MotorCortexCurrentTask = MotorCortexAction{"Drink water", "Walk", Location{closestWater.Location.X, closestWater.Location.Y}, false, false}
    }
}



// HasFruitsThatAreEdible checks the plant produces fruit and theyre ripe and edible
func (e *Entity) HasFruitsThatAreEdible(plant *Plant) bool {
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

// EatFruit - Eat fruit from a fruit tree
func (e *Entity) EatFruit() {
    vision := e.Owner.WorldProvider.GetPlantsInVision(e.Owner.Location.X, e.Owner.Location.Y, e.Owner.VisionRange)
    if len(vision) == 0 {
        fmt.Println("I can't see any food. Do I remember where I saw food last time?")
        return
    }

    fruitTrees := []*Plant{}
    for _, PlantInVision := range vision {
        if PlantInVision.ProducesFruit {
            if e.HasFruitsThatAreEdible(PlantInVision) {
                fruitTrees = append(fruitTrees, PlantInVision)
            }
        }
    }

    closestFruitTree := e.Owner.FindTheClosestPlant(fruitTrees)
    if closestFruitTree != nil {
        path := e.DecidePathTo(closestFruitTree.Location.X, closestFruitTree.Location.Y)
        if path == nil {
            fmt.Println("I can't find a path to the fruit tree.")
        } else {
            e.WalkOverPath(closestFruitTree.Location.X, closestFruitTree.Location.Y)
            fruit := closestFruitTree.Fruit[0]
            e.Owner.Eat(fruit)
            e.PhysiologicalNeeds.WayOfGettingFood = true
        }
    }
}

// GetFoodForStorage - Get food for storage
func (e *Entity) GetFoodForStorage(action TargetedAction) {
    hasFoodStorage := e.FindInOwnedItems("Food Box")

    if hasFoodStorage == nil {
        fmt.Println("I need a food storage to store food.")
        e.GetWood()

        foodBox := e.Craft("Food Box")
        e.Owner.GrabRight(foodBox)
        e.Owner.OwnedItems = append(e.Owner.OwnedItems, foodBox)
        e.RemoveActionFromActionList(action)
        e.PhysiologicalNeeds.FoodSupply = true
        fmt.Println("I got the wood and crafted a food box.")
    }
}

// Craft - Craft an item
func (e *Entity) Craft(item string) *Item {
    switch item {
    case "Stone Axe":
            stoneAxe := Item{"Stone Axe", 1, 10, 5, []Material{materials[0]}, make([]Residue, 0), Location{e.Owner.Location.X, e.Owner.Location.Y}}
        return &stoneAxe
    case "Food Box":
        foodBox := Item{"Food Box", 1, 1, 1, []Material{materials[0]}, make([]Residue, 0), Location{e.Owner.Location.X, e.Owner.Location.Y}}
        return &foodBox
    }
    return nil
}

// GetWood - Get wood
func (e *Entity) GetWood() *Item {
    nearestTree := e.FindWood()
    if nearestTree != nil {
        fmt.Println("I found wood. Now I need to get it.")
        e.WalkOverPath(nearestTree.Location.X, nearestTree.Location.Y)
        tree := e.Owner.WorldProvider.GetTile(nearestTree.Location.X, nearestTree.Location.Y).Plant
        woodLog := e.ChopDownTree(tree)
        return woodLog
    } else {
        fmt.Println("I can't see any wood. Do I remember where I saw wood last time?")
        return nil
    }
}


// FindWood - Find wood
func (e *Entity) FindWood() *Plant {
    fmt.Println("Check vision for wood")
    vision := e.Owner.WorldProvider.GetPlantsInVision(e.Owner.Location.X, e.Owner.Location.Y, e.Owner.VisionRange)
    lumberTrees := []*Plant{}

    for _, PlantInVision := range vision {
        if PlantInVision.Name == "Oak Tree" {
            lumberTrees = append(lumberTrees, PlantInVision)
        }
    }
    closestLumberTree := e.Owner.FindTheClosestPlant(lumberTrees)
    if closestLumberTree != nil {
        fmt.Println("Found wood at", closestLumberTree.Location.X, closestLumberTree.Location.Y)
        return closestLumberTree
    }
    return nil
}

// ChopDownTree - Chop down a tree
func (e *Entity) ChopDownTree(tree *Plant) *Item {
    if e.HasItemEquippedInRight("Stone Axe") {
        e.Owner.WorldProvider.RemovePlant(tree)
        e.Owner.DropRight("Stone Axe")
        wood := CreateNewItem("Wood log")
        e.Owner.GrabRight(wood)
        fmt.Println("I chopped down the tree.")
        return wood
    } else if e.HasItemEquippedInLeft("Stone Axe") {
        e.Owner.WorldProvider.RemovePlant(tree)
        e.Owner.DropLeft("Stone Axe")
        wood := CreateNewItem("Wood log")
        e.Owner.GrabLeft(wood)
        fmt.Println("I chopped down the tree.")
        return wood
    } else {
        fmt.Println("I need a stone axe to chop down the tree.")
        return nil
    }
}

// ConstructShelter - Construct a shelter
func (e *Entity) ConstructShelter() bool {
    // Find a Grass tile that is empty and construct a shelter there
    vision := e.Owner.WorldProvider.GetGrassInVision(e.Owner.Location.X, e.Owner.Location.Y, e.Owner.VisionRange)
    closestGrass := e.Owner.FindClosestEmptyGrass(vision)

    path := e.DecidePathTo(closestGrass.Location.X, closestGrass.Location.Y)
    if path == nil {
        fmt.Println("I can't find a path to the grass tile.")
        return false
    }

    e.WalkOverPath(closestGrass.Location.X, closestGrass.Location.Y)

    newShelter := NewShelter(e.Owner.Location.X, e.Owner.Location.Y, e.Owner)
    hasWoodLog := e.FindInOwnedItems("Wood log")
    if hasWoodLog != nil {
        e.Owner.WorldProvider.AddShelter(e.Owner.Location.X, e.Owner.Location.Y, newShelter)
        return true
    } else {
        return false
    }
}

// MakeShelter - Make a shelter
func (e *Entity) MakeShelter(action TargetedAction) {
    hasWoodLog := e.FindInOwnedItems("Wood log")

    if hasWoodLog == nil {
        fmt.Println("I need wood to make a shelter.")
        e.GetWood()
        e.ConstructShelter()
    } else {
        fmt.Println("I have wood. Now I can make a shelter.")
        success := e.ConstructShelter()
        if success {
            e.RemoveActionFromActionList(action)
            e.PhysiologicalNeeds.HasShelter = true
            e.AddMemoryToLongTerm("Made a shelter", "Shelter", e.Owner.Location)
            fmt.Println("I made a shelter.")
        } else {
            fmt.Println("I couldnt make a shelter here.")
        }
    }
}

// ImproveDefense - Improve defense
func (e *Entity) ImproveDefense() bool {
    // Check the tile it's standing on
    isTileEmpty := e.IsTileEmpty(e.Owner.Location.X, e.Owner.Location.Y)

    if isTileEmpty {
        e.Owner.CombatSkill += 1
        fmt.Println("My combat skill is now", e.Owner.CombatSkill)
    } else {

    vision := e.Owner.WorldProvider.GetGrassInVision(e.Owner.Location.X, e.Owner.Location.Y, e.Owner.VisionRange)
    closestGrass := e.Owner.FindClosestEmptyGrass(vision)
    path := e.DecidePathTo(closestGrass.Location.X, closestGrass.Location.Y)
    
    if path == nil {
        fmt.Println("I can't find a path to the grass tile.")
        return false
    }
        e.WalkOverPath(closestGrass.Location.X, closestGrass.Location.Y)
        e.ImproveDefense()
    }

    return true
}

