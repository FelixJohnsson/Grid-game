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
	appleTree := NewPlant("Apple Tree", &world.Tiles[5][5])
	world.AddPlant(5, 5, appleTree)
	appleTree.PlantLife.turnOn()

	// Test the attack function
	TestAttack(world, newPerson1, newPerson2, 2*time.Second)

	return world
}

func TestAttack(w *World, person1 *Person, person2 *Person, d time.Duration) {
	// Test the attack function every 2 seconds
	go func() {
		for {
			time.Sleep(d)
			if person2.Body.Head == nil {
				fmt.Println(person2.FullName, " doesnt have a head anymore.")
				break
			}
			person1.Attack(person2, "Head")
		}
	}()
}