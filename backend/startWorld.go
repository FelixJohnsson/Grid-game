package main

import (
	"time"
)

// AddPlantToTheWorld adds a plant to the world at the given location.
func (w *World) AddPlantToTheWorld(x, y int, plant string) *Plant {
	newPlant := NewPlant(plant, &w.Tiles[x][y], x, y)
	w.AddPlant(x, y, newPlant)
	newPlant.PlantLife.turnOn()

	return newPlant
}

func initializeWorld() *World {
	world := NewWorld(10, 10) 

	// Create people
	newPerson1 := world.createNewPerson(1, 1)
	newPerson1.Title = "Leader"
	newPerson1.Thinking = "I am the leader of this group."
	newPerson1.Body.Head.Brain.PhysiologicalNeeds.Thirst = 40

	world.AddPerson(2, 2, newPerson1)

	// Turn on the brain for the people
	newPerson1.Body.Head.Brain.turnOn()

	stoneAxe := CreateNewItem("Stone Axe")
	newPerson1.GrabRight(stoneAxe)


	// Add a woven grass basket to the world
	wovenGrassBasket := items[6]
	world.AddItem(1, 1, &wovenGrassBasket)

	newPerson1.OwnedItems = append(newPerson1.OwnedItems, &wovenGrassBasket)

	// Add a plant
	appleTree := world.AddPlantToTheWorld(5, 5, "Apple Tree")
	appleTree.Fruit = append(appleTree.Fruit, Fruit{"Apple", "Sweet", 3, 3, true, make([]Nutrients, 0)})

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