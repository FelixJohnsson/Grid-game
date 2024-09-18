package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type PlantType string

// All kinda of available plant types
const (
    AppleTree PlantType = "Apple Tree"
    OakTree PlantType = "Oak Tree"
    HighGrass PlantType = "High Grass"
    LowGrass PlantType = "Low Grass"
    Flower PlantType = "Flower"
)

// NewPlant creates a new plant with the given name.
func (w *World) NewPlant(name PlantType, tile *Tile, x, y int) *Plant {
	newPlant := &Plant{
		Name:          name,
		PlantStage:    0,
		Age:           0,
		Health:        100,
		IsAlive:       true,
		ProducesFruit: false,
		Fruit:         []Fruit{},
		PlantLife:     nil,
        Location:      Location{x, y},
        Nutrients: 0,
        Tile: tile,

        WorldProvider: w,
	}

    newPlant.PlantLife = NewPlantLife(tile, newPlant)

    if name == AppleTree {
        newPlant.ProducesFruit = true
        newPlant.Fruit = append(newPlant.Fruit, Fruit{"Apple", 5, true, 20})
    }

    return newPlant
}

// NewPlantLife creates a new plant life and assigns an owner to it.
func NewPlantLife(tile *Tile, owner *Plant) *PlantLife {
    ctx, cancel := context.WithCancel(context.Background())
    return &PlantLife{
        active:  false,
        ctx:     ctx,
        cancel:  cancel,
        currentAction: PlantAction{},
        actions: []PlantAction{
            {"Gathering", tile, 10},
        },
        Owner: owner,
    }
}

func (pl *PlantLife) turnOn() {
    if pl.active {
        return
    }

    pl.active = true

    go pl.mainLoop()
}

func (pl *PlantLife) turnOff() {
    fmt.Println(Red, "Plant life is shutting down.", Reset)
    pl.active = false
    pl.cancel()
}

func (pl *PlantLife) mainLoop() {
    for {
        select {
        case <-pl.ctx.Done():
            pl.active = false
            return
        default:
            time.Sleep(PLANT_SIMULATION_RATE * time.Millisecond)
            pl.SustainVitals()
            if pl.Owner.Nutrients < -50 {
                pl.Owner.WorldProvider.RemovePlant(pl.Owner)
                pl.turnOff()
            }
            pl.ActionDecider()
            pl.ActionHandler()
            pl.ClearActions()
        }
    }
}
// Rank the tasks based on priority
func (pl *PlantLife) RankTasks() PlantAction {
    highestPriority := 0
    highestPriorityAction := PlantAction{"Idle", pl.Owner.Tile, 0}

    for _, action := range pl.actions {
        if action.Priority > highestPriority {
            highestPriority = action.Priority
            highestPriorityAction = action
        }
    }

    return highestPriorityAction
}

func (pl *PlantLife) ClearActions() {
    pl.actions = []PlantAction{}
}

func (pl *PlantLife) AddAction(action PlantAction) {
    pl.actions = append(pl.actions, action)
}

func (pl *PlantLife) SustainVitals() {
    pl.Owner.Nutrients -= 1
}


func (pl *PlantLife) ActionDecider() {
    if pl.Owner.Nutrients > 50 && pl.Owner.PlantStage > 5 && len(pl.Owner.Fruit) < 10 {
        action := PlantAction{"Fruiting", pl.Owner.Tile, 3}
        pl.AddAction(action)
    }
    if pl.Owner.Nutrients < 100 {
        action := PlantAction{"Gathering", pl.Owner.Tile, 2}
        pl.AddAction(action)
    }
    if pl.Owner.Nutrients > 50 && pl.Owner.PlantStage <= 5 {
        action := PlantAction{"Grow", pl.Owner.Tile, 1}
        pl.AddAction(action)
    } 
    if len(pl.Owner.Fruit) >= 10 && pl.Owner.PlantStage > 5 && pl.Owner.Nutrients > 50 {
        action := PlantAction{"Drop fruits", pl.Owner.Tile, 3}
        pl.AddAction(action)
    }
    if pl.Owner.Nutrients < 0 {
        fmt.Println(Yellow, "Nutrients are running low.", pl.Owner.Tile.NutritionalValue, Reset)
    }
}

func (pl *PlantLife) ActionHandler() {
    highestRankedAction := pl.RankTasks()

    switch highestRankedAction.Name {
    case "Gathering":
        pl.GatheringActionHandler()
    case "Fruiting":
        pl.FruitingActionHandler()
    case "Grow":
        pl.GrowActionHandler()
    case "Drop fruits":
        pl.DropFruitsActionHandler()
        fmt.Println(Blue, "Dropping fruits:", len(pl.Owner.Fruit), Reset)
    }
}

func (pl *PlantLife) GatheringActionHandler() {
    nutrientsTaken := pl.Owner.Tile.NutritionalValue/100
    if nutrientsTaken > 5 {
        nutrientsTaken = 5
    }
    pl.Owner.Nutrients += nutrientsTaken
    pl.Owner.Tile.NutritionalValue -= 1
}

func (pl *PlantLife) FruitingActionHandler() {
    if pl.Owner.Fruit != nil {
        newFruit := CreateNewFruit("Apple", 5, false, 10)
        pl.Owner.Fruit = append(pl.Owner.Fruit, newFruit)
        pl.Owner.Nutrients -= 30
    }
}

func (pl *PlantLife) GrowActionHandler() {
    pl.Owner.PlantStage += 0.1
    pl.Owner.Nutrients -= 10
}

func (pl *PlantLife) DropFruitsActionHandler() {
    pl.Owner.Fruit = []Fruit{}

    randomNumber := rand.Intn(100)
    if randomNumber < 99 {
        pl.Owner.WorldProvider.PlantFruit(pl.Owner, pl.Owner.Location.X, pl.Owner.Location.Y)
    }
}

func CreateNewFruit(name string, ripeAge int, isRipe bool, nutritionalValue int) Fruit {
    return Fruit{
        Name: name,
        RipeAge: ripeAge,
        IsRipe: isRipe,
        NutritionalValue: nutritionalValue,
    }
}