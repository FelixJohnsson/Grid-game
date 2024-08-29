package main

import (
	"context"
	"fmt"
	"time"
)

type PlantAction struct {
	Name string
	Target *Tile
	Priority int
}

type PlantLife struct {
	active bool
	ctx    context.Context
	cancel context.CancelFunc
	actions []PlantAction
}

type Nutrients struct {
	Calories int
	Carbs    int
	Protein  int
	Fat      int

	Vitamins int
	Minerals int
}

type Fruit struct {
	Name      string
	Taste     string
	Age       int
	RipeAge   int
	IsRipe    bool
	Nutrients []Nutrients
}

type PlantStage int

const (
	Seed PlantStage = iota
	Sprout
	Vegetative
	Flowering
	Fruiting
)

type Plant struct {
	Name          string
	Age           int
	Health        int
	IsAlive       bool
	ProducesFruit bool
	Fruit         []Fruit
	PlantStage    PlantStage
	PlantLife     *PlantLife
}

// NewPlant creates a new plant with the given name.
func NewPlant(name string, tile *Tile) *Plant {
	plantLife := NewPlantLife(tile)

	return &Plant{
		Name:          name,
		PlantStage:    0,
		Age:           0,
		Health:        100,
		IsAlive:       true,
		ProducesFruit: false,
		Fruit:         []Fruit{},
		PlantLife:     plantLife,
	}
}

// NewPlantLife creates a new plant life and assigns an owner to it.
func NewPlantLife(tile *Tile) *PlantLife {
    ctx, cancel := context.WithCancel(context.Background())
    return &PlantLife{
        active:  false,
        ctx:     ctx,
        cancel:  cancel,
        actions: []PlantAction{
            {"Gathering resources from tile", tile, 10},
        },
    }
}

func (l *PlantLife) turnOn() {
    if l.active {
        fmt.Println("Plant life is already active.")
        return
    }

    fmt.Println("Plant life is now active.")
    l.active = true

    go l.mainLoop()
}

func (l *PlantLife) mainLoop() {
    for {
        select {
        case <-l.ctx.Done():
            fmt.Println("Plant life is shutting down.")
            l.active = false
            return
        default:
            // PlantLife logic goes here

            // Sleep or yield for a bit to prevent CPU hogging
            time.Sleep(5000 * time.Millisecond)
        }
    }
}
