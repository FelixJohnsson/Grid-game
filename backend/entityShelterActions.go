package main

import (
	"math"
	"sort"
)

const (
	stickItemName                = "Sticks"
	stoneItemName                = "Stone"
	grassItemName                = "Grass"
	foodCarryItemName            = "Foraged Food"
	foodStockpileKey             = "Food"
	stoneStockpileKey            = "Stone"
	grassStockpileKey            = "Grass"
	stickBundlesNeededForShelter = 6
	baseStickMaintenanceTarget   = 5
	baseFoodStockTarget          = 5
	baseStoneStockTarget         = 5
	baseGrassStockTarget         = 5
)

func (b *Brain) ClaimBaseTask(_ TargetedAction) {
	if b.Base.Claimed {
		return
	}

	baseLocation, ok := b.SelectBaseLocation()
	if ok {
		b.ensureBaseStockpile()
		b.Base.Claimed = true
		b.Base.Location = baseLocation
		b.AddMemoryToLongTerm("Claimed base", "Base", baseLocation)
		return
	}

	// Gather enough world knowledge to claim a meaningful base.
	_ = b.FindWaterSupply()
	_ = b.FindFoodSupply()
	_ = b.FindLumberTrees()
	b.GoSearchFor("Base location")
}

func (b *Brain) MakeShelterTask(_ TargetedAction) {
	if b.PhysiologicalNeeds.HasShelter {
		b.completeShelterTaskWithFallback()
		return
	}

	if !b.Base.Claimed {
		b.ClaimBaseTask(TargetedAction{Action: ClaimBase})
		return
	}

	if b.baseHasShelter() {
		b.PhysiologicalNeeds.HasShelter = true
		return
	}

	b.ensureBaseStockpile()

	if b.Base.Stockpile[stickItemName] >= stickBundlesNeededForShelter {
		b.tryBuildShelterAtBase()
		return
	}

	if b.countOwnedItemsByName(stickItemName) > 0 {
		b.deliverSticksToBase()
		return
	}

	b.harvestSticksForShelter()
}

func (b *Brain) MaintainStickStockTask(action TargetedAction) {
	if !b.PhysiologicalNeeds.HasShelter || !b.Base.Claimed {
		return
	}

	b.ensureBaseStockpile()
	if b.Base.Stockpile[stickItemName] >= baseStickMaintenanceTarget {
		b.RemoveActionFromActionList(action)
		b.ClearCurrentTask()
		b.MotorCortexCurrentTask = MotorCortexAction{
			ActionReason:   "Idle",
			ActionType:     "Idle",
			TargetLocation: b.Owner.Location,
			IsActive:       false,
			Finished:       true,
		}
		return
	}

	if b.countOwnedItemsByName(stickItemName) > 0 {
		b.deliverSticksToBase()
		return
	}

	b.harvestSticksForShelter()
}

func (b *Brain) StockpileResourcesTask(action TargetedAction) {
	if !b.Base.Claimed {
		b.ClaimBaseTask(TargetedAction{Action: ClaimBase})
		return
	}
	if !b.PhysiologicalNeeds.HasShelter {
		return
	}

	b.ensureBaseStockpile()

	// Prefer depositing already-carried resources before starting a new gather loop.
	if b.countOwnedItemsByName(stickItemName) > 0 {
		b.deliverSticksToBase()
		return
	}
	if b.countOwnedItemsByName(stoneItemName) > 0 {
		b.deliverStoneToBase()
		return
	}
	if b.countOwnedItemsByName(grassItemName) > 0 {
		b.deliverGrassToBase()
		return
	}

	if b.Base.Stockpile[stickItemName] < baseStickMaintenanceTarget {
		b.harvestSticksForShelter()
		return
	}
	if b.Base.Stockpile[stoneStockpileKey] < baseStoneStockTarget {
		b.CollectStoneForBaseTask(TargetedAction{Action: HaveStone})
		return
	}
	if b.Base.Stockpile[grassStockpileKey] < baseGrassStockTarget {
		b.CollectGrassForBaseTask(TargetedAction{Action: HaveGrass})
		return
	}

	b.RemoveActionFromActionList(action)
	b.ClearCurrentTask()
	b.MotorCortexCurrentTask = MotorCortexAction{
		ActionReason:   "Idle",
		ActionType:     "Idle",
		TargetLocation: b.Owner.Location,
		IsActive:       false,
		Finished:       true,
	}
}

func (b *Brain) MaintainStoneStockTask(action TargetedAction) {
	if !b.PhysiologicalNeeds.HasShelter || !b.Base.Claimed {
		return
	}

	b.ensureBaseStockpile()
	if b.Base.Stockpile[stoneStockpileKey] >= baseStoneStockTarget {
		b.RemoveActionFromActionList(action)
		b.ClearCurrentTask()
		b.MotorCortexCurrentTask = MotorCortexAction{
			ActionReason:   "Idle",
			ActionType:     "Idle",
			TargetLocation: b.Owner.Location,
			IsActive:       false,
			Finished:       true,
		}
		return
	}

	b.CollectStoneForBaseTask(action)
}

func (b *Brain) MaintainGrassStockTask(action TargetedAction) {
	if !b.PhysiologicalNeeds.HasShelter || !b.Base.Claimed {
		return
	}

	b.ensureBaseStockpile()
	if b.Base.Stockpile[grassStockpileKey] >= baseGrassStockTarget {
		b.RemoveActionFromActionList(action)
		b.ClearCurrentTask()
		b.MotorCortexCurrentTask = MotorCortexAction{
			ActionReason:   "Idle",
			ActionType:     "Idle",
			TargetLocation: b.Owner.Location,
			IsActive:       false,
			Finished:       true,
		}
		return
	}

	b.CollectGrassForBaseTask(action)
}

func (b *Brain) SelectBaseLocation() (Location, bool) {
	knownTiles := b.GetKnownTilesSnapshot()
	if len(knownTiles) == 0 {
		return Location{}, false
	}

	waterLocations := b.getKnownWaterLocations(knownTiles)
	foodLocations := b.getKnownFoodPlantLocations(knownTiles)
	lumberLocations := b.getKnownLumberLocations(knownTiles)

	if len(waterLocations) == 0 || len(foodLocations) == 0 {
		return Location{}, false
	}

	candidates := b.getKnownBaseCandidates(knownTiles)
	if len(candidates) == 0 {
		if b.isValidBaseTile(b.Owner.Location) {
			return b.Owner.Location, true
		}
		return Location{}, false
	}

	bestScore := math.MaxFloat64
	bestLocation := Location{}
	found := false

	for _, candidate := range candidates {
		// Skip clearly unreachable candidates.
		if b.DecidePathTo(candidate.X, candidate.Y) == nil {
			continue
		}

		dWater, _ := minDistance(candidate, waterLocations)
		dFood, _ := minDistance(candidate, foodLocations)
		dLumber, hasLumber := minDistance(candidate, lumberLocations)

		lumberScore := float64(dLumber)
		if !hasLumber {
			// Prefer tiles with known lumber nearby but still allow fallback.
			lumberScore = float64(dFood) + 8
		}

		ownerDistance := b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, candidate)
		score := 0.45*float64(dWater) + 0.35*float64(dFood) + 0.20*lumberScore + 0.05*float64(ownerDistance)

		if !found || score < bestScore {
			bestScore = score
			bestLocation = candidate
			found = true
		}
	}

	if !found && b.isValidBaseTile(b.Owner.Location) {
		return b.Owner.Location, true
	}

	return bestLocation, found
}

func (b *Brain) isValidBaseTile(location Location) bool {
	tile := b.Owner.WorldProvider.GetTile(location.X, location.Y)
	if tile.Type != Grass {
		return false
	}
	if tile.Plant != nil || tile.Shelter != nil {
		return false
	}
	return tile.Entity == nil || tile.Entity.FullName == b.Owner.FullName
}

func (b *Brain) getKnownBaseCandidates(knownTiles map[Location]CognitiveMapTile) []Location {
	candidates := make([]Location, 0)
	for location, tile := range knownTiles {
		if tile.TileType != Grass {
			continue
		}
		if b.isValidBaseTile(location) {
			candidates = append(candidates, location)
		}
	}
	sortLocations(candidates)
	return candidates
}

func (b *Brain) getKnownWaterLocations(knownTiles map[Location]CognitiveMapTile) []Location {
	locations := make([]Location, 0)
	for location, tile := range knownTiles {
		if tile.TileType != Water {
			continue
		}
		if b.hasWaterAtLocation(location) {
			locations = append(locations, location)
		}
	}
	sortLocations(locations)
	return locations
}

func (b *Brain) getKnownFoodPlantLocations(knownTiles map[Location]CognitiveMapTile) []Location {
	locations := make([]Location, 0)
	for location, tile := range knownTiles {
		if !tile.Plant.IsAlive || !tile.Plant.ProducesFruit {
			continue
		}
		if b.hasFoodPlantAtLocation(location) {
			locations = append(locations, location)
		}
	}
	sortLocations(locations)
	return locations
}

func (b *Brain) getKnownLumberLocations(knownTiles map[Location]CognitiveMapTile) []Location {
	locations := make([]Location, 0)
	for location, tile := range knownTiles {
		if !tile.Plant.IsAlive || tile.Plant.Name != OakTree {
			continue
		}
		if b.hasLumberAtLocation(location) {
			locations = append(locations, location)
		}
	}
	sortLocations(locations)
	return locations
}

func (b *Brain) hasFoodPlantAtLocation(location Location) bool {
	plants := b.Owner.WorldProvider.GetPlantsInVision(location.X, location.Y, 0)
	for _, plant := range plants {
		if plant != nil && plant.IsAlive && plant.ProducesFruit {
			return true
		}
	}
	return false
}

func sortLocations(locations []Location) {
	sort.Slice(locations, func(i, j int) bool {
		if locations[i].Y == locations[j].Y {
			return locations[i].X < locations[j].X
		}
		return locations[i].Y < locations[j].Y
	})
}

func minDistance(from Location, to []Location) (int, bool) {
	if len(to) == 0 {
		return 0, false
	}

	best := 0
	found := false
	for _, location := range to {
		distance := abs(from.X-location.X) + abs(from.Y-location.Y)
		if !found || distance < best {
			best = distance
			found = true
		}
	}

	return best, found
}

func (b *Brain) ensureBaseStockpile() {
	if b.Base.Stockpile == nil {
		b.Base.Stockpile = map[string]int{}
	}
}

func (b *Brain) baseHasShelter() bool {
	if !b.Base.Claimed {
		return false
	}
	tile := b.Owner.WorldProvider.GetTile(b.Base.Location.X, b.Base.Location.Y)
	return tile.Shelter != nil
}

func (b *Brain) tryBuildShelterAtBase() {
	if b.Owner.Location != b.Base.Location {
		b.MotorCortexCurrentTask = MotorCortexAction{"Build shelter", "Walk", b.Base.Location, false, false}
		return
	}

	tile := b.Owner.WorldProvider.GetTile(b.Base.Location.X, b.Base.Location.Y)
	if tile.Shelter != nil {
		b.PhysiologicalNeeds.HasShelter = true
		b.completeShelterTaskWithFallback()
		return
	}

	// If this tile is no longer buildable, force a fresh base claim.
	if tile.Plant != nil {
		b.Base.Claimed = false
		return
	}

	shelter := NewShelter(b.Base.Location.X, b.Base.Location.Y, b.Owner)
	if err := b.Owner.WorldProvider.AddShelter(b.Base.Location.X, b.Base.Location.Y, shelter); err != nil {
		return
	}

	b.consumeGroundItemsAtBase(stickItemName, stickBundlesNeededForShelter)

	b.Base.Stockpile[stickItemName] -= stickBundlesNeededForShelter
	if b.Base.Stockpile[stickItemName] < 0 {
		b.Base.Stockpile[stickItemName] = 0
	}

	b.PhysiologicalNeeds.HasShelter = true
	b.AddMemoryToLongTerm("Built shelter", "Shelter", b.Base.Location)
	b.completeShelterTaskWithFallback()
}

func (b *Brain) harvestSticksForShelter() {
	memory := b.GetLumberSupplyInMemory()
	if memory.Event != "Found lumber tree" {
		if !b.FindLumberTrees() {
			b.GoSearchFor("Lumber tree")
		}
		return
	}

	targetLocation := memory.Location
	if b.Owner.Location != targetLocation {
		b.MotorCortexCurrentTask = MotorCortexAction{"Harvest sticks", "Walk", targetLocation, false, false}
		return
	}

	tile := b.Owner.WorldProvider.GetTile(targetLocation.X, targetLocation.Y)
	if tile.Plant == nil || !tile.Plant.IsAlive || tile.Plant.Name != OakTree {
		b.RemoveMemoriesByEventAtLocation("Found lumber tree", targetLocation)
		return
	}
	if !b.hasFreeHand() && b.countOwnedItemsByName(stickItemName) > 0 {
		b.deliverSticksToBase()
		return
	}

	sticks := CreateNewItem(stickItemName)
	if sticks == nil {
		sticks = &Item{
			Name:      stickItemName,
			Sharpness: 1,
			Bluntness: 1,
			Weight:    1,
			Material:  []Material{GetMaterialByName("Sticks")},
			Residues:  make([]Residue, 0),
			Location:  b.Owner.Location,
		}
	}
	sticks.Location = b.Owner.Location
	if !b.carryItemInHand(sticks) {
		return
	}
	b.MotorCortexCurrentTask = MotorCortexAction{"Carry sticks to base", "Walk", b.Base.Location, false, false}
}

func (b *Brain) deliverSticksToBase() {
	if b.Owner.Location != b.Base.Location {
		b.MotorCortexCurrentTask = MotorCortexAction{"Carry sticks to base", "Walk", b.Base.Location, false, false}
		return
	}

	delivered := b.removeOwnedItemsByName(stickItemName)
	if delivered == 0 {
		return
	}

	deliveredToGround := 0
	for i := 0; i < delivered; i++ {
		sticks := CreateNewItem(stickItemName)
		if sticks == nil {
			continue
		}
		if err := b.Owner.WorldProvider.AddItem(b.Base.Location.X, b.Base.Location.Y, sticks); err != nil {
			continue
		}
		deliveredToGround++
	}

	b.ensureBaseStockpile()
	b.Base.Stockpile[stickItemName] += deliveredToGround
}

func (b *Brain) CollectFoodForBaseTask(action TargetedAction) {
	if !b.Base.Claimed {
		b.ClaimBaseTask(TargetedAction{Action: ClaimBase})
		return
	}

	b.ensureBaseStockpile()
	if b.Base.Stockpile[foodStockpileKey] >= baseFoodStockTarget {
		b.RemoveActionFromActionList(action)
		return
	}

	if b.countOwnedItemsByName(foodCarryItemName) > 0 {
		b.deliverFoodToBase()
		return
	}

	memory := b.GetFoodSupplyInMemory()
	if memory.Event != "Found food supply" {
		if !b.FindFoodSupply() {
			b.GoSearchFor("Food supply")
			return
		}
		memory = b.GetFoodSupplyInMemory()
		if memory.Event != "Found food supply" {
			return
		}
	}

	targetLocation := memory.Location
	if b.Owner.Location != targetLocation {
		b.MotorCortexCurrentTask = MotorCortexAction{"Collect food for base", "Walk", targetLocation, false, false}
		return
	}

	if !b.hasFreeHand() && b.countOwnedItemsByName(foodCarryItemName) > 0 {
		b.deliverFoodToBase()
		return
	}
	if !b.ensureCarryCapacityFor(foodCarryItemName) {
		return
	}

	fruit, err := b.Owner.WorldProvider.ConsumeRipeFruitAt(targetLocation.X, targetLocation.Y)
	if err != nil {
		b.RemoveMemoriesByEventAtLocation("Found food supply", targetLocation)
		return
	}

	foodItem := b.createCarriedFoodItem(fruit)
	if !b.carryItemInHand(foodItem) {
		return
	}

	b.MotorCortexCurrentTask = MotorCortexAction{"Carry food to base", "Walk", b.Base.Location, false, false}
}

func (b *Brain) deliverFoodToBase() {
	if b.Owner.Location != b.Base.Location {
		b.MotorCortexCurrentTask = MotorCortexAction{"Carry food to base", "Walk", b.Base.Location, false, false}
		return
	}

	delivered := b.removeOwnedItemsByName(foodCarryItemName)
	if delivered == 0 {
		return
	}

	deliveredToGround := 0
	for i := 0; i < delivered; i++ {
		foodItem := b.createCarriedFoodItem(Fruit{Name: foodCarryItemName})
		if foodItem == nil {
			continue
		}
		if err := b.Owner.WorldProvider.AddItem(b.Base.Location.X, b.Base.Location.Y, foodItem); err != nil {
			continue
		}
		deliveredToGround++
	}

	b.ensureBaseStockpile()
	b.Base.Stockpile[foodStockpileKey] += deliveredToGround
}

func (b *Brain) CollectStoneForBaseTask(action TargetedAction) {
	if !b.Base.Claimed {
		b.ClaimBaseTask(TargetedAction{Action: ClaimBase})
		return
	}

	b.ensureBaseStockpile()
	if b.Base.Stockpile[stoneStockpileKey] >= baseStoneStockTarget {
		b.RemoveActionFromActionList(action)
		return
	}

	if b.countOwnedItemsByName(stoneItemName) > 0 {
		b.deliverStoneToBase()
		return
	}

	memory := b.GetStoneSupplyInMemory()
	if memory.Event != "Found stone supply" {
		if !b.FindStoneSupply() {
			b.GoSearchFor("Stone")
			return
		}
		memory = b.GetStoneSupplyInMemory()
		if memory.Event != "Found stone supply" {
			return
		}
	}

	targetLocation := memory.Location
	if b.Owner.Location != targetLocation {
		b.MotorCortexCurrentTask = MotorCortexAction{"Collect stone for base", "Walk", targetLocation, false, false}
		return
	}

	if !b.ensureCarryCapacityFor(stoneItemName) {
		return
	}

	tile := b.Owner.WorldProvider.GetTile(targetLocation.X, targetLocation.Y)
	var stoneOnGround *Item
	for _, item := range tile.Items {
		if item != nil && item.Name == stoneItemName {
			stoneOnGround = item
			break
		}
	}
	if stoneOnGround == nil {
		b.RemoveMemoriesByEventAtLocation("Found stone supply", targetLocation)
		return
	}

	if err := b.Owner.WorldProvider.DestroyItem(stoneOnGround); err != nil {
		return
	}
	b.RemoveKnownItemFromCognitiveMap(targetLocation, stoneItemName, 1)

	carriedStone := CreateNewItem(stoneItemName)
	if carriedStone == nil {
		carriedStone = &Item{
			Name:      stoneItemName,
			Sharpness: 0,
			Bluntness: 1,
			Weight:    2,
			Material:  []Material{GetMaterialByName("Stone")},
			Residues:  make([]Residue, 0),
			Location:  b.Owner.Location,
		}
	}
	carriedStone.Location = b.Owner.Location
	if !b.carryItemInHand(carriedStone) {
		_ = b.Owner.WorldProvider.AddItem(targetLocation.X, targetLocation.Y, carriedStone)
		return
	}

	b.MotorCortexCurrentTask = MotorCortexAction{"Carry stone to base", "Walk", b.Base.Location, false, false}
}

func (b *Brain) deliverStoneToBase() {
	if b.Owner.Location != b.Base.Location {
		b.MotorCortexCurrentTask = MotorCortexAction{"Carry stone to base", "Walk", b.Base.Location, false, false}
		return
	}

	delivered := b.removeOwnedItemsByName(stoneItemName)
	if delivered == 0 {
		return
	}

	deliveredToGround := 0
	for i := 0; i < delivered; i++ {
		stone := CreateNewItem(stoneItemName)
		if stone == nil {
			continue
		}
		if err := b.Owner.WorldProvider.AddItem(b.Base.Location.X, b.Base.Location.Y, stone); err != nil {
			continue
		}
		deliveredToGround++
	}
	b.SyncKnownTileFromWorld(b.Base.Location)

	b.ensureBaseStockpile()
	b.Base.Stockpile[stoneStockpileKey] += deliveredToGround
}

func (b *Brain) CollectGrassForBaseTask(action TargetedAction) {
	if !b.Base.Claimed {
		b.ClaimBaseTask(TargetedAction{Action: ClaimBase})
		return
	}

	b.ensureBaseStockpile()
	if b.Base.Stockpile[grassStockpileKey] >= baseGrassStockTarget {
		b.RemoveActionFromActionList(action)
		return
	}

	if b.countOwnedItemsByName(grassItemName) > 0 {
		b.deliverGrassToBase()
		return
	}

	memory := b.GetGrassSupplyInMemory()
	if memory.Event != "Found grass supply" {
		if !b.FindGrassSupply() {
			b.GoSearchFor("High grass")
			return
		}
		memory = b.GetGrassSupplyInMemory()
		if memory.Event != "Found grass supply" {
			return
		}
	}

	targetLocation := memory.Location
	if b.Owner.Location != targetLocation {
		b.MotorCortexCurrentTask = MotorCortexAction{"Collect grass for base", "Walk", targetLocation, false, false}
		return
	}

	if !b.ensureCarryCapacityFor(grassItemName) {
		return
	}

	tile := b.Owner.WorldProvider.GetTile(targetLocation.X, targetLocation.Y)
	if tile.Plant == nil || !tile.Plant.IsAlive || tile.Plant.Name != HighGrass {
		b.RemoveMemoriesByEventAtLocation("Found grass supply", targetLocation)
		b.SyncKnownTileFromWorld(targetLocation)
		return
	}

	if _, err := b.Owner.WorldProvider.RemovePlant(tile.Plant); err != nil {
		return
	}
	b.SyncKnownTileFromWorld(targetLocation)

	carriedGrass := b.createCarriedGrassItem()
	if carriedGrass == nil {
		return
	}
	carriedGrass.Location = b.Owner.Location
	if !b.carryItemInHand(carriedGrass) {
		_ = b.Owner.WorldProvider.AddItem(targetLocation.X, targetLocation.Y, carriedGrass)
		b.SyncKnownTileFromWorld(targetLocation)
		return
	}

	b.MotorCortexCurrentTask = MotorCortexAction{"Carry grass to base", "Walk", b.Base.Location, false, false}
}

func (b *Brain) deliverGrassToBase() {
	if b.Owner.Location != b.Base.Location {
		b.MotorCortexCurrentTask = MotorCortexAction{"Carry grass to base", "Walk", b.Base.Location, false, false}
		return
	}

	delivered := b.removeOwnedItemsByName(grassItemName)
	if delivered == 0 {
		return
	}

	deliveredToGround := 0
	for i := 0; i < delivered; i++ {
		grass := b.createCarriedGrassItem()
		if grass == nil {
			continue
		}
		if err := b.Owner.WorldProvider.AddItem(b.Base.Location.X, b.Base.Location.Y, grass); err != nil {
			continue
		}
		deliveredToGround++
	}
	b.SyncKnownTileFromWorld(b.Base.Location)

	b.ensureBaseStockpile()
	b.Base.Stockpile[grassStockpileKey] += deliveredToGround
}

func (b *Brain) createCarriedFoodItem(fruit Fruit) *Item {
	itemName := foodCarryItemName
	if fruit.Name != "" && fruit.Name != foodCarryItemName {
		itemName = foodCarryItemName
	}

	return &Item{
		Name:      itemName,
		Sharpness: 0,
		Bluntness: 0,
		Weight:    1,
		Material:  []Material{GetMaterialByName("Grass")},
		Residues:  make([]Residue, 0),
		Location:  b.Owner.Location,
	}
}

func (b *Brain) createCarriedGrassItem() *Item {
	grass := CreateNewItem(grassItemName)
	if grass != nil {
		return grass
	}

	return &Item{
		Name:      grassItemName,
		Sharpness: 0,
		Bluntness: 0,
		Weight:    1,
		Material:  []Material{GetMaterialByName("Grass")},
		Residues:  make([]Residue, 0),
		Location:  b.Owner.Location,
	}
}

func (b *Brain) hasFreeHand() bool {
	if b.Owner.Body == nil {
		return false
	}
	if b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.Hand != nil && len(b.Owner.Body.RightArm.Hand.Items) == 0 {
		return true
	}
	if b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.Hand != nil && len(b.Owner.Body.LeftArm.Hand.Items) == 0 {
		return true
	}
	return false
}

func isCriticalCarryItem(itemName string) bool {
	switch itemName {
	case "Stone Axe", "Food Box", "Wooden Box", "Wooden Crate", "Woven Grass Basket":
		return true
	default:
		return false
	}
}

func (b *Brain) ensureCarryCapacityFor(targetItemName string) bool {
	if b.hasFreeHand() {
		return true
	}

	dropFromRight := ""
	if b.Owner.Body != nil && b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.Hand != nil {
		for _, item := range b.Owner.Body.RightArm.Hand.Items {
			if item == nil {
				continue
			}
			if item.Name == targetItemName {
				// Keep already collected target resources in hand.
				continue
			}
			if isCriticalCarryItem(item.Name) {
				continue
			}
			dropFromRight = item.Name
			break
		}
	}
	if dropFromRight != "" {
		b.Owner.DropFromRightHand(dropFromRight)
		return b.hasFreeHand()
	}

	dropFromLeft := ""
	if b.Owner.Body != nil && b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.Hand != nil {
		for _, item := range b.Owner.Body.LeftArm.Hand.Items {
			if item == nil {
				continue
			}
			if item.Name == targetItemName {
				// Keep already collected target resources in hand.
				continue
			}
			if isCriticalCarryItem(item.Name) {
				continue
			}
			dropFromLeft = item.Name
			break
		}
	}
	if dropFromLeft != "" {
		b.Owner.DropFromLeftHand(dropFromLeft)
		return b.hasFreeHand()
	}

	return false
}

func (b *Brain) carryItemInHand(item *Item) bool {
	if item == nil {
		return false
	}
	if !b.ensureCarryCapacityFor(item.Name) {
		return false
	}

	if b.Owner.Body != nil && b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.Hand != nil && len(b.Owner.Body.RightArm.Hand.Items) == 0 {
		b.Owner.GrabWithRightHand(item)
		return true
	}
	if b.Owner.Body != nil && b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.Hand != nil && len(b.Owner.Body.LeftArm.Hand.Items) == 0 {
		b.Owner.GrabWithLeftHand(item)
		return true
	}

	return false
}

func (b *Brain) countOwnedItemsByName(itemName string) int {
	count := 0
	for _, item := range b.Owner.OwnedItems {
		if item != nil && item.Name == itemName {
			count++
		}
	}
	return count
}

func filterItemsByName(items []*Item, itemName string) ([]*Item, int) {
	filtered := make([]*Item, 0, len(items))
	removed := 0
	for _, item := range items {
		if item != nil && item.Name == itemName {
			removed++
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered, removed
}

func (b *Brain) removeOwnedItemsByName(itemName string) int {
	filteredOwnedItems, removedOwnedItems := filterItemsByName(b.Owner.OwnedItems, itemName)
	b.Owner.OwnedItems = filteredOwnedItems

	if b.Owner.Body != nil && b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.Hand != nil {
		filteredRightHand, _ := filterItemsByName(b.Owner.Body.RightArm.Hand.Items, itemName)
		b.Owner.Body.RightArm.Hand.Items = filteredRightHand
	}
	if b.Owner.Body != nil && b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.Hand != nil {
		filteredLeftHand, _ := filterItemsByName(b.Owner.Body.LeftArm.Hand.Items, itemName)
		b.Owner.Body.LeftArm.Hand.Items = filteredLeftHand
	}

	return removedOwnedItems
}

func (b *Brain) completeShelterTaskWithFallback() {
	b.RemoveActionFromActionList(TargetedAction{Action: MakeShelter, Target: ""})
	b.ClearCurrentTask()

	b.MotorCortexCurrentTask = MotorCortexAction{
		ActionReason:   "Idle",
		ActionType:     "Idle",
		TargetLocation: b.Owner.Location,
		IsActive:       false,
		Finished:       true,
	}

	// Fall back to exploration when immediate survival stats are stable.
	if b.Owner.Body != nil && b.Owner.Body.Blood.Water >= 40 && b.Owner.Body.Blood.Glucose >= 40 {
		b.GoSearchFor("Nearby resources")
	}
}

func (b *Brain) consumeGroundItemsAtBase(itemName string, amount int) int {
	if amount <= 0 || !b.Base.Claimed {
		return 0
	}

	tile := b.Owner.WorldProvider.GetTile(b.Base.Location.X, b.Base.Location.Y)
	if len(tile.Items) == 0 {
		return 0
	}
	itemsSnapshot := append([]*Item(nil), tile.Items...)

	consumed := 0
	for _, item := range itemsSnapshot {
		if consumed >= amount {
			break
		}
		if item == nil || item.Name != itemName {
			continue
		}
		if err := b.Owner.WorldProvider.DestroyItem(item); err != nil {
			continue
		}
		consumed++
	}

	return consumed
}
