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
	newPlant := NewPlant(plant, &w.Tiles[x][y], x, y)
	w.AddPlant(x, y, newPlant)
	//newPlant.PlantLife.turnOn()

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

func InitializeWorld() *World {
	world := NewWorld(SIZE_OF_MAP, SIZE_OF_MAP)
	// Lets time how long this function takes to run
	start := time.Now()
	// Create people
	newPerson1 := world.CreateNewPersonEntity(2, 2, Human)
	newPerson1.Title = "Leader"
	newPerson1.Thinking = "I am the leader of this group."
	newPerson1.Brain.PhysiologicalNeeds.Thirst = 70
	newPerson1.Brain.turnOn()

	// Create wolf
	wolf1 := world.CreateNewAnimalEntity(Wolf, 50, 10)
	wolf2 := world.CreateNewAnimalEntity(Wolf, 50, 12)
	wolf3 := world.CreateNewAnimalEntity(Wolf, 48, 10)

	// Add relationships of the wolves to each other
	wolf1.AddRelationship(wolf2, "Pack member", 100)
	wolf1.AddRelationship(wolf3, "Pack member", 100)
	wolf2.AddRelationship(wolf1, "Pack member", 100)
	wolf2.AddRelationship(wolf3, "Pack member", 100)
	wolf3.AddRelationship(wolf1, "Pack member", 100)
	wolf3.AddRelationship(wolf2, "Pack member", 100)

    wolf1.Brain.turnOn()
    wolf2.Brain.turnOn()
    wolf3.Brain.turnOn()

	stoneAxe := CreateNewItem("Stone Axe")
	newPerson1.GrabWithRightHand(stoneAxe)

	// Add a woven grass basket to the world
	wovenGrassBasket := items[6]
	world.AddItem(1, 1, &wovenGrassBasket)

	newPerson1.OwnedItems = append(newPerson1.OwnedItems, &wovenGrassBasket)

	world.MakeLakeAroundLocation(20, 20, 5)
	world.MakeLakeAroundLocation(30, 70, 20)
	world.MakeLakeAroundLocation(50, 30, 10)
	world.MakeLakeAroundLocation(80, 50, 15)

	// Add some lumber trees
    world.MakePlantsAroundLocation(10, 50, 20, OakTree)

	// Add some apple trees
	world.MakePlantsAroundLocation(30, 25, 5, AppleTree)

	// Add some high grass
	world.MakePlantsAroundLocation(10, 10, 10, HighGrass)
	world.MakePlantsAroundLocation(20, 20, 30, HighGrass)

	end := time.Now()
	fmt.Println("Time taken to initialize world: ", end.Sub(start))

	world.LaunchGame(newPerson1)
	
	return world
}

func TestAttack(w *World, person1 *Entity, person2 *Entity, d time.Duration) {
	damage := person1.AttackWithArm(person2, "Head", person1.Body.RightArm.Hand)
		// This should probably return a result of the attack
		if damage.AmountBluntDamage > 0 || damage.AmountSharpDamage > 0 {
			bloodResidue := Residue{"Blood", 1}
			person1.CombatExperience += 1
			person1.AddResidue("RightHand", bloodResidue)
		}
}