package main

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)
type Item struct {
	Name string
}

type WorldState struct {
	Map 	 	[][]string;
  	Buildings	[]Building;
  	Persons	 	[]Person;
  	Resources 	Resources;
}

// Enum for Jobs
type Jobs string

const (
	Farmer     Jobs = "Farmer"
	Miner      Jobs = "Miner"
	Lumberjack Jobs = "Lumberjack"
	Builder    Jobs = "Builder"
	Soldier    Jobs = "Soldier"
	Unemployed Jobs = "Unemployed"
)

type Relationship struct {
	WithPerson Person
	Relationship string
	Intensity int
}
type Person struct {
	Age              int
	Title 		     string
	FirstName        string
	FamilyName       string
	FullName 	     string
	Initials         string
	IsChild          bool
	Gender           string
	Description      string
	Icon             string
	Occupation       Jobs
	IsWorkingAt      *Building
	Color            string
	Location         Location
	IsMoving         bool
	IsTalking        bool
	IsSitting        bool
	IsHolding        bool
	IsEating         bool
	IsSleeping       bool
	IsWorking        bool
	Thinking         string
	WantsTo          string
	Inventory        []Item
	Relationships    []Relationship
	CurrentWorldState World
	Genes            []string
	Brain			 Brain
	VisionRange 	 int
	World 			 World
}

func NewPerson() *Person {
	age := rand.Intn(63) + 2
	firstName := gofakeit.FirstName()
	familyName := gofakeit.LastName()
	gender := gofakeit.Gender()
	brain := NewBrain()

	person := &Person{
		Age:              age,
		Title:            "",
		FirstName:        firstName,
		FamilyName:       familyName,
		FullName:         firstName + " " + familyName,
		Initials:         string(firstName[0]) + string(familyName[0]),
		IsChild:          age < 18,
		Gender:           gender,
		Description:      "",
		Icon:             "P",
		Occupation:       Unemployed,
		IsWorkingAt:      nil,
		Color:            "",
		Location:         Location{X: 0, Y: 0},
		IsMoving:         false,
		IsTalking:        false,
		IsSitting:        false,
		IsHolding:        false,
		IsEating:         false,
		IsSleeping:       false,
		IsWorking:        false,
		Thinking:         "",
		WantsTo:          "",
		Inventory:        []Item{},
		Relationships:    []Relationship{},
		Genes:            []string{},

		Brain:            *brain,
		VisionRange:      10,
	}

	person.Brain.owner = person
	fmt.Printf("%s has been created\n", person.FullName)

	person.turnOnBrain()

	return person
}

func (p *Person) turnOnBrain() {
	fmt.Printf("%s is turning on their brain\n", p.FullName)
	p.Brain.turnOn()
}

func (p *Person) addTask(task action) {
	p.Brain.addTask(task)
}

func (p *Person) addEmployer(building *Building) {
	if p.IsChild {
		return
	}
	p.IsWorkingAt = building
	building.Workers = append(building.Workers, *p)
	switch building.Type {
	case "Lumberjack":
		p.Occupation = Lumberjack
	case "Mine":
		p.Occupation = Miner
	case "Farm":
		p.Occupation = Farmer
	default:
		p.Occupation = Unemployed
	}
}

func (p *Person) startWorking() {
	if !p.IsWorking {
		p.IsWorking = true
	}
}

func (p *Person) stopWorking() {
	if p.IsWorking {
		p.IsWorking = false
		fmt.Printf("%s has stopped working\n", p.FullName)
	}
}

func (p *Person) eat() {
	if !p.IsEating {
		p.IsEating = true
		fmt.Printf("%s is eating\n", p.FullName)
	}
}

func (p *Person) sleep() {
	if !p.IsSleeping {
		p.IsSleeping = true
		fmt.Printf("%s is sleeping\n", p.FullName)
	}
}

func (p *Person) wakeUp() {
	if p.IsSleeping {
		p.IsSleeping = false
		fmt.Printf("%s has woken up\n", p.FullName)
	}
}

func (p *Person) talk() {
	if !p.IsTalking {
		p.IsTalking = true
		fmt.Printf("%s is talking\n", p.FullName)
	}
}

func (p *Person) stopTalking() {
	if p.IsTalking {
		p.IsTalking = false
		fmt.Printf("%s has stopped talking\n", p.FullName)
	}
}

func createNewPerson() {
	person := NewPerson()
	addPerson(*person)
}

func getPersons() []Person {
	persons, error := loadPersonsFromFile()
	if error != nil {
		return []Person{}
	}
	return persons
}