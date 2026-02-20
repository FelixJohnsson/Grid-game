package main

import (
	"fmt"
	"math/rand"
)

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
	b.CollectFoodForBaseTask(action)
}

func (b *Brain) Craft(item string) *Item {
	crafted := CreateNewItem(item)
	if crafted == nil {
		return nil
	}
	crafted.Location = b.Owner.Location
	return crafted
}

func (b *Brain) DrinkWaterTask(TargetedAction TargetedAction) {
	if b.CheckIfCurrentMotorTaskIsDone(b.MotorCortexCurrentTask, "Drink water") {
		// Validate we actually reached water before drinking.
		if b.Owner.WorldProvider.GetTileType(b.Owner.Location.X, b.Owner.Location.Y) == Water {
			water := Liquid{"Water"}
			b.Owner.Drink(water)
			return
		}

		// Target is stale; clear water memory for this location and continue searching.
		b.RemoveMemoriesByEventAtLocation("Found water supply", b.Owner.Location)
	}

	// If the remembered closest water tile is blocked, prune it and retry
	// so we can fall back to the next closest known tile.
	for attempt := 0; attempt < 4; attempt++ {
		memorySuccess := b.GetWaterSupplyInMemory()
		if memorySuccess.Event != "Found water supply" {
			break
		}
		if b.canOccupyLocation(memorySuccess.Location) {
			b.MotorCortexCurrentTask = MotorCortexAction{"Drink water", "Walk", Location{memorySuccess.Location.X, memorySuccess.Location.Y}, false, false}
			return
		}
		b.RemoveMemoriesByEventAtLocation("Found water supply", memorySuccess.Location)
	}

	success := b.FindWaterSupply()

	if success {
		water := b.GetWaterInVision()
		if closestWater, ok := b.FindClosestReachableWaterSupply(water); ok {
			b.MotorCortexCurrentTask = MotorCortexAction{"Drink water", "Walk", Location{closestWater.Location.X, closestWater.Location.Y}, false, false}
			return
		}
	}
	b.GoSearchFor("Water supply")
}

func (b *Brain) EatFoodTask() {
	if b.CheckIfCurrentMotorTaskIsDone(b.MotorCortexCurrentTask, "Eat food") {
		edible, err := b.Owner.WorldProvider.ConsumeRipeFruitAt(b.Owner.Location.X, b.Owner.Location.Y)
		if err == nil {
			b.Owner.Eat(edible)
			b.RemoveMemoriesByEventAtLocation("Found food supply", b.Owner.Location)
			return
		}

		// The remembered/targeted food source is stale; clear it and continue searching.
		b.RemoveMemoriesByEventAtLocation("Found food supply", b.Owner.Location)
		fmt.Println("No ripe food at target location anymore; clearing stale food memory.")
	}

	memorySuccess := b.GetFoodSupplyInMemory()

	if memorySuccess.Event == "Found food supply" {
		b.MotorCortexCurrentTask = MotorCortexAction{"Eat food", "Walk", Location{memorySuccess.Location.X, memorySuccess.Location.Y}, false, false}
		return
	}

	visionSuccess := b.FindFoodSupply()

	if visionSuccess {
		plants := b.GetFoodInVision()
		closestFood := b.FindClosestPlant(plants)
		b.MotorCortexCurrentTask = MotorCortexAction{"Eat food", "Walk", Location{closestFood.Location.X, closestFood.Location.Y}, false, false}
	} else {
		b.GoSearchFor("Food supply")
	}
}

func (b *Brain) GetLumberTask() {
	memorySuccess := b.GetLumberSupplyInMemory()
	if memorySuccess.Event == "Found lumber tree" {
		b.MotorCortexCurrentTask = MotorCortexAction{"Get lumber", "Walk", Location{memorySuccess.Location.X, memorySuccess.Location.Y}, false, false}
		return
	}

	success := b.FindLumberTrees()
	if success {
		trees := b.GetLumberInVision()
		closestTree := b.FindClosestPlant(trees)
		b.MotorCortexCurrentTask = MotorCortexAction{"Get lumber", "Walk", Location{closestTree.Location.X, closestTree.Location.Y}, false, false}
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
	e.Body.Torso.Stomach.Contains = append(e.Body.Torso.Stomach.Contains, food)
}

// Drink - Consume a liquid
func (e *Entity) Drink(liquid Liquid) {
	switch liquid.Name {
	case "Water":
		e.Body.Blood.Water += 50
	}
}

func (b *Brain) canOccupyLocation(location Location) bool {
	if location == b.Owner.Location {
		return true
	}

	tile := b.Owner.WorldProvider.GetTile(location.X, location.Y)
	if tile.Entity != nil && tile.Entity.FullName != b.Owner.FullName {
		return false
	}

	path := b.DecidePathTo(location.X, location.Y)
	return len(path) > 0
}

func (b *Brain) FindClosestReachableWaterSupply(water []Tile) (Tile, bool) {
	bestDistance := 0
	bestTile := Tile{}
	found := false

	for _, tile := range water {
		if !b.canOccupyLocation(tile.Location) {
			continue
		}
		distance := b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, tile.Location)
		if !found || distance < bestDistance {
			bestDistance = distance
			bestTile = tile
			found = true
		}
	}

	return bestTile, found
}
