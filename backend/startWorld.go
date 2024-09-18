package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var Reset = "\033[0m" 
var Red = "\033[31m" 
var Green = "\033[32m" 
var Yellow = "\033[33m" 
var Blue = "\033[34m" 
var Magenta = "\033[35m" 
var Cyan = "\033[36m" 
var Gray = "\033[37m" 
var White = "\033[97m"

// AddPlantToTheWorld adds a plant to the world at the given location.
func (w *World) AddPlantToTheWorld(x, y int, plant PlantType) *Plant {
	newPlant := w.NewPlant(plant, w.Tiles[x][y], x, y)
	w.AddPlant(x, y, newPlant)

	return newPlant
}

func (w *World) MakeLakeAroundLocation(x, y, radius int) {
    xMin := x - radius
    yMin := y - radius
    xMax := x + radius
    yMax := y + radius

    // Set the center tile as water
    w.SetTileType(x, y, 1)

    for i := xMin; i <= xMax; i++ {
        for j := yMin; j <= yMax; j++ {
            if i == x && j == y {
                continue
            }

            // Calculate distance from center
            dx := float64(i - x)
            dy := float64(j - y)
            distance := math.Sqrt(dx*dx + dy*dy)

            // Probability of being water decreases with distance
            probability := 1.0 - (distance / float64(radius))

            // Add some randomness
            if rand.Float64() < probability {
                w.SetTileType(i, j, 1)
            }
        }
    }

    // Smooth the edges
    w.SmoothLakeEdges(xMin, yMin, xMax, yMax)
}

func (w *World) SmoothLakeEdges(xMin, yMin, xMax, yMax int) {
    for i := xMin; i <= xMax; i++ {
        for j := yMin; j <= yMax; j++ {
			if i > 0 && j > 0 && i < SIZE_OF_MAP && j < SIZE_OF_MAP {
				if w.GetTileType(i, j) == 1 { // If it's water
                // Count water neighbors
                waterNeighbors := 0
                for di := -1; di <= 1; di++ {
                    for dj := -1; dj <= 1; dj++ {
                        if di == 0 && dj == 0 {
                            continue
                        }
                        if w.GetTileType(i+di, j+dj) == 1 {
                            waterNeighbors++
                        }
                    }
                }
                // If fewer than 4 water neighbors, chance to revert to land
                if waterNeighbors < 4 && rand.Float64() < 0.5 {
                    w.SetTileType(i, j, 0) // Set to land
                }
            }
			}
        }
    }
}

func (w *World) MakePlantsAroundLocation(x, y, radius int, plantType PlantType) {
    xMin := x - radius
    yMin := y - radius
    xMax := x + radius
    yMax := y + radius

    for i := xMin; i <= xMax; i++ {
        for j := yMin; j <= yMax; j++ {
            if i == x && j == y {
                continue
            }

            // Calculate distance from center
            dx := float64(i - x)
            dy := float64(j - y)
            distance := math.Sqrt(dx*dx + dy*dy)

            // Probability of being water decreases with distance
            probability := 0.8 - (distance / float64(radius))

            // Add some randomness
            if rand.Float64() < probability && i > 0 && j > 0 && i < SIZE_OF_MAP && j < SIZE_OF_MAP {
                if w.GetTileType(i, j) == 0 && w.GetTile(i, j).Plant == nil {
                    w.AddPlantToTheWorld(i, j, plantType)
                }
            }
        }
    }
}

func (w *World) MakeRiver() {
    // Just a straight horizontal line for now, x is veritcal, y is horizontal
    for x := 0; x < SIZE_OF_MAP; x++ {
        if x < SIZE_OF_MAP {
            w.SetTileType(20, x, 1)
            w.SetTileType(21, x, 1)
            w.SetTileType(22, x, 1)
        }
    }
}

func (w *World) PopulateWorld() {
    w.MakeRiver()
    w.AddPlantToTheWorld(15, 15, AppleTree)
}

func InitializeWorld() *World {
	world := NewWorld(SIZE_OF_MAP, SIZE_OF_MAP)
	// Lets time how long this function takes to run
	start := time.Now()

	end := time.Now()
	fmt.Println("Time taken to initialize world: ", end.Sub(start))

    world.PopulateWorld()

    go world.TileNutritionalValueTicker()

	world.LaunchGame()
	
	return world
}
