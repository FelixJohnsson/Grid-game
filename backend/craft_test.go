package main

import "testing"

func TestGetCraftingRecipeIncludesStoneAgeEntries(t *testing.T) {
	recipe, ok := GetCraftingRecipe("Stone Pickaxe")
	if !ok {
		t.Fatalf("expected stone pickaxe recipe to exist")
	}
	if recipe.OutputItem != "Stone Pickaxe" {
		t.Fatalf("expected output Stone Pickaxe, got %q", recipe.OutputItem)
	}
	if recipe.Ingredients[stoneStockpileKey] <= 0 || recipe.Ingredients[stickItemName] <= 0 {
		t.Fatalf("expected stone pickaxe recipe to require stone and sticks, got %+v", recipe.Ingredients)
	}
}

func TestCraftItemConsumesStockpileAndDropsResultAtBase(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 5, Y: 5}
	brain.ensureBaseStockpile()
	brain.Base.Stockpile[stoneStockpileKey] = 2
	brain.Base.Stockpile[stickItemName] = 2

	if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
		t.Fatalf("failed to move entity to base: %v", err)
	}

	action := TargetedAction{Action: CraftItem, Target: "Stone Axe"}
	brain.CraftItem(action)

	if brain.Base.Stockpile[stoneStockpileKey] != 0 {
		t.Fatalf("expected stone stockpile to be consumed, got %d", brain.Base.Stockpile[stoneStockpileKey])
	}
	if brain.Base.Stockpile[stickItemName] != 0 {
		t.Fatalf("expected stick stockpile to be consumed, got %d", brain.Base.Stockpile[stickItemName])
	}

	baseTile := world.GetTile(brain.Base.Location.X, brain.Base.Location.Y)
	if got := countItemsOnTileByName(baseTile, "Stone Axe"); got != 1 {
		t.Fatalf("expected one crafted Stone Axe on base tile, got %d", got)
	}
}

func TestCraftItemRequiresShelterWhenRecipeDemandsIt(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 5, Y: 5}
	brain.ensureBaseStockpile()
	brain.PhysiologicalNeeds.HasShelter = false
	brain.Base.Stockpile[stoneStockpileKey] = 3
	brain.Base.Stockpile[stickItemName] = 6

	if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
		t.Fatalf("failed to move entity to base: %v", err)
	}

	action := TargetedAction{Action: CraftItem, Target: "Food Box"}
	brain.CraftItem(action)

	baseTile := world.GetTile(brain.Base.Location.X, brain.Base.Location.Y)
	if got := countItemsOnTileByName(baseTile, "Food Box"); got != 0 {
		t.Fatalf("expected no crafted Food Box without shelter, got %d", got)
	}
	if brain.Base.Stockpile[stoneStockpileKey] != 3 || brain.Base.Stockpile[stickItemName] != 6 {
		t.Fatalf("expected stockpile unchanged when shelter requirement fails")
	}
}
