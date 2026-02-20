package main

import (
	"testing"
	"time"
)

func TestRankTasksPrefersStockpileResourcesAfterShelter(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)

	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 6, Y: 6}
	brain.Base.Stockpile[stickItemName] = 0
	brain.Base.Stockpile[stoneStockpileKey] = 0
	brain.Base.Stockpile[grassStockpileKey] = 0

	brain.PhysiologicalNeeds.HasShelter = true
	brain.PhysiologicalNeeds.WayOfGettingWater = true
	brain.PhysiologicalNeeds.WayOfGettingFood = true
	brain.PhysiologicalNeeds.IsSufficientlyWarm = true
	brain.PhysiologicalNeeds.IsInSafeArea = true
	brain.PhysiologicalNeeds.IsInPain = false
	brain.PhysiologicalNeeds.NeedToExcrete = false
	brain.PhysiologicalNeeds.Rested = 100
	brain.Owner.Body.Blood.Water = 100
	brain.Owner.Body.Blood.Glucose = 100

	brain.TranslateWantToTaskList()
	ranked := brain.RankTasks()
	if ranked.Action != StockpileResources {
		t.Fatalf("expected highest-priority task %q, got %q", StockpileResources, ranked.Action)
	}
}

func TestStockpileResourcesTaskHarvestsGrassWhenStickAndStoneAreFull(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)

	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 6, Y: 6}
	brain.PhysiologicalNeeds.HasShelter = true
	brain.Base.Stockpile[stickItemName] = baseStickMaintenanceTarget
	brain.Base.Stockpile[stoneStockpileKey] = baseStoneStockTarget
	brain.Base.Stockpile[grassStockpileKey] = 0

	grassLocation := Location{X: 8, Y: 6}
	highGrass := &Plant{Name: HighGrass, IsAlive: true}
	if err := world.AddPlant(grassLocation.X, grassLocation.Y, highGrass); err != nil {
		t.Fatalf("failed to add high grass plant: %v", err)
	}
	brain.AddLocationToCognitiveMap(grassLocation, CognitiveMapTile{
		TileType: Grass,
		Plant: CognitiveMapPlant{
			Name:    HighGrass,
			IsAlive: true,
		},
		LastSeenUnixMs: time.Now().UnixMilli(),
	})
	brain.AddMemoryToLongTerm("Found grass supply", "Grass", grassLocation)

	if err := world.MoveEntity(brain.Owner, grassLocation.X, grassLocation.Y); err != nil {
		t.Fatalf("failed to move entity to grass source: %v", err)
	}

	brain.StockpileResourcesTask(TargetedAction{Action: StockpileResources})
	if brain.countOwnedItemsByName(grassItemName) != 1 {
		t.Fatalf("expected stockpile task to harvest one grass item")
	}
}
