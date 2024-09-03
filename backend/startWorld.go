package main

import (
	"fmt"
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
func (w *World) AddPlantToTheWorld(x, y int, plant string) *Plant {
	newPlant := NewPlant(plant, &w.Tiles[x][y], x, y)
	w.AddPlant(x, y, newPlant)
	newPlant.PlantLife.turnOn()

	return newPlant
}

// Display the map of the world in the console, super simple for now
func (w *World) DisplayMap() {
	if true {
		fmt.Println()
		fmt.Println()
		for y := 0; y < w.Height; y++ {
			for x := 0; x < w.Width; x++ {
				if w.Tiles[x][y].Person != nil {
					if w.Tiles[x][y].Person.Brain.IsAlive {
						fmt.Print(Red + "P" + Reset)
					} else {
						fmt.Print("X")
					}
				} else if w.Tiles[x][y].Plant != nil {
					if w.Tiles[x][y].Plant.Name == "Apple Tree" {
						fmt.Print("A")
					} else if w.Tiles[x][y].Plant.Name == "Oak Tree" {
						fmt.Print("T")
					}
				} else if w.Tiles[x][y].Type == 1 {
					fmt.Print(Blue + "W" + Reset)
				} else if w.Tiles[x][y].Shelter != nil {
					fmt.Print(Yellow + "S" + Reset)
				} else {
					fmt.Print(Green + "G" + Reset)
				}

			}
			fmt.Println()
		}
	}
	time.Sleep(1 * time.Second)
}

func InitializeWorld() *World {
	world := NewWorld(100, 100)

	// Create people
	newPerson1 := world.createNewPerson(2, 2)
	newPerson1.Title = "Leader"
	newPerson1.Thinking = "I am the leader of this group."
	newPerson1.Brain.PhysiologicalNeeds.Thirst = 70

	world.AddPerson(2, 2, newPerson1)


	// Create wolf
	wolf := world.CreateNewAnimalByType("Wolf", 3, 3)
	world.AddAnimal(3, 3, wolf)


	stoneAxe := CreateNewItem("Stone Axe")
	newPerson1.GrabRight(stoneAxe)


	// Add a woven grass basket to the world
	wovenGrassBasket := items[6]
	world.AddItem(1, 1, &wovenGrassBasket)

	newPerson1.OwnedItems = append(newPerson1.OwnedItems, &wovenGrassBasket)

	// Add a plant
	appleTree := world.AddPlantToTheWorld(5, 5, "Apple Tree")

	// Add 10 fruits to the apple tree
	apple := CreateNewFruit("Apple", 5, true, 20)

	for i := 0; i < 10; i++ {
		appleTree.Fruit = append(appleTree.Fruit, apple)
	}

	// Add some lumber trees
	world.AddPlantToTheWorld(3, 6, "Oak Tree")
	world.AddPlantToTheWorld(4, 4, "Oak Tree")
	world.AddPlantToTheWorld(8, 6, "Oak Tree")
	world.AddPlantToTheWorld(9, 7, "Oak Tree")
	world.AddPlantToTheWorld(8, 8, "Oak Tree")

	// Add water to the map to test the A* algorithm
	world.SetTileType(0, 1, 1)
	world.SetTileType(0, 2, 1)
	world.SetTileType(0, 3, 1)
	world.SetTileType(1, 2, 1)
	world.SetTileType(1, 3, 1)

	newPerson1.Brain.turnOn()
	wolf.Brain.turnOn()

	return world
}

func TestAttack(w *World, person1 *Person, person2 *Person, d time.Duration) {
	damage := person1.AttackWithArm(person2, "Head", person1.Body.RightArm.Hand)
			// This should probably return a result of the attack
			if damage.AmountBluntDamage > 0 || damage.AmountSharpDamage > 0 {
				bloodResidue := Residue{"Blood", 1}
				person1.CombatExperience += 1
				person1.AddResidue("RightHand", bloodResidue)
			}
}