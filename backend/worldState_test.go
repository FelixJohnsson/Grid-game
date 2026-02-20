package main

import (
	"errors"
	"testing"
)

func newTestEntity(fullName string, x, y int) *Entity {
	return &Entity{
		FullName: fullName,
		Location: Location{X: x, Y: y},
	}
}

func TestGetPersonByFullNameHandlesEmptyTiles(t *testing.T) {
	world := NewWorld(3, 3)

	if got := world.GetPersonByFullName("missing"); got != nil {
		t.Fatalf("expected nil for missing person, got %+v", got)
	}

	person := newTestEntity("Alice", 1, 1)
	if err := world.AddEntity(1, 1, person); err != nil {
		t.Fatalf("failed to add person: %v", err)
	}

	got := world.GetPersonByFullName("Alice")
	if got != person {
		t.Fatalf("expected %p, got %p", person, got)
	}
}

func TestGetAllPersonsFiltersNilTiles(t *testing.T) {
	world := NewWorld(3, 3)

	personA := newTestEntity("A", 0, 0)
	personB := newTestEntity("B", 2, 2)
	if err := world.AddEntity(0, 0, personA); err != nil {
		t.Fatalf("failed to add first person: %v", err)
	}
	if err := world.AddEntity(2, 2, personB); err != nil {
		t.Fatalf("failed to add second person: %v", err)
	}

	all := world.GetAllPersons()
	if len(all) != 2 {
		t.Fatalf("expected 2 persons, got %d", len(all))
	}
}

func TestMoveEntityValidatesAndMoves(t *testing.T) {
	world := NewWorld(3, 3)
	mover := newTestEntity("Mover", 1, 1)

	if err := world.AddEntity(1, 1, mover); err != nil {
		t.Fatalf("failed to add mover: %v", err)
	}

	if err := world.MoveEntity(nil, 1, 2); !errors.Is(err, ErrNilEntity) {
		t.Fatalf("expected ErrNilEntity, got %v", err)
	}

	if err := world.MoveEntity(mover, -1, 2); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("expected ErrOutOfBounds, got %v", err)
	}

	world.SetTileType(2, 1, Mountain)
	if err := world.MoveEntity(mover, 2, 1); !errors.Is(err, ErrUnwalkableTile) {
		t.Fatalf("expected ErrUnwalkableTile, got %v", err)
	}
	world.SetTileType(2, 1, Grass)

	blocker := newTestEntity("Blocker", 2, 1)
	if err := world.AddEntity(2, 1, blocker); err != nil {
		t.Fatalf("failed to add blocker: %v", err)
	}
	if err := world.MoveEntity(mover, 2, 1); !errors.Is(err, ErrTileOccupied) {
		t.Fatalf("expected ErrTileOccupied, got %v", err)
	}

	if err := world.MoveEntity(blocker, 2, 2); err != nil {
		t.Fatalf("failed to move blocker away: %v", err)
	}
	if err := world.MoveEntity(mover, 2, 1); err != nil {
		t.Fatalf("failed to move mover: %v", err)
	}

	if world.GetPersons(1, 1) != nil {
		t.Fatalf("expected old tile to be empty")
	}
	if world.GetPersons(2, 1) != mover {
		t.Fatalf("expected mover on destination tile")
	}
	if mover.Location.X != 2 || mover.Location.Y != 1 {
		t.Fatalf("expected mover location updated to (2,1), got (%d,%d)", mover.Location.X, mover.Location.Y)
	}
}

func TestItemMutatorsReturnErrorsAndUpdateTile(t *testing.T) {
	world := NewWorld(2, 2)
	item := &Item{Name: "Rock"}

	if err := world.AddItem(0, 0, nil); !errors.Is(err, ErrNilItem) {
		t.Fatalf("expected ErrNilItem, got %v", err)
	}
	if err := world.AddItem(-1, 0, item); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("expected ErrOutOfBounds, got %v", err)
	}
	if err := world.AddItem(0, 0, item); err != nil {
		t.Fatalf("failed to add item: %v", err)
	}

	if item.Location.X != 0 || item.Location.Y != 0 {
		t.Fatalf("expected item location updated to (0,0), got (%d,%d)", item.Location.X, item.Location.Y)
	}

	missing := &Item{Name: "Missing"}
	if _, err := world.RemoveItem(missing, 0, 0); !errors.Is(err, ErrItemNotFoundOnTile) {
		t.Fatalf("expected ErrItemNotFoundOnTile, got %v", err)
	}

	if err := world.DestroyItem(nil); !errors.Is(err, ErrNilItem) {
		t.Fatalf("expected ErrNilItem from DestroyItem, got %v", err)
	}
	if err := world.DestroyItem(item); err != nil {
		t.Fatalf("failed to destroy item: %v", err)
	}
	if got := len(world.GetItems(0, 0)); got != 0 {
		t.Fatalf("expected no items left, got %d", got)
	}
}

func TestPlantMutatorsReturnErrorsAndUpdateTile(t *testing.T) {
	world := NewWorld(2, 2)
	plant := &Plant{Name: AppleTree, Location: Location{X: 0, Y: 0}}

	if _, err := world.RemovePlant(nil); !errors.Is(err, ErrNilPlant) {
		t.Fatalf("expected ErrNilPlant, got %v", err)
	}
	if _, err := world.RemovePlant(plant); !errors.Is(err, ErrPlantNotFoundOnTile) {
		t.Fatalf("expected ErrPlantNotFoundOnTile, got %v", err)
	}
	if err := world.AddPlant(-1, 0, plant); !errors.Is(err, ErrOutOfBounds) {
		t.Fatalf("expected ErrOutOfBounds from AddPlant, got %v", err)
	}
	if err := world.AddPlant(0, 0, plant); err != nil {
		t.Fatalf("failed to add plant: %v", err)
	}

	tile, err := world.RemovePlant(plant)
	if err != nil {
		t.Fatalf("failed to remove plant: %v", err)
	}
	if tile.Plant != nil {
		t.Fatalf("expected plant to be removed from tile")
	}
	if world.GetPlants(0, 0) != nil {
		t.Fatalf("expected no plant on tile after removal")
	}
}

func TestOutOfBoundsReadGuards(t *testing.T) {
	world := NewWorld(2, 2)

	tile := world.GetTile(-1, 0)
	if tile.Type != Mountain {
		t.Fatalf("expected out-of-bounds tile to be Mountain, got %v", tile.Type)
	}
	if world.IsTileEmpty(-1, 0) {
		t.Fatalf("expected out-of-bounds tile to be non-empty (guarded false)")
	}
	if world.CanWalk(99, 99) {
		t.Fatalf("expected out-of-bounds tile to be non-walkable")
	}
}

func TestConsumeRipeFruitAtConsumesOneFruitAndKeepsPlant(t *testing.T) {
	world := NewWorld(3, 3)
	plant := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Location:      Location{X: 1, Y: 1},
		Fruit: []Fruit{
			CreateNewFruit("Apple", 0, false, 20),
			CreateNewFruit("Apple", 0, true, 20),
		},
	}
	if err := world.AddPlant(1, 1, plant); err != nil {
		t.Fatalf("failed to add plant: %v", err)
	}

	fruit, err := world.ConsumeRipeFruitAt(1, 1)
	if err != nil {
		t.Fatalf("expected ripe fruit consumption to succeed, got %v", err)
	}
	if !fruit.IsRipe {
		t.Fatalf("expected consumed fruit to be ripe")
	}

	remaining := world.GetPlants(1, 1)
	if remaining == nil {
		t.Fatalf("expected plant to remain on tile after consuming fruit")
	}
	if got := len(remaining.Fruit); got != 1 {
		t.Fatalf("expected 1 fruit remaining on plant, got %d", got)
	}
}

func TestTickPlantLifecycleSpawnsAndRipensAppleFruit(t *testing.T) {
	world := NewWorld(3, 3)
	plant := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Age:           appleTreeSpawnEveryTicks - 1,
		Location:      Location{X: 1, Y: 1},
		Fruit:         []Fruit{},
	}
	if err := world.AddPlant(1, 1, plant); err != nil {
		t.Fatalf("failed to add plant: %v", err)
	}

	world.TickPlantLifecycle() // age reaches spawn boundary, should spawn one unripe fruit
	tickedPlant := world.GetPlants(1, 1)
	if tickedPlant == nil {
		t.Fatalf("expected plant at (1,1)")
	}
	if got := len(tickedPlant.Fruit); got != 1 {
		t.Fatalf("expected 1 spawned fruit, got %d", got)
	}
	if tickedPlant.Fruit[0].IsRipe {
		t.Fatalf("expected spawned fruit to start unripe")
	}

	for i := 0; i < appleFruitRipenAfterTicks; i++ {
		world.TickPlantLifecycle()
	}

	tickedPlant = world.GetPlants(1, 1)
	if tickedPlant == nil || len(tickedPlant.Fruit) == 0 {
		t.Fatalf("expected fruit to remain on plant")
	}
	if !tickedPlant.Fruit[0].IsRipe {
		t.Fatalf("expected fruit to ripen by age threshold")
	}
}

func TestIntentEngineRoutesWorldMutationsThroughIntentLoop(t *testing.T) {
	world := NewWorld(4, 4)
	world.EnableIntentEngine()

	mover := newTestEntity("IntentMover", 1, 1)
	if err := world.AddEntity(1, 1, mover); err != nil {
		t.Fatalf("failed to add mover: %v", err)
	}

	if err := world.MoveEntity(mover, 1, 2); err != nil {
		t.Fatalf("failed to move with intent engine: %v", err)
	}
	if mover.Location.X != 1 || mover.Location.Y != 2 {
		t.Fatalf("expected mover location to be (1,2), got (%d,%d)", mover.Location.X, mover.Location.Y)
	}

	item := &Item{Name: "Stone"}
	if err := world.AddItem(1, 2, item); err != nil {
		t.Fatalf("failed to add item with intent engine: %v", err)
	}
	if got := len(world.GetItems(1, 2)); got != 1 {
		t.Fatalf("expected 1 item on tile, got %d", got)
	}
	if err := world.DestroyItem(item); err != nil {
		t.Fatalf("failed to destroy item with intent engine: %v", err)
	}
	if got := len(world.GetItems(1, 2)); got != 0 {
		t.Fatalf("expected tile to be empty after destroy, got %d", got)
	}

	plant := &Plant{
		Name:          AppleTree,
		IsAlive:       true,
		ProducesFruit: true,
		Location:      Location{X: 2, Y: 2},
		Fruit: []Fruit{
			CreateNewFruit("Apple", 0, true, 20),
		},
	}
	if err := world.AddPlant(2, 2, plant); err != nil {
		t.Fatalf("failed to add plant: %v", err)
	}
	if _, err := world.ConsumeRipeFruitAt(2, 2); err != nil {
		t.Fatalf("failed to consume ripe fruit with intent engine: %v", err)
	}
	if got := len(world.GetPlants(2, 2).Fruit); got != 0 {
		t.Fatalf("expected plant fruit to be consumed, got %d remaining", got)
	}
}
