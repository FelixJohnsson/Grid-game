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
	fmt.Println("Trying to clear airway of the mouth")

    if action.Target == "Mouth" && randomNumber < 20 {
        e.Body.Head.Mouth.IsObstructed = false
        e.Brain.RemoveActionFromActionList(action)
        fmt.Println(e.FullName + " cleared the airway of the mouth.")
        return
    }
    if action.Target == "Nose" && randomNumber < 20 {
        e.Body.Head.Nose.IsObstructed = false
        e.Brain.RemoveActionFromActionList(action)
        fmt.Println(e.FullName + " cleared the airway of the nose.")
        return
    }
}

// FixNose - Fix the nose of the person
func (e *Entity) FixBrokenNose(action TargetedAction) {
    randomNumber := rand.Intn(100)
	fmt.Println("Trying to fix the nose")

    if randomNumber < 20 {
        e.Body.Head.Nose.IsBroken = false
        e.Brain.RemoveActionFromActionList(action)
        fmt.Println(e.FullName + " fixed the nose.")
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

	success := b.FindWaterSupply()

	if success {
		water := b.GetWaterInVision()
		closestWater := b.FindClosestWaterSupply(water)
		b.MotorCortexCurrentTask = MotorCortexAction{"Drink water", "Walk", Location{closestWater.Location.X, closestWater.Location.Y}, false, false}
	} else {
		fmt.Println("I can't find a water supply.")
		b.GoSearchFor("Water supply")
	}
}

func (b *Brain) EatFoodTask(TargetedAction TargetedAction) {
		if b.CheckIfCurrentMotorTaskIsDone(b.MotorCortexCurrentTask, "Eat food") {
		
		return
	}

	success := b.FindFoodSupply()

	if success {
		plants := b.GetFoodInVision()
		closestFood := b.FindClosestPlant(plants)
		b.MotorCortexCurrentTask = MotorCortexAction{"Eat food", "Walk", Location{closestFood.Location.X, closestFood.Location.Y}, false, false}
	} else {
		fmt.Println("I can't find a food supply.")
		b.GoSearchFor("Food supply")
	}
}

func (b *Brain) GetLumberTask(TargetedAction TargetedAction) {
	success := b.FindLumberTrees()
	fmt.Println(Red + "Get Lumber Task: " + Reset)
	if success {
		trees := b.GetLumberInVision()
		closestTree := b.FindClosestPlant(trees)
		b.MotorCortexCurrentTask = MotorCortexAction{"Get lumber", "Walk", Location{closestTree.Location.X, closestTree.Location.Y}, false, false}
		fmt.Println(Green + "I can find a lumber tree." + Reset)
	} else {
		fmt.Println("I can't find a lumber tree.")
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
        fmt.Println("I chopped down the tree.")
        return wood
    } else if b.HasItemEquippedInLeft("Stone Axe") {
        b.Owner.WorldProvider.RemovePlant(tree)
        b.Owner.DropFromLeftHand("Stone Axe")
        wood := CreateNewItem("Wood log")
        b.Owner.GrabWithLeftHand(wood)
		b.Owner.OwnedItems = append(b.Owner.OwnedItems, wood)
        fmt.Println("I chopped down the tree.")
        return wood
    } else {
        fmt.Println("I need a stone axe to chop down the tree.")
        return nil
    }
}

// ---------------- Food and water end tasks ----------------

// Eat - Consume food
func (e *Entity) Eat(food Food) {
	switch food.GetName() {
	case "Apple":
		fmt.Println("Eating an apple")
		e.Brain.DecreaseHungerLevel(food.GetNutritionalValue())
	}
}
// Drink - Consume a liquid
func (e *Entity) Drink(liquid Liquid) {
    switch liquid.Name {
    case "Water":
        fmt.Println("Drinking water") 
        e.Brain.DecreaseThirstLevel(50)
    }
}