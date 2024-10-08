package main

import (
	"fmt"
	"math/rand"
)

// ---------------- Actions ----------------

// ---------------- General actions ------------

// ClearAirway - Clear the airway of the person - Nose or Mouth
func (e *Entity) ClearAirway(action TargetedAction) {
    randomNumber := rand.Intn(100)

    if action.Target == "Mouth" && randomNumber < 20 {
        e.Body.Head.Mouth.IsObstructed = false
        e.Brain.RemoveActionFromActionList(action)
        return
    }
    if action.Target == "Nose" && randomNumber < 20 {
        e.Body.Head.Nose.IsObstructed = false
        e.Brain.RemoveActionFromActionList(action)
        return
    }
}

// FixNose - Fix the nose of the person
func (e *Entity) FixBrokenNose(action TargetedAction) {
    randomNumber := rand.Intn(100)

    if randomNumber < 20 {
        e.Body.Head.Nose.IsBroken = false
        e.Brain.RemoveActionFromActionList(action)
        e.Brain.ApplyPain(101)
    }
}

func (b *Brain) GetFoodForStorage(action TargetedAction) {

}




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

func (b *Brain) DrinkWaterTask(TargetedAction TargetedAction) {
	if b.CheckIfCurrentMotorTaskIsDone(b.MotorCortexCurrentTask, "Drink water") {
		water := Liquid{"Water"}
		b.Owner.Drink(water)
		return
	}
    water := b.GetAllWaterTilesFromCognitiveMap()

    if len(water) > 0 {
        closestWater := b.GetClosestWaterTileFromCognitiveMap(water)
        b.AddMotorCortexTask("Drink water", "Walk", Location{closestWater.Location.X, closestWater.Location.Y})
    } else {
        b.GoSearchFor("Water supply")
    }
}

func (b *Brain) EatFoodTask() {
    if b.CheckIfCurrentMotorTaskIsDone(b.MotorCortexCurrentTask, "Eat food") {
        food := b.Owner.WorldProvider.GetTile(b.Owner.Location.X, b.Owner.Location.Y)
        if food.Plant != nil && len(food.Plant.Fruit) > 0 {
            food := food.Plant.Fruit[0]
            b.Owner.Eat(food)
        } else {
            fmt.Println("I can't find food where I am, but the motor cortex thinks I've found food.")
        }
	}

	memorySuccess := b.GetFoodSupplyInMemory() 

    if memorySuccess.Event == "Found food supply" {
        b.AddMotorCortexTask("Eat food", "Walk", Location{memorySuccess.Location.X, memorySuccess.Location.Y})
        return
    }

    visionSuccess := b.FindFoodSupply()

    if visionSuccess {
        plants := b.GetFoodInVision()
        closestFood := b.FindClosestPlant(plants)
        b.AddMotorCortexTask("Eat food", "Walk", Location{closestFood.Location.X, closestFood.Location.Y})
    } else {
        b.GoSearchFor("Food supply")
    }
}

func (b *Brain) GetLumberTask() {
	success := b.FindLumberTrees()
	if success {
		trees := b.GetLumberInVision()
		closestTree := b.FindClosestPlant(trees)
		b.AddMotorCortexTask("Get lumber", "Walk", Location{closestTree.Location.X, closestTree.Location.Y})
	} else {
		b.GoSearchFor("Lumber tree")
	}
}

func (b *Brain) ChopDownTree(tree *Plant) *Item {
    if b.HasItemEquippedInRight("Stone Axe") {
        b.Owner.WorldProvider.RemovePlant(tree)
        b.Owner.DropFromRightHand("Stone Axe")
        wood := CreateNewItem("Wood log")
        b.Owner.GrabWithRightHand(wood)
		b.Owner.OwnedItems = append(b.Owner.OwnedItems, wood)
        return wood
    } else if b.HasItemEquippedInLeft("Stone Axe") {
        b.Owner.WorldProvider.RemovePlant(tree)
        b.Owner.DropFromLeftHand("Stone Axe")
        wood := CreateNewItem("Wood log")
        b.Owner.GrabWithLeftHand(wood)
		b.Owner.OwnedItems = append(b.Owner.OwnedItems, wood)
        return wood
    } else {
        return nil
    }
}

// ---------------- Food and water end tasks ----------------

// Eat - Consume food
func (e *Entity) Eat(food Food) {
    e.Brain.DecreaseHungerLevel(food.GetNutritionalValue())
}
// Drink - Consume a liquid
func (e *Entity) Drink(liquid Liquid) {
    switch liquid.Name {
    case "Water":
        e.Brain.DecreaseThirstLevel(50)
    }
}