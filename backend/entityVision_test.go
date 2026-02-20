package main

import (
	"testing"
	"time"
)

func newVisionTestBrain(t *testing.T, world *World, x, y int) *Brain {
	t.Helper()

	person := NewPersonEntity(world, x, y, Human)
	if person == nil {
		t.Fatalf("expected person to be created")
	}
	if err := world.AddEntity(x, y, person); err != nil {
		t.Fatalf("failed to add person to world: %v", err)
	}
	return person.Brain
}

func TestGetWaterSupplyInMemoryPrunesStaleLocations(t *testing.T) {
	world := NewWorld(8, 8)
	brain := newVisionTestBrain(t, world, 1, 1)

	stale := Location{X: 6, Y: 6} // default grass tile
	brain.AddMemoryToLongTerm("Found water supply", "Water", stale)

	memory := brain.GetWaterSupplyInMemory()
	if memory.Event != "" {
		t.Fatalf("expected no valid water memory, got %+v", memory)
	}
	if len(brain.Memories.LongTermMemory) != 0 {
		t.Fatalf("expected stale water memory to be removed, got %d entries", len(brain.Memories.LongTermMemory))
	}
}

func TestDrinkWaterTaskUsesKnownWaterMemory(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 1, 1)
	brain.Owner.VisionRange = 2

	knownWater := Location{X: 10, Y: 10}
	world.SetTileType(knownWater.X, knownWater.Y, Water)
	brain.AddMemoryToLongTerm("Found water supply", "Water", knownWater)

	brain.DrinkWaterTask(TargetedAction{})

	if brain.MotorCortexCurrentTask.ActionReason != "Drink water" {
		t.Fatalf("expected motor cortex to pursue drinking, got %q", brain.MotorCortexCurrentTask.ActionReason)
	}
	if brain.MotorCortexCurrentTask.TargetLocation != knownWater {
		t.Fatalf("expected target water location %+v, got %+v", knownWater, brain.MotorCortexCurrentTask.TargetLocation)
	}
}

func TestDrinkWaterTaskFallsBackToSecondClosestWhenClosestIsOccupied(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 1, 1)
	_ = newVisionTestBrain(t, world, 3, 1) // blocker entity on nearest water tile

	blockedWater := Location{X: 3, Y: 1}
	fallbackWater := Location{X: 4, Y: 1}
	world.SetTileType(blockedWater.X, blockedWater.Y, Water)
	world.SetTileType(fallbackWater.X, fallbackWater.Y, Water)

	brain.AddMemoryToLongTerm("Found water supply", "Water", blockedWater)
	brain.AddMemoryToLongTerm("Found water supply", "Water", fallbackWater)

	brain.DrinkWaterTask(TargetedAction{})

	if brain.MotorCortexCurrentTask.ActionReason != "Drink water" {
		t.Fatalf("expected motor cortex to pursue drinking, got %q", brain.MotorCortexCurrentTask.ActionReason)
	}
	if brain.MotorCortexCurrentTask.TargetLocation != fallbackWater {
		t.Fatalf("expected fallback water location %+v, got %+v", fallbackWater, brain.MotorCortexCurrentTask.TargetLocation)
	}
}

func TestGetWaterSupplyInMemoryReturnsClosestValidLocation(t *testing.T) {
	world := NewWorld(30, 30)
	brain := newVisionTestBrain(t, world, 1, 1)

	farWater := Location{X: 20, Y: 20}
	nearWater := Location{X: 3, Y: 2}
	world.SetTileType(farWater.X, farWater.Y, Water)
	world.SetTileType(nearWater.X, nearWater.Y, Water)

	// Intentionally store the farther location in short-term to ensure distance wins over recency.
	brain.AddMemoryToShortTerm("Found water supply", "Water", farWater)
	brain.AddMemoryToLongTerm("Found water supply", "Water", nearWater)

	memory := brain.GetWaterSupplyInMemory()
	if memory.Event != "Found water supply" {
		t.Fatalf("expected water memory, got %+v", memory)
	}
	if memory.Location != nearWater {
		t.Fatalf("expected closest water memory %+v, got %+v", nearWater, memory.Location)
	}
}

func TestGetFoodSupplyInMemoryReturnsClosestValidLocation(t *testing.T) {
	world := NewWorld(30, 30)
	brain := newVisionTestBrain(t, world, 1, 1)

	farFood := Location{X: 18, Y: 18}
	nearFood := Location{X: 4, Y: 3}
	farPlant := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Fruit:         []Fruit{CreateNewFruit("Apple", 0, true, 20)},
	}
	nearPlant := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Fruit:         []Fruit{CreateNewFruit("Apple", 0, true, 20)},
	}
	if err := world.AddPlant(farFood.X, farFood.Y, farPlant); err != nil {
		t.Fatalf("failed to add far food plant: %v", err)
	}
	if err := world.AddPlant(nearFood.X, nearFood.Y, nearPlant); err != nil {
		t.Fatalf("failed to add near food plant: %v", err)
	}

	brain.AddMemoryToShortTerm("Found food supply", "Food", farFood)
	brain.AddMemoryToLongTerm("Found food supply", "Food", nearFood)

	memory := brain.GetFoodSupplyInMemory()
	if memory.Event != "Found food supply" {
		t.Fatalf("expected food memory, got %+v", memory)
	}
	if memory.Location != nearFood {
		t.Fatalf("expected closest food memory %+v, got %+v", nearFood, memory.Location)
	}
}

func TestGetFoodSupplyInMemoryPrefersClosestCognitiveMapLocation(t *testing.T) {
	world := NewWorld(30, 30)
	brain := newVisionTestBrain(t, world, 2, 2)

	nearFood := Location{X: 4, Y: 4}
	farFood := Location{X: 20, Y: 20}
	nearPlant := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Fruit:         []Fruit{CreateNewFruit("Apple", 0, true, 20)},
	}
	farPlant := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Fruit:         []Fruit{CreateNewFruit("Apple", 0, true, 20)},
	}
	if err := world.AddPlant(nearFood.X, nearFood.Y, nearPlant); err != nil {
		t.Fatalf("failed to add near food plant: %v", err)
	}
	if err := world.AddPlant(farFood.X, farFood.Y, farPlant); err != nil {
		t.Fatalf("failed to add far food plant: %v", err)
	}

	nowMs := time.Now().UnixMilli()
	brain.AddLocationToCognitiveMap(nearFood, CognitiveMapTile{
		TileType: Grass,
		Plant: CognitiveMapPlant{
			Name:          AppleTree,
			IsAlive:       true,
			ProducesFruit: true,
			FruitCount:    2,
			HasRipeFruit:  true,
		},
		LastSeenUnixMs: nowMs,
	})
	brain.AddLocationToCognitiveMap(farFood, CognitiveMapTile{
		TileType: Grass,
		Plant: CognitiveMapPlant{
			Name:          AppleTree,
			IsAlive:       true,
			ProducesFruit: true,
			FruitCount:    1,
			HasRipeFruit:  true,
		},
		LastSeenUnixMs: nowMs,
	})

	// Keep a farther episodic memory to ensure cognitive-map location is preferred.
	brain.AddMemoryToLongTerm("Found food supply", "Food", farFood)

	memory := brain.GetFoodSupplyInMemory()
	if memory.Event != "Found food supply" {
		t.Fatalf("expected food memory event, got %+v", memory)
	}
	if memory.Location != nearFood {
		t.Fatalf("expected closest cognitive-map food location %+v, got %+v", nearFood, memory.Location)
	}
}

func TestGetWaterSupplyInMemoryFallsBackToEventMemoryWhenMapObservationIsStale(t *testing.T) {
	world := NewWorld(30, 30)
	brain := newVisionTestBrain(t, world, 2, 2)

	staleMapWater := Location{X: 12, Y: 12}
	freshMemoryWater := Location{X: 5, Y: 5}
	world.SetTileType(freshMemoryWater.X, freshMemoryWater.Y, Water)

	brain.AddLocationToCognitiveMap(staleMapWater, CognitiveMapTile{
		TileType:       Water,
		LastSeenUnixMs: time.Now().Add(-cognitiveMapWaterMaxAge - time.Second).UnixMilli(),
	})
	brain.AddMemoryToLongTerm("Found water supply", "Water", freshMemoryWater)

	memory := brain.GetWaterSupplyInMemory()
	if memory.Event != "Found water supply" {
		t.Fatalf("expected water memory event, got %+v", memory)
	}
	if memory.Location != freshMemoryWater {
		t.Fatalf("expected fallback event memory %+v, got %+v", freshMemoryWater, memory.Location)
	}
}

func TestFindLumberTreesDetectsOakAndStoresMemory(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 5, 5)
	brain.Owner.VisionRange = 4

	oak := &Plant{Name: OakTree, IsAlive: true}
	if err := world.AddPlant(7, 5, oak); err != nil {
		t.Fatalf("failed to add oak tree: %v", err)
	}

	found := brain.FindLumberTrees()
	if !found {
		t.Fatalf("expected lumber tree to be found when oak is in vision")
	}

	memory := brain.GetLumberSupplyInMemory()
	if memory.Event != "Found lumber tree" {
		t.Fatalf("expected lumber memory to be recorded, got %+v", memory)
	}
	if memory.Location != (Location{X: 7, Y: 5}) {
		t.Fatalf("expected lumber memory at (7,5), got %+v", memory.Location)
	}
}

func TestGetLumberSupplyInMemoryReturnsClosestValidLocation(t *testing.T) {
	world := NewWorld(30, 30)
	brain := newVisionTestBrain(t, world, 2, 2)

	farLumber := Location{X: 19, Y: 19}
	nearLumber := Location{X: 5, Y: 3}
	farOak := &Plant{Name: OakTree, IsAlive: true}
	nearOak := &Plant{Name: OakTree, IsAlive: true}
	if err := world.AddPlant(farLumber.X, farLumber.Y, farOak); err != nil {
		t.Fatalf("failed to add far oak: %v", err)
	}
	if err := world.AddPlant(nearLumber.X, nearLumber.Y, nearOak); err != nil {
		t.Fatalf("failed to add near oak: %v", err)
	}

	brain.AddMemoryToShortTerm("Found lumber tree", "Lumber", farLumber)
	brain.AddMemoryToLongTerm("Found lumber tree", "Lumber", nearLumber)

	memory := brain.GetLumberSupplyInMemory()
	if memory.Event != "Found lumber tree" {
		t.Fatalf("expected lumber memory, got %+v", memory)
	}
	if memory.Location != nearLumber {
		t.Fatalf("expected closest lumber memory %+v, got %+v", nearLumber, memory.Location)
	}
}

func TestGetLumberTaskUsesKnownLumberMemory(t *testing.T) {
	world := NewWorld(25, 25)
	brain := newVisionTestBrain(t, world, 2, 2)
	brain.Owner.VisionRange = 2

	knownOak := Location{X: 15, Y: 15}
	oak := &Plant{Name: OakTree, IsAlive: true}
	if err := world.AddPlant(knownOak.X, knownOak.Y, oak); err != nil {
		t.Fatalf("failed to add oak tree: %v", err)
	}
	brain.AddMemoryToLongTerm("Found lumber tree", "Lumber", knownOak)

	brain.GetLumberTask()

	if brain.MotorCortexCurrentTask.ActionReason != "Get lumber" {
		t.Fatalf("expected motor cortex to pursue lumber, got %q", brain.MotorCortexCurrentTask.ActionReason)
	}
	if brain.MotorCortexCurrentTask.TargetLocation != knownOak {
		t.Fatalf("expected lumber target %+v, got %+v", knownOak, brain.MotorCortexCurrentTask.TargetLocation)
	}
}

func TestDecidePathToAvoidsOccupiedTiles(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 1, 1)
	_ = newVisionTestBrain(t, world, 2, 1) // block direct route

	path := brain.DecidePathTo(3, 1)
	if len(path) == 0 {
		t.Fatalf("expected a path around occupied tile")
	}

	for _, node := range path {
		if node.X == 2 && node.Y == 1 {
			t.Fatalf("expected path not to include occupied tile (2,1)")
		}
	}
}

func TestFindStoneSupplyStoresClosestStoneMemory(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 5, 5)
	brain.Owner.VisionRange = 4

	farStone := Location{X: 9, Y: 9}
	nearStone := Location{X: 6, Y: 5}
	if err := world.AddItem(farStone.X, farStone.Y, CreateNewItem("Stone")); err != nil {
		t.Fatalf("failed to add far stone item: %v", err)
	}
	if err := world.AddItem(nearStone.X, nearStone.Y, CreateNewItem("Stone")); err != nil {
		t.Fatalf("failed to add near stone item: %v", err)
	}

	found := brain.FindStoneSupply()
	if !found {
		t.Fatalf("expected stone supply to be found in vision")
	}

	memory := brain.GetStoneSupplyInMemory()
	if memory.Event != "Found stone supply" {
		t.Fatalf("expected stone memory event, got %+v", memory)
	}
	if memory.Location != nearStone {
		t.Fatalf("expected closest stone memory %+v, got %+v", nearStone, memory.Location)
	}
}

func TestGetStoneSupplyInMemoryPrunesStaleLocations(t *testing.T) {
	world := NewWorld(12, 12)
	brain := newVisionTestBrain(t, world, 2, 2)

	stale := Location{X: 10, Y: 10}
	brain.AddMemoryToLongTerm("Found stone supply", "Stone", stale)

	memory := brain.GetStoneSupplyInMemory()
	if memory.Event != "" {
		t.Fatalf("expected no valid stone memory, got %+v", memory)
	}
	if len(brain.Memories.LongTermMemory) != 0 {
		t.Fatalf("expected stale stone memory to be removed, got %d entries", len(brain.Memories.LongTermMemory))
	}
}

func TestGetStoneSupplyInMemoryIgnoresStonesOnOwnBaseTile(t *testing.T) {
	world := NewWorld(20, 20)
	brain := newVisionTestBrain(t, world, 6, 5)
	brain.Base.Claimed = true
	brain.Base.Location = Location{X: 6, Y: 6}

	baseStone := brain.Base.Location
	sourceStone := Location{X: 8, Y: 6}
	if err := world.AddItem(baseStone.X, baseStone.Y, CreateNewItem(stoneItemName)); err != nil {
		t.Fatalf("failed to add base stone: %v", err)
	}
	if err := world.AddItem(sourceStone.X, sourceStone.Y, CreateNewItem(stoneItemName)); err != nil {
		t.Fatalf("failed to add source stone: %v", err)
	}

	brain.AddMemoryToLongTerm("Found stone supply", "Stone", baseStone)
	brain.AddMemoryToLongTerm("Found stone supply", "Stone", sourceStone)

	memory := brain.GetStoneSupplyInMemory()
	if memory.Event != "Found stone supply" {
		t.Fatalf("expected stone memory event, got %+v", memory)
	}
	if memory.Location != sourceStone {
		t.Fatalf("expected non-base stone source %+v, got %+v", sourceStone, memory.Location)
	}
}
