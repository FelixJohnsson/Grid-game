package main

import (
	"fmt"
	"math/rand"
	"time"

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
type Person struct {
	Age              int
	Name             string
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
	CurrentWorldState WorldState
	Genes            []string
	Brain            
}

func NewPerson(x, y int) *Person {
	rand.Seed(time.Now().UnixNano())
	age := rand.Intn(63) + 2
	name := gofakeit.Name()
	initials := string(name[0]) + string(name[1])
	gender := gofakeit.Gender()
	brain := NewBrain()

	person := &Person{
		Age:              age,
		Name:             name,
		Initials:         initials,
		IsChild:          age < 18,
		Gender:           gender,
		Description:      "",
		Icon:             "P",
		Occupation:       Unemployed,
		IsWorkingAt:      nil,
		Color:            "",
		Location:         Location{X: x, Y: y},
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
		Genes:            []string{},
		Brain:            *brain,
	}

	person.Brain.owner = person
	fmt.Printf("%s has been created\n", person.Name)

	person.turnOnBrain()

	return person
}

func (p *Person) turnOnBrain() {
	fmt.Printf("%s is turning on their brain\n", p.Name)
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
		fmt.Printf("%s has stopped working\n", p.Name)
	}
}

func (p *Person) eat() {
	if !p.IsEating {
		p.IsEating = true
		fmt.Printf("%s is eating\n", p.Name)
	}
}

func (p *Person) sleep() {
	if !p.IsSleeping {
		p.IsSleeping = true
		fmt.Printf("%s is sleeping\n", p.Name)
	}
}

func (p *Person) wakeUp() {
	if p.IsSleeping {
		p.IsSleeping = false
		fmt.Printf("%s has woken up\n", p.Name)
	}
}

func (p *Person) talk() {
	if !p.IsTalking {
		p.IsTalking = true
		fmt.Printf("%s is talking\n", p.Name)
	}
}

func (p *Person) stopTalking() {
	if p.IsTalking {
		p.IsTalking = false
		fmt.Printf("%s has stopped talking\n", p.Name)
	}
}

func createNewPerson() {
	person := NewPerson(0, 0)
	addPerson(*person)
}

func getPersons() []Person {
	persons, error := loadPersonsFromFile()
	if error != nil {
		return []Person{}
	}
	return persons
}