package main

import (
	"testing"
	"time"
)

func TestClaimBaseTaskSelectsLocationNearWaterFoodAndLumber(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)

	nearBase := Location{X: 6, Y: 6}
	farBase := Location{X: 15, Y: 15}
	waterLocation := Location{X: 4, Y: 4}
	foodLocation := Location{X: 5, Y: 4}
	lumberLocation := Location{X: 6, Y: 4}

	world.SetTileType(waterLocation.X, waterLocation.Y, Water)

	appleTree := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Fruit:         []Fruit{CreateNewFruit("Apple", 0, true, 20)},
	}
	if err := world.AddPlant(foodLocation.X, foodLocation.Y, appleTree); err != nil {
		t.Fatalf("failed to add food plant: %v", err)
	}

	oakTree := &Plant{Name: OakTree, IsAlive: true}
	if err := world.AddPlant(lumberLocation.X, lumberLocation.Y, oakTree); err != nil {
		t.Fatalf("failed to add lumber tree: %v", err)
	}

	now := time.Now().UnixMilli()
	brain.AddLocationToCognitiveMap(nearBase, CognitiveMapTile{TileType: Grass, LastSeenUnixMs: now})
	brain.AddLocationToCognitiveMap(farBase, CognitiveMapTile{TileType: Grass, LastSeenUnixMs: now})
	brain.AddLocationToCognitiveMap(waterLocation, CognitiveMapTile{TileType: Water, LastSeenUnixMs: now})
	brain.AddLocationToCognitiveMap(foodLocation, CognitiveMapTile{
		TileType: Grass,
		Plant: CognitiveMapPlant{
			Name:          AppleTree,
			IsAlive:       true,
			ProducesFruit: true,
			FruitCount:    1,
			HasRipeFruit:  true,
		},
		LastSeenUnixMs: now,
	})
	brain.AddLocationToCognitiveMap(lumberLocation, CognitiveMapTile{
		TileType: Grass,
		Plant: CognitiveMapPlant{
			Name:       OakTree,
			IsAlive:    true,
			PlantStage: Vegetative,
		},
		LastSeenUnixMs: now,
	})

	brain.ClaimBaseTask(TargetedAction{Action: ClaimBase})

	if !brain.Base.Claimed {
		t.Fatalf("expected base to be claimed")
	}
	if brain.Base.Location != nearBase {
		t.Fatalf("expected near base %+v, got %+v", nearBase, brain.Base.Location)
	}
}

func TestMakeShelterTaskBuildsShelterFromHarvestedSticksWithoutStoneAxe(t *testing.T) {
	world := NewWorld(30, 30)
	brain := newVisionTestBrain(t, world, 2, 2)

	if brain.FindInOwnedItems("Stone Axe") != nil {
		t.Fatalf("expected test subject to start without stone axe")
	}

	baseLocation := Location{X: 6, Y: 6}
	waterLocation := Location{X: 4, Y: 4}
	foodLocation := Location{X: 5, Y: 4}
	lumberLocation := Location{X: 8, Y: 6}

	world.SetTileType(waterLocation.X, waterLocation.Y, Water)

	appleTree := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Fruit:         []Fruit{CreateNewFruit("Apple", 0, true, 20)},
	}
	if err := world.AddPlant(foodLocation.X, foodLocation.Y, appleTree); err != nil {
		t.Fatalf("failed to add food plant: %v", err)
	}

	oakTree := &Plant{Name: OakTree, IsAlive: true}
	if err := world.AddPlant(lumberLocation.X, lumberLocation.Y, oakTree); err != nil {
		t.Fatalf("failed to add lumber tree: %v", err)
	}

	now := time.Now().UnixMilli()
	brain.AddLocationToCognitiveMap(baseLocation, CognitiveMapTile{TileType: Grass, LastSeenUnixMs: now})
	brain.AddLocationToCognitiveMap(waterLocation, CognitiveMapTile{TileType: Water, LastSeenUnixMs: now})
	brain.AddLocationToCognitiveMap(foodLocation, CognitiveMapTile{
		TileType: Grass,
		Plant: CognitiveMapPlant{
			Name:          AppleTree,
			IsAlive:       true,
			ProducesFruit: true,
			FruitCount:    1,
			HasRipeFruit:  true,
		},
		LastSeenUnixMs: now,
	})
	brain.AddLocationToCognitiveMap(lumberLocation, CognitiveMapTile{
		TileType: Grass,
		Plant: CognitiveMapPlant{
			Name:    OakTree,
			IsAlive: true,
		},
		LastSeenUnixMs: now,
	})

	brain.ClaimBaseTask(TargetedAction{Action: ClaimBase})
	if !brain.Base.Claimed {
		t.Fatalf("expected base to be claimed before building shelter")
	}

	for i := 0; i < stickBundlesNeededForShelter; i++ {
		if err := world.MoveEntity(brain.Owner, lumberLocation.X, lumberLocation.Y); err != nil {
			t.Fatalf("failed to move entity to lumber location: %v", err)
		}
		brain.AddMemoryToLongTerm("Found lumber tree", "Lumber", lumberLocation)
		brain.MakeShelterTask(TargetedAction{Action: MakeShelter})
		if brain.countOwnedItemsByName(stickItemName) == 0 {
			t.Fatalf("expected sticks to be carried after harvesting")
		}
		if !hasItemInHands(brain.Owner, stickItemName) {
			t.Fatalf("expected harvested sticks to be carried in hand")
		}

		if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
			t.Fatalf("failed to move entity to base location: %v", err)
		}
		brain.MakeShelterTask(TargetedAction{Action: MakeShelter})

		baseTileAfterDeposit := world.GetTile(brain.Base.Location.X, brain.Base.Location.Y)
		groundStickCount := 0
		for _, item := range baseTileAfterDeposit.Items {
			if item != nil && item.Name == stickItemName {
				groundStickCount++
			}
		}
		expectedGroundSticks := i + 1
		if groundStickCount != expectedGroundSticks {
			t.Fatalf("expected %d sticks on base tile, got %d", expectedGroundSticks, groundStickCount)
		}
	}

	if brain.Base.Stockpile[stickItemName] < stickBundlesNeededForShelter {
		t.Fatalf("expected enough sticks in stockpile, got %d", brain.Base.Stockpile[stickItemName])
	}

	// One more tick at base should consume sticks and construct shelter.
	brain.MakeShelterTask(TargetedAction{Action: MakeShelter})

	baseTile := world.GetTile(brain.Base.Location.X, brain.Base.Location.Y)
	if baseTile.Shelter == nil {
		t.Fatalf("expected shelter to be built at base location")
	}
	if !brain.PhysiologicalNeeds.HasShelter {
		t.Fatalf("expected physiological shelter need to be satisfied")
	}
	if brain.Base.Stockpile[stickItemName] != 0 {
		t.Fatalf("expected stockpile sticks to be consumed, got %d", brain.Base.Stockpile[stickItemName])
	}

	groundStickCount := 0
	for _, item := range baseTile.Items {
		if item != nil && item.Name == stickItemName {
			groundStickCount++
		}
	}
	if groundStickCount != 0 {
		t.Fatalf("expected consumed sticks to be removed from ground, got %d", groundStickCount)
	}
	if brain.CurrentTask.Action != None {
		t.Fatalf("expected current task to fall back to idle, got %q", brain.CurrentTask.Action)
	}
	if brain.MotorCortexCurrentTask.ActionReason != "Idle" &&
		brain.MotorCortexCurrentTask.ActionReason != "Searching for Nearby resources" {
		t.Fatalf("expected fallback motor action to be idle/explore, got %q", brain.MotorCortexCurrentTask.ActionReason)
	}
}

func TestEnsureCarryCapacityForDropsNonCriticalItemFirst(t *testing.T) {
	world := NewWorld(12, 12)
	brain := newVisionTestBrain(t, world, 2, 2)

	stoneAxe := CreateNewItem("Stone Axe")
	woodenStaff := CreateNewItem("Wooden Staff")
	if stoneAxe == nil || woodenStaff == nil {
		t.Fatalf("expected test items to be created")
	}

	brain.Owner.GrabWithRightHand(stoneAxe) // critical
	brain.Owner.GrabWithLeftHand(woodenStaff)

	if !brain.ensureCarryCapacityFor(stickItemName) {
		t.Fatalf("expected capacity helper to free one hand")
	}
	if !hasItemInHands(brain.Owner, "Stone Axe") {
		t.Fatalf("expected critical item to be kept in hand")
	}
	if hasItemInHands(brain.Owner, "Wooden Staff") {
		t.Fatalf("expected non-critical item to be dropped")
	}
	if got := countItemsOnTileByName(world.GetTile(2, 2), "Wooden Staff"); got != 1 {
		t.Fatalf("expected dropped wooden staff on ground, got %d", got)
	}
}

func TestGetFoodForStorageCollectsAndDepositsFoodAtBase(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)
	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 6, Y: 6}

	foodLocation := Location{X: 8, Y: 6}
	appleTree := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Fruit:         []Fruit{CreateNewFruit("Apple", 0, true, 20)},
	}
	if err := world.AddPlant(foodLocation.X, foodLocation.Y, appleTree); err != nil {
		t.Fatalf("failed to add apple tree: %v", err)
	}
	brain.AddMemoryToLongTerm("Found food supply", "Food", foodLocation)

	if err := world.MoveEntity(brain.Owner, foodLocation.X, foodLocation.Y); err != nil {
		t.Fatalf("failed to move entity to food source: %v", err)
	}
	brain.GetFoodForStorage(TargetedAction{Action: HaveFood})
	if brain.countOwnedItemsByName(foodCarryItemName) != 1 {
		t.Fatalf("expected one carried food item after collection")
	}

	if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
		t.Fatalf("failed to move entity to base: %v", err)
	}
	brain.GetFoodForStorage(TargetedAction{Action: HaveFood})

	if brain.Base.Stockpile[foodStockpileKey] != 1 {
		t.Fatalf("expected base food stockpile to be 1, got %d", brain.Base.Stockpile[foodStockpileKey])
	}
	baseTile := world.GetTile(brain.Base.Location.X, brain.Base.Location.Y)
	if got := countItemsOnTileByName(baseTile, foodCarryItemName); got != 1 {
		t.Fatalf("expected one deposited food item on base tile, got %d", got)
	}
}

func TestMaintainStickStockTaskCapsAtFive(t *testing.T) {
	world := NewWorld(24, 24)
	brain := newVisionTestBrain(t, world, 3, 3)
	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 7, Y: 7}
	brain.PhysiologicalNeeds.HasShelter = true

	lumberLocation := Location{X: 9, Y: 7}
	oakTree := &Plant{Name: OakTree, IsAlive: true}
	if err := world.AddPlant(lumberLocation.X, lumberLocation.Y, oakTree); err != nil {
		t.Fatalf("failed to add oak tree: %v", err)
	}
	brain.AddMemoryToLongTerm("Found lumber tree", "Lumber", lumberLocation)

	for i := 0; i < 10; i++ {
		if err := world.MoveEntity(brain.Owner, lumberLocation.X, lumberLocation.Y); err != nil {
			t.Fatalf("failed to move entity to lumber location: %v", err)
		}
		brain.MaintainStickStockTask(TargetedAction{Action: HaveLumber})

		if brain.countOwnedItemsByName(stickItemName) > 0 {
			if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
				t.Fatalf("failed to move entity to base location: %v", err)
			}
			brain.MaintainStickStockTask(TargetedAction{Action: HaveLumber})
		}
	}

	if brain.Base.Stockpile[stickItemName] != baseStickMaintenanceTarget {
		t.Fatalf("expected stick stockpile cap %d, got %d", baseStickMaintenanceTarget, brain.Base.Stockpile[stickItemName])
	}
}

func TestCollectStoneForBaseTaskCollectsAndDepositsStoneAtBase(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)
	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 6, Y: 6}

	stoneLocation := Location{X: 8, Y: 6}
	if err := world.AddItem(stoneLocation.X, stoneLocation.Y, CreateNewItem(stoneItemName)); err != nil {
		t.Fatalf("failed to add stone item: %v", err)
	}
	brain.AddLocationToCognitiveMap(stoneLocation, CognitiveMapTile{
		TileType:       Grass,
		Items:          map[string]int{stoneItemName: 1},
		LastSeenUnixMs: time.Now().UnixMilli(),
	})
	brain.AddMemoryToLongTerm("Found stone supply", "Stone", stoneLocation)

	if err := world.MoveEntity(brain.Owner, stoneLocation.X, stoneLocation.Y); err != nil {
		t.Fatalf("failed to move entity to stone source: %v", err)
	}
	brain.CollectStoneForBaseTask(TargetedAction{Action: HaveStone})
	if brain.countOwnedItemsByName(stoneItemName) != 1 {
		t.Fatalf("expected one carried stone after collection")
	}
	knownStoneTile := brain.GetLocationFromCognitiveMap(stoneLocation)
	if knownStoneTile.Items[stoneItemName] != 0 {
		t.Fatalf("expected picked stone to be removed from cognitive map, got count %d", knownStoneTile.Items[stoneItemName])
	}

	if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
		t.Fatalf("failed to move entity to base: %v", err)
	}
	brain.CollectStoneForBaseTask(TargetedAction{Action: HaveStone})

	if brain.Base.Stockpile[stoneStockpileKey] != 1 {
		t.Fatalf("expected base stone stockpile to be 1, got %d", brain.Base.Stockpile[stoneStockpileKey])
	}
	baseTile := world.GetTile(brain.Base.Location.X, brain.Base.Location.Y)
	if got := countItemsOnTileByName(baseTile, stoneItemName); got != 1 {
		t.Fatalf("expected one deposited stone on base tile, got %d", got)
	}
	knownBaseTile := brain.GetLocationFromCognitiveMap(brain.Base.Location)
	if knownBaseTile.Items[stoneItemName] != 1 {
		t.Fatalf("expected base cognitive-map stone count to be 1 after deposit, got %d", knownBaseTile.Items[stoneItemName])
	}
}

func TestCollectStoneForBaseTaskSkipsBaseStoneAsSource(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)
	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 6, Y: 6}

	baseStoneLocation := brain.Base.Location
	remoteStoneLocation := Location{X: 9, Y: 6}
	if err := world.AddItem(baseStoneLocation.X, baseStoneLocation.Y, CreateNewItem(stoneItemName)); err != nil {
		t.Fatalf("failed to add base stone item: %v", err)
	}
	if err := world.AddItem(remoteStoneLocation.X, remoteStoneLocation.Y, CreateNewItem(stoneItemName)); err != nil {
		t.Fatalf("failed to add remote stone item: %v", err)
	}

	brain.AddMemoryToLongTerm("Found stone supply", "Stone", baseStoneLocation)
	brain.AddMemoryToLongTerm("Found stone supply", "Stone", remoteStoneLocation)
	if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
		t.Fatalf("failed to move entity to base: %v", err)
	}

	brain.CollectStoneForBaseTask(TargetedAction{Action: HaveStone})

	if brain.MotorCortexCurrentTask.ActionReason != "Collect stone for base" {
		t.Fatalf("expected collect-stone walk task, got %q", brain.MotorCortexCurrentTask.ActionReason)
	}
	if brain.MotorCortexCurrentTask.TargetLocation != remoteStoneLocation {
		t.Fatalf("expected remote stone target %+v, got %+v", remoteStoneLocation, brain.MotorCortexCurrentTask.TargetLocation)
	}
}

func TestMaintainStoneStockTaskCapsAtFive(t *testing.T) {
	world := NewWorld(24, 24)
	brain := newVisionTestBrain(t, world, 3, 3)
	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 7, Y: 7}
	brain.PhysiologicalNeeds.HasShelter = true

	stoneLocation := Location{X: 9, Y: 7}
	for i := 0; i < 20; i++ {
		if err := world.AddItem(stoneLocation.X, stoneLocation.Y, CreateNewItem(stoneItemName)); err != nil {
			t.Fatalf("failed to add stone item: %v", err)
		}
	}
	brain.AddMemoryToLongTerm("Found stone supply", "Stone", stoneLocation)

	for i := 0; i < 12; i++ {
		if err := world.MoveEntity(brain.Owner, stoneLocation.X, stoneLocation.Y); err != nil {
			t.Fatalf("failed to move entity to stone location: %v", err)
		}
		brain.MaintainStoneStockTask(TargetedAction{Action: HaveStone})

		if brain.countOwnedItemsByName(stoneItemName) > 0 {
			if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
				t.Fatalf("failed to move entity to base location: %v", err)
			}
			brain.MaintainStoneStockTask(TargetedAction{Action: HaveStone})
		}
	}

	if brain.Base.Stockpile[stoneStockpileKey] != baseStoneStockTarget {
		t.Fatalf("expected stone stockpile cap %d, got %d", baseStoneStockTarget, brain.Base.Stockpile[stoneStockpileKey])
	}
}

func TestCollectGrassForBaseTaskHarvestsHighGrassAndDepositsAtBase(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 2, 2)
	brain.ensureBaseStockpile()
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 6, Y: 6}
	brain.PhysiologicalNeeds.HasShelter = true

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
	brain.CollectGrassForBaseTask(TargetedAction{Action: HaveGrass})
	if brain.countOwnedItemsByName(grassItemName) != 1 {
		t.Fatalf("expected one carried grass item after harvest")
	}
	if world.GetPlants(grassLocation.X, grassLocation.Y) != nil {
		t.Fatalf("expected high grass plant to be harvested and removed from tile")
	}
	knownGrassTile := brain.GetLocationFromCognitiveMap(grassLocation)
	if knownGrassTile.Plant.IsAlive {
		t.Fatalf("expected harvested grass tile to be synced as plant-free in cognitive map")
	}

	if err := world.MoveEntity(brain.Owner, brain.Base.Location.X, brain.Base.Location.Y); err != nil {
		t.Fatalf("failed to move entity to base: %v", err)
	}
	brain.CollectGrassForBaseTask(TargetedAction{Action: HaveGrass})

	if brain.Base.Stockpile[grassStockpileKey] != 1 {
		t.Fatalf("expected base grass stockpile to be 1, got %d", brain.Base.Stockpile[grassStockpileKey])
	}
	baseTile := world.GetTile(brain.Base.Location.X, brain.Base.Location.Y)
	if got := countItemsOnTileByName(baseTile, grassItemName); got != 1 {
		t.Fatalf("expected one deposited grass item on base tile, got %d", got)
	}
}

func hasItemInHands(entity *Entity, itemName string) bool {
	if entity == nil || entity.Body == nil {
		return false
	}

	if entity.Body.RightArm != nil && entity.Body.RightArm.Hand != nil {
		for _, item := range entity.Body.RightArm.Hand.Items {
			if item != nil && item.Name == itemName {
				return true
			}
		}
	}

	if entity.Body.LeftArm != nil && entity.Body.LeftArm.Hand != nil {
		for _, item := range entity.Body.LeftArm.Hand.Items {
			if item != nil && item.Name == itemName {
				return true
			}
		}
	}

	return false
}

func countItemsOnTileByName(tile Tile, itemName string) int {
	count := 0
	for _, item := range tile.Items {
		if item != nil && item.Name == itemName {
			count++
		}
	}
	return count
}
