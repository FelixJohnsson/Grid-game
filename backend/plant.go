package main

import (
	"context"
	"time"
)

type PlantType string

// All kinda of available plant types
const (
	AppleTree PlantType = "Apple Tree"
	OakTree   PlantType = "Oak Tree"
	HighGrass PlantType = "High Grass"
	LowGrass  PlantType = "Low Grass"
	Flower    PlantType = "Flower"
)

// NewPlant creates a new plant with the given name.
func NewPlant(name PlantType, tile *Tile, x, y int) *Plant {
	plantLife := NewPlantLife(tile)

	newPlant := &Plant{
		Name:          name,
		PlantStage:    0,
		Age:           0,
		Health:        100,
		IsAlive:       true,
		ProducesFruit: false,
		Fruit:         []Fruit{},
		PlantLife:     plantLife,
		Location:      Location{x, y},
	}

	if name == AppleTree {
		newPlant.ProducesFruit = true
		// Seed each apple tree with one immediately edible apple.
		AddFruitToPlant(newPlant, CreateNewFruit("Apple", 0, true, 20))
	}

	return newPlant
}
func AddFruitToPlant(plant *Plant, fruit Fruit) {
	plant.Fruit = append(plant.Fruit, fruit)
}

func HasRipeFruit(plant *Plant) bool {
	if plant == nil {
		return false
	}
	for _, fruit := range plant.Fruit {
		if fruit.IsRipe {
			return true
		}
	}
	return false
}

func RemoveFruitFromPlant(plant *Plant, number int) {
	if plant == nil || number <= 0 || len(plant.Fruit) == 0 {
		return
	}
	if number >= len(plant.Fruit) {
		plant.Fruit = []Fruit{}
		return
	}
	plant.Fruit = plant.Fruit[:len(plant.Fruit)-number]
}

// NewPlantLife creates a new plant life and assigns an owner to it.
func NewPlantLife(tile *Tile) *PlantLife {
	ctx, cancel := context.WithCancel(context.Background())
	return &PlantLife{
		active: false,
		ctx:    ctx,
		cancel: cancel,
		actions: []PlantAction{
			{"Gathering resources from tile", tile, 10},
		},
	}
}

func (l *PlantLife) turnOn() {
	if l.active {
		return
	}

	l.active = true

	go l.mainLoop()
}

func (l *PlantLife) mainLoop() {
	for {
		select {
		case <-l.ctx.Done():
			l.active = false
			return
		default:
			// PlantLife logic goes here

			// Sleep or yield for a bit to prevent CPU hogging
			time.Sleep(5000 * time.Millisecond)
		}
	}
}

func CreateNewFruit(name string, ripeAge int, isRipe bool, nutritionalValue int) Fruit {
	return Fruit{
		Name:             name,
		RipeAge:          ripeAge,
		IsRipe:           isRipe,
		NutritionalValue: nutritionalValue,
	}
}
