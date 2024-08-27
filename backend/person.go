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
	WithPerson string
	Relationship string
	Intensity int
}

type TargetedAction struct {
	Action string
	Target string
	IsActive bool
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
	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsHolding        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction
	IsWorking        TargetedAction
	Thinking         string
	WantsTo          string
	Inventory        []Item
	Relationships    []Relationship
	Personality 	 string
	Genes            []string
	Brain			 Brain
	VisionRange 	 int
	WorldProvider    WorldAccessor
}

func NewPerson(worldAccessor WorldAccessor, x, y int) *Person {
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
		Location:         Location{X: x, Y: y},
		IsMoving:         TargetedAction{},
		IsTalking:        TargetedAction{ Action: "Talk", Target: "", IsActive: false },
		IsSitting:        TargetedAction{},
		IsHolding:        TargetedAction{},
		IsEating:         TargetedAction{},
		IsSleeping:       TargetedAction{},
		IsWorking:        TargetedAction{},
		Thinking:         "",
		WantsTo:          "",
		Inventory:        []Item{},
		Relationships:    []Relationship{},
		Personality:      "Talkative",
		Genes:            []string{},

		Brain:            *brain,
		VisionRange:      5,
		WorldProvider:    worldAccessor,

	}

	person.Brain.owner = person
	fmt.Printf("%s has been created\n", person.FullName)

	return person
}
// UpdateLocation updates the location of the person
func (p *Person) UpdateLocation(x, y int) {
	p.Location.X = x
	p.Location.Y = y
}

func (p *Person) GetVision() Vision {
    return p.WorldProvider.GetVision(p.Location.X, p.Location.Y, p.VisionRange)
}

func (p *Person) GetPersonByFullName(FullName string) *Person {
	return p.WorldProvider.GetPersonByFullName(FullName)
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

func (p *Person) addRelationship(person PersonCleaned, relationship string, intensity int) {
	p.Relationships = append(p.Relationships, Relationship{WithPerson: person.FullName, Relationship: relationship, Intensity: intensity})
}

func (p *Person) hasRelationship(fullName string) bool {
	for _, relationship := range p.Relationships {
		if relationship.WithPerson == fullName {
			return true
		}
	}
	return false
}

func (p *Person) removeRelationship(person Person) {
	for i, relationship := range p.Relationships {
		if relationship.WithPerson == person.FullName {
			p.Relationships = append(p.Relationships[:i], p.Relationships[i+1:]...)
			break
		}
	}
}

func (p *Person) updateRelationship(fullName string, relationship string, intensity int) {
	for i, rel := range p.Relationships {
		if rel.WithPerson == fullName {
			p.Relationships[i].Relationship = relationship
			p.Relationships[i].Intensity = intensity
			break
		}
	}
}

func (w *World) createNewPerson(x, y int) *Person {
    person := NewPerson(w, x, y)
    return person
}
