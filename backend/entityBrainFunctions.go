package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ----------------- Brain Functions -----------------

func (b *Brain) turnOn() {
    if b.Active {
        fmt.Println("Brain is already active.")
        return
    }

    b.IsConscious = true
    b.Active = true
    go b.MotorCortex()
    go b.MainLoop()
}


func (b *Brain) MainLoop() {
    for {
        select {
        case <-b.Ctx.Done():
            b.Active = false
            return
        default:
            if !b.Active {
                return
            }

            if !b.IsConscious {
                return
            }

            if b.IsUnderAttack.Active {
                b.IsUnderAttackHandler()
            }

            b.OxygenHandler()

            b.PainHandler()
            b.FoodHandler()
            b.ThirstHandler()

            b.ClearWants()
            b.HomoSapiensCalculateWant()
            b.TranslateWantToTaskList()

            obs := b.Owner.WorldProvider.GetVision(b.Owner.Location.X, b.Owner.Location.Y, b.Owner.VisionRange)
            b.CognitiveMapHandler(obs)
            b.ActionHandler()

            // Sleep for 2 seconds
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func (b *Brain) turnOff() {
    if !b.Active {
        fmt.Println("Brain is already inactive.")
        return
    }

	b.Active = false
	b.IsConscious = false
    b.Cancel()
}

// LooseConsciousness makes the person loose consciousness - duration is in seconds
func (b *Brain) LooseConsciousness(duration int) {
    b.IsConscious = false
    go func() {
        time.Sleep(time.Duration(duration) * time.Second)
        b.IsConscious = true
    }()
}


// ----------------- Cognitive Map -----------------

func (b *Brain) CognitiveMapHandler(obs []Tile) {
    for _, tile := range obs {
        cognitiveMapTile := CognitiveMapTile{
            TileType: tile.Type,
            Location: tile.Location,
        }
        if tile.Entity != nil && tile.Entity.FullName != b.Owner.FullName {
            cognitiveMapTile.Entity = CognitiveMapEntity{
                FullName: tile.Entity.FullName,
                SpeciesType: tile.Entity.Species,
                IsAlive: tile.Entity.Brain.IsAlive,
            }
        }
        if tile.Plant != nil {
            cognitiveMapTile.Plant = CognitiveMapPlant{
                Name: tile.Plant.Name,
                IsAlive: tile.Plant.IsAlive,
                ProducesFruit: tile.Plant.ProducesFruit,
                PlantStage: tile.Plant.PlantStage,
            }
        }
        b.AddLocationToCognitiveMap(tile.Location, cognitiveMapTile)
    }
}

func (b *Brain) AddLocationToCognitiveMap(location Location, tileInfo CognitiveMapTile) {
    b.CognitiveMap.KnownTiles[location] = tileInfo
}

func (b *Brain) IsTileKnown(location Location) bool {
    if _, ok := b.CognitiveMap.KnownTiles[location]; ok {
        return true
    } else {
        return false
    }
}

func (b *Brain) GetLocationFromCognitiveMap(location Location) CognitiveMapTile {
    if _, ok := b.CognitiveMap.KnownTiles[location]; !ok {
        return CognitiveMapTile{}
    } else {
        return b.CognitiveMap.KnownTiles[location]
    }
}

func (b *Brain) GetAllWaterTilesFromCognitiveMap() []CognitiveMapTile {
    var waterTiles []CognitiveMapTile
    for _, tile := range b.CognitiveMap.KnownTiles {
        if tile.TileType == 1 {
            waterTiles = append(waterTiles, tile)
        }
    }
    return waterTiles
}

func (b *Brain) GetAllPlantsFromCognitiveMap() []CognitiveMapTile {
    var plants []CognitiveMapTile
    for _, tile := range b.CognitiveMap.KnownTiles {
        if tile.TileType == 2 {
            plants = append(plants, tile)
        }
    }
    return plants
}

func (b *Brain) GetAllFruitingPlantsFromCognitiveMap() []CognitiveMapTile {
    var fruitingPlants []CognitiveMapTile
    for _, tile := range b.CognitiveMap.KnownTiles {
        if tile.Plant.ProducesFruit {
            fruitingPlants = append(fruitingPlants, tile)
        }
    }
    return fruitingPlants
}

func (b *Brain) GetClosestFruitingPlantFromCognitiveMap() CognitiveMapTile {
    var closestPlant CognitiveMapTile
    for _, tile := range b.CognitiveMap.KnownTiles {
        if tile.Plant.ProducesFruit {
            if b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, tile.Location) < b.Owner.WorldProvider.CalculateDistance(b.Owner.Location, closestPlant.Location) {
                closestPlant = tile
            }
        }
    } 
    return closestPlant
}

// ----------------- Pain -----------------------------

// PainHandler is a function that handles the pain level of the person
func (b *Brain) PainHandler() {
    b.CalculatePainLevel()

    if b.PainLevel > b.PainTolerance {
        // Get a random number for duration of unconsciousness
        durationInSeconds := rand.Intn(60)
        
        b.LooseConsciousness(durationInSeconds)
    }
}

// ApplyPain - Apply pain to the person
func (b *Brain) ApplyPain(amount int) {
	b.PainLevel += amount
}

// CalculatePainLevel - Calculate the pain level of the person
func (b *Brain) CalculatePainLevel() {
    // We need to loop over the body parts and check if it's broken or bleeding

    if b.Owner.Body.Head != nil {
        if b.Owner.Body.Head.IsBroken {
            b.ApplyPain(5)
        } 
        if b.Owner.Body.Head.IsBleeding {
            b.ApplyPain(2)
        }
        if b.Owner.Body.Head.Ears != nil && b.Owner.Body.Head.Ears.IsBleeding || b.Owner.Body.Head.Ears.IsBroken {
           b.ApplyPain(1)
        } 
        if b.Owner.Body.Head.Eyes != nil && b.Owner.Body.Head.Eyes.IsBleeding {
           b.ApplyPain(5)
        }
        if b.Owner.Body.Head.Nose != nil && b.Owner.Body.Head.Nose.IsBleeding || b.Owner.Body.Head.Nose.IsBroken {
            b.ApplyPain(2)
        }
        if b.Owner.Body.Head.Mouth != nil  &&b.Owner.Body.Head.Mouth.IsBleeding || b.Owner.Body.Head.Mouth.IsBroken {
            b.ApplyPain(2)
        }
    }
    if b.Owner.Body.Torso.IsBleeding || b.Owner.Body.Torso.IsBroken {
        b.ApplyPain(5)
    }
    if b.Owner.Body.RightArm != nil {
        if  b.Owner.Body.RightArm.IsBleeding || b.Owner.Body.RightArm.IsBroken {
            b.ApplyPain(5)
        }
    }
    if b.Owner.Body.LeftArm != nil {
        if  b.Owner.Body.LeftArm.IsBleeding || b.Owner.Body.LeftArm.IsBroken {
            b.ApplyPain(5)
        }
    }
    if b.Owner.Body.RightLeg != nil {
        if b.Owner.Body.RightLeg.IsBleeding || b.Owner.Body.RightLeg.IsBroken {
            b.ApplyPain(5)
        }
    }
    if b.Owner.Body.LeftLeg != nil {
        if b.Owner.Body.LeftLeg.IsBleeding || b.Owner.Body.LeftLeg.IsBroken {
            b.ApplyPain(5)  
        }
    }
    if b.Owner.Body.RightArm != nil {
        if b.Owner.Body.RightArm.Hand != nil {
            if b.Owner.Body.RightArm.Hand.IsBleeding || b.Owner.Body.RightArm.Hand.IsBroken {
                b.ApplyPain(2)
            }
        }
    }
    if b.Owner.Body.LeftArm != nil {
        if  b.Owner.Body.LeftArm.Hand != nil {
            if  b.Owner.Body.LeftArm.Hand.IsBleeding || b.Owner.Body.LeftArm.Hand.IsBroken {
                b.ApplyPain(2)
            }
        }
    }
    if b.Owner.Body.RightLeg != nil {
        if b.Owner.Body.RightLeg.Foot != nil {
            if b.Owner.Body.RightLeg.Foot.IsBleeding || b.Owner.Body.RightLeg.Foot.IsBroken {
                b.ApplyPain(2)
            }
        }
    }
    if b.Owner.Body.LeftLeg != nil {
        if b.Owner.Body.LeftLeg.Foot != nil {
            if b.Owner.Body.LeftLeg.Foot.IsBleeding || b.Owner.Body.LeftLeg.Foot.IsBroken {
                b.ApplyPain(2)
            }
        }
    }
}

// ----------------- Oxygen levels -----------------

// Breath 
func (b *Brain) Breath() {
    b.OxygenLevel += 10
}

// ConsumeOxygen
func (b *Brain) ConsumeOxygen() {
    b.OxygenLevel -= 10
}

// CheckIfCanBreah - Check if the person can breath
func (b *Brain) CheckIfCanBreath() bool {
    mouthCanBreath := b.Owner.Body.Head.Mouth != nil && !b.Owner.Body.Head.Mouth.IsObstructed
    noseCanBreath := b.Owner.Body.Head.Nose != nil && !b.Owner.Body.Head.Nose.IsObstructed && !b.Owner.Body.Head.Nose.IsBroken

    return mouthCanBreath || noseCanBreath
}

// ----------------- Safety ---------------------

// UnderAttack is called when the person is being attacked
func (b *Brain) UnderAttack(attacker *Entity, targettedLimb BodyPartType, attackersLimb BodyPartType) {
	// Decide between fight or flight
	if !b.Owner.Body.RightArm.IsBroken {
		b.Owner.AttackWithArm(attacker, targettedLimb, b.Owner.Body.RightArm.Hand)
	} else if !b.Owner.Body.LeftArm.IsBroken {
		b.Owner.AttackWithArm(attacker, targettedLimb, b.Owner.Body.LeftArm.Hand)
	} else if !b.Owner.Body.RightLeg.IsBroken {
		b.Owner.AttackWithLeg(attacker, targettedLimb, b.Owner.Body.RightLeg)
	} else if !b.Owner.Body.LeftLeg.IsBroken {
		b.Owner.AttackWithLeg(attacker, targettedLimb, b.Owner.Body.LeftLeg)
	} else {
		//
	}
}

// ----------------- Pathfinding ---------------------

// Decide a path to the target location - Check first if it's physically possible to walk.
func (b *Brain) DecidePathTo(x, y int) []*Node {
    path := b.AStar(b.Owner.Location.X, b.Owner.Location.Y, x, y)
    return path
}

// TakeStepOverPath - Take a step over the path that was decided
func (b *Brain) TakeStepOverPath(MotorCortexAction MotorCortexAction) bool {
    path := b.DecidePathTo(MotorCortexAction.TargetLocation.X, MotorCortexAction.TargetLocation.Y)
    if path == nil {
        fmt.Println(b.Owner.FullName + " could not find a path to the location.")
        return false
    }
    targetNode := path[1]
    b.Owner.WalkStepTo(targetNode.X, targetNode.Y)

    return true
}

// WalkOverPath - Walk over the path that was decided
func (b *Brain) WalkOverPath(MotorCortexAction MotorCortexAction) bool {
    path := b.DecidePathTo(MotorCortexAction.TargetLocation.X, MotorCortexAction.TargetLocation.Y)
    if path == nil {
        return false
    }
    fmt.Println(b.Owner.FullName + " is walking to ", MotorCortexAction.TargetLocation.X, MotorCortexAction.TargetLocation.Y)
    for _, node := range path {
        b.MotorCortexCurrentTask.IsActive = true
		// Wait for a half second before walking to the next node
        b.Owner.WalkStepTo(node.X, node.Y)
    }
	b.MotorCortexCurrentTask.Finished = true
    b.MotorCortexCurrentTask.IsActive = false
    return true
}

// ----------------- Items -------------------------

// FindInOwnedItems - Find an item in the owned items
func (b *Brain) FindInOwnedItems(itemName string) *Item {
    for _, item := range b.Owner.OwnedItems {
        if item.Name == itemName {
            return item
        }
    }
    return nil
}

// HasItemEquippedInRight - Check if the person has an item equipped in right hand
func (b *Brain) HasItemEquippedInRight(itemName string) bool {
    for _, item := range b.Owner.Body.RightArm.Hand.Items {
        if item.Name == itemName {
            return true
        }
    }
    return false
}

// HasItemEquippedInLeft - Check if the person has an item equipped in left hand
func (b *Brain) HasItemEquippedInLeft(itemName string) bool {
    for _, item := range b.Owner.Body.LeftArm.Hand.Items {
        if item.Name == itemName {
            return true
        }
    }
    return false
}