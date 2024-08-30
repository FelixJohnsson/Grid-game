package main

import (
	"fmt"
	"time"
)

func initializeWorld() *World {
	world := NewWorld(10, 10) 

	// Create people
	newPerson1 := world.createNewPerson(2, 2)
	newPerson1.Title = "Leader"
	newPerson1.Thinking = "I am the leader of this group."

	newPerson2 := world.createNewPerson(9, 9)
	newPerson2.Title = "Follower"
	newPerson2.Thinking = "I follow the leader."

	world.AddPerson(2, 2, newPerson1)
	world.AddPerson(9, 9, newPerson2)


	// Turn on the brain for the people
	newPerson1.Body.Head.Brain.turnOn()
	newPerson2.Body.Head.Brain.turnOn()

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
	appleTree := NewPlant("Apple Tree", &world.Tiles[5][5], 5, 5)
	world.AddPlant(5, 5, appleTree)
	appleTree.PlantLife.turnOn()

	// Add 4 Oak trees
	oakTree1 := NewPlant("Oak Tree", &world.Tiles[2][3], 2, 3)
	world.AddPlant(2, 3, oakTree1)
	oakTree1.PlantLife.turnOn()

	oakTree2 := NewPlant("Oak Tree", &world.Tiles[8][6], 8, 6)
	world.AddPlant(8, 6, oakTree2)
	oakTree2.PlantLife.turnOn()

	oakTree3 := NewPlant("Oak Tree", &world.Tiles[9][7], 9, 7)
	world.AddPlant(9, 7, oakTree3)
	oakTree3.PlantLife.turnOn()

	oakTree4 := NewPlant("Oak Tree", &world.Tiles[8][8], 8, 8)
	world.AddPlant(8, 8, oakTree4)
	oakTree4.PlantLife.turnOn()

	// Test the FindLumberTrees function
	lumberTreesFound := newPerson1.FindLumberTrees()
	if lumberTreesFound != nil {
		closestLumberTree := newPerson1.FindTheClosestPlant(lumberTreesFound)
		fmt.Println("And the closest lumber tree is at", closestLumberTree.Location.X, closestLumberTree.Location.Y)
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