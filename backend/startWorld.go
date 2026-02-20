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
	newPlant := NewPlant(plant, &w.Tiles[y][x], x, y)
	if err := w.AddPlant(x, y, newPlant); err != nil {
		return nil
	}
	//newPlant.PlantLife.turnOn()

	return newPlant
}

func (w *World) initializeTiles(width, height int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Tiles = make([][]Tile, height)
	w.Width = width
	w.Height = height

	for y := range w.Tiles {
		w.Tiles[y] = make([]Tile, width)
		for x := range w.Tiles[y] {
			w.Tiles[y][x] = NewTile(Grass, x, y)
		}
	}
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

func (w *World) MakeItemsAroundLocation(x, y, radius int, itemName string) {
	xMin := x - radius
	yMin := y - radius
	xMax := x + radius
	yMax := y + radius

	for i := xMin; i <= xMax; i++ {
		for j := yMin; j <= yMax; j++ {
			if i == x && j == y {
				continue
			}

			dx := float64(i - x)
			dy := float64(j - y)
			distance := math.Sqrt(dx*dx + dy*dy)
			probability := 0.75 - (distance / float64(radius))

			if probability <= 0 || rand.Float64() >= probability {
				continue
			}
			if i <= 0 || j <= 0 || i >= SIZE_OF_MAP || j >= SIZE_OF_MAP {
				continue
			}
			if w.GetTileType(i, j) != Grass {
				continue
			}
			tile := w.GetTile(i, j)
			if tile.Plant != nil || tile.Shelter != nil {
				continue
			}

			item := CreateNewItem(itemName)
			if item == nil {
				continue
			}
			_ = w.AddItem(i, j, item)
		}
	}
}

func (w *World) seedInitialState() *Entity {
	// Lets time how long this function takes to run
	start := time.Now()
	// Create people
	newPerson1 := w.CreateNewPersonEntity(2, 2, Human)
	newPerson2 := w.CreateNewPersonEntity(50, 50, Human)
	newPerson3 := w.CreateNewPersonEntity(60, 20, Human)
	newPerson4 := w.CreateNewPersonEntity(50, 5, Human)

	if newPerson1 == nil || newPerson2 == nil || newPerson3 == nil || newPerson4 == nil {
		end := time.Now()
		fmt.Println("Time taken to initialize world: ", end.Sub(start))
		return nil
	}
	newPerson1.Title = "Leader"
	newPerson2.Title = "Person 2"
	newPerson3.Title = "Person 3"
	newPerson4.Title = "Person 4"
	newPerson1.Thinking = "I am the leader of this group."
	newPerson2.Thinking = "I am the follower of the leader."
	newPerson3.Thinking = "I am the follower of the leader."
	newPerson4.Thinking = "I am the follower of the leader."

	stoneAxe := CreateNewItem("Stone Axe")
	newPerson1.GrabWithRightHand(stoneAxe)

	// Add a woven grass basket to the world
	wovenGrassBasket := CreateNewItem("Woven Grass Basket")
	if wovenGrassBasket != nil {
		_ = w.AddItem(1, 1, wovenGrassBasket)
		newPerson1.OwnedItems = append(newPerson1.OwnedItems, wovenGrassBasket)
	}

	w.MakeLakeAroundLocation(20, 20, 5)
	w.MakeLakeAroundLocation(30, 70, 20)
	w.MakeLakeAroundLocation(50, 30, 10)
	w.MakeLakeAroundLocation(50, 5, 15)
	w.MakeLakeAroundLocation(82, 50, 7)

	// Add some lumber trees
	w.MakePlantsAroundLocation(10, 50, 15, OakTree)
	w.MakePlantsAroundLocation(60, 65, 15, OakTree)
	w.MakePlantsAroundLocation(80, 13, 10, OakTree)
	w.MakePlantsAroundLocation(89, 50, 7, OakTree)

	// Add some apple trees
	w.MakePlantsAroundLocation(30, 25, 5, AppleTree)
	w.MakePlantsAroundLocation(37, 40, 5, AppleTree)
	w.MakePlantsAroundLocation(50, 50, 5, AppleTree)
	w.MakePlantsAroundLocation(71, 49, 7, AppleTree)

	// Add some high grass
	w.MakePlantsAroundLocation(10, 10, 10, HighGrass)
	w.MakePlantsAroundLocation(20, 20, 15, HighGrass)
	w.MakePlantsAroundLocation(50, 30, 10, HighGrass)
	w.MakePlantsAroundLocation(82, 50, 10, HighGrass)

	// Add stone clusters as ground items.
	w.MakeItemsAroundLocation(61, 39, 5, "Stone")
	w.MakeItemsAroundLocation(38, 50, 5, "Stone")

	// Start world-level plant lifecycle (fruit growth/spawn).
	w.StartPlantLifecycle()

	// Route runtime brain writes through a single intent-applier loop.
	w.EnableIntentEngine()
	newPerson1.Brain.turnOn()
	newPerson2.Brain.turnOn()
	newPerson3.Brain.turnOn()
	newPerson4.Brain.turnOn()

	end := time.Now()
	fmt.Println("Time taken to initialize world: ", end.Sub(start))

	return newPerson1
}

func (w *World) ResetWorld() *Entity {
	w.StopPlantLifecycle()

	for _, person := range w.GetAllPersons() {
		if person == nil || person.Brain == nil {
			continue
		}
		if person.Brain.Active {
			person.Brain.turnOff("world reset")
		}
	}

	width := w.Width
	height := w.Height
	if width <= 0 {
		width = SIZE_OF_MAP
	}
	if height <= 0 {
		height = SIZE_OF_MAP
	}

	w.initializeTiles(width, height)

	return w.seedInitialState()
}

func InitializeWorld() (*World, *Entity) {
	world := NewWorld(SIZE_OF_MAP, SIZE_OF_MAP)
	newPerson1 := world.seedInitialState()

	return world, newPerson1
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
