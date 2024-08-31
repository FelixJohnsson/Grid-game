package main

import (
	"fmt"
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

	newPerson2 := world.createNewPerson(9, 9)
	newPerson2.Title = "Follower"
	newPerson2.Thinking = "I follow the leader."

	world.AddPerson(2, 2, newPerson1)
	world.AddPerson(9, 9, newPerson2)


	// Turn on the brain for the people
	newPerson1.Body.Head.Brain.turnOn()

	// Create a Wooden spear item from items
	woodenSpear := items[0]
	woodenSpear.Residues = append(woodenSpear.Residues, Residue{"Dirt", 1})

	// Create a Wooden staff item from items
	woodenStaff := items[1]
	woodenStaff.Residues = append(woodenStaff.Residues, Residue{"Blood", 1})
	stoneAxe := items[2]

	newPerson1.GrabRight(&stoneAxe)

	// Add the wooden spear to the world
	world.AddItem(1, 1, &woodenSpear)

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

	// Test the FindLumberTrees function
	lumberTreesFound := newPerson1.FindLumberTrees()
	if lumberTreesFound != nil {
		closestLumberTree := newPerson1.FindTheClosestPlant(lumberTreesFound)
		newPerson1.Body.Head.Brain.WalkOverPath(closestLumberTree.Location.X, closestLumberTree.Location.Y)
	} else {
		fmt.Println("No lumber trees found in vision")
	}

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