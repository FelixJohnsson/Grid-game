package main

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

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
	RequiresLimb []string
}

type Wearable struct {
	Name string
	Material string
	Protection int
}

// Body status
type LimbStatus struct {
	Damage int
	IsBleeding bool
	IsBroken bool
	Residues []string
	CoveredWith []Wearable
}

type LimbThatCanHold struct {
	LimbStatus
	Items []Item
	WeightOfItems int
}

// Available actions
var actions = []string {
	// World actions
	"Move",
	"Talk",

	"Grab",
	"Drop",

	"Sit",
	"Hold",
	"Eat",
	"Sleep",
	"Work",
	"Throw",
	"Build",
	"Dig",
	"Plant",
	"Harvest",
	"Chop",
	"Mine",

	"Open",
	"Close",
	"Enter",
	"Exit",
	"Use",

	// Hostile actions
	"Attack",
	"Steal",
	"Destroy",

	// Friendly actions
	"Help",
	"Gift",
	"Protect",
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
	Personality 	 string
	Genes            []string

	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsHolding        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction
	IsWorking        TargetedAction

	Thinking         string
	WantsTo          string
	FeelingSafe 	 int
	FeelingScared	 int

	RightHand        LimbThatCanHold
	LeftHand 	     LimbThatCanHold
	Back 		     LimbStatus
	LeftFoot 	     LimbStatus
	RightFoot 	     LimbStatus
	Head 		     LimbStatus
	Torso 		     LimbStatus
	Legs 		     LimbStatus

	Relationships    []Relationship

	Brain			 Brain
	VisionRange 	 int
	WorldProvider    WorldAccessor
	Location         Location
	OnTileType 	     TileType
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
		IsMoving:         TargetedAction{},
		IsTalking:        TargetedAction{},
		IsSitting:        TargetedAction{},
		IsHolding:        TargetedAction{},
		IsEating:         TargetedAction{},
		IsSleeping:       TargetedAction{},
		IsWorking:        TargetedAction{},
		Thinking:         "",
		WantsTo:          "",
		FeelingSafe: 	  0,
		FeelingScared:    0,
		Relationships:    []Relationship{},
		Personality:      "",
		Genes:            []string{},

		Brain:            *brain,
		VisionRange:      5,
		Location:         Location{X: x, Y: y},
		WorldProvider:    worldAccessor,
		OnTileType:       0,

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

// Grab in the right hand
func (p *Person) GrabRight(item Item) {
	if p.RightHand.Items == nil {
		p.RightHand.Items = []Item{item}
		// If the item has residues, add them to the limb
		if item.Residues != nil {
			for _, residue := range item.Residues {
				p.AddResidue("RightHand", residue)
			}
		}
	} else {
		fmt.Println("Right hand is already holding something")
	}
}

// Drop from the right hand
func (p *Person) DropRight() {
	if p.RightHand.Items != nil {
		p.RightHand.Items = nil
	} else {
		fmt.Println("Right hand is empty")
	}
}

// Grab in the left hand
func (p *Person) GrabLeft(item Item) {
	if p.LeftHand.Items == nil {
		p.LeftHand.Items = []Item{item}
		// If the item has residues, add them to the limb
		if item.Residues != nil {
			for _, residue := range item.Residues {
				p.AddResidue("RightHand", residue)
			}
		}
	} else {
		fmt.Println("Left hand is already holding something")
	}
}

// Drop from the left hand
func (p *Person) DropLeft() {
	if p.LeftHand.Items != nil {
		p.LeftHand.Items = nil
	} else {
		fmt.Println("Left hand is empty")
	}
}

// AddResidue adds a residue to the limb
func (p *Person) AddResidue(limb string, residue string) {
	switch limb {
	case "Back":
		p.Back.Residues = append(p.Back.Residues, residue)
	case "LeftFoot":
		p.LeftFoot.Residues = append(p.LeftFoot.Residues, residue)
	case "RightFoot":
		p.RightFoot.Residues = append(p.RightFoot.Residues, residue)
	case "Head":
		p.Head.Residues = append(p.Head.Residues, residue)
	case "Torso":
		p.Torso.Residues = append(p.Torso.Residues, residue)
	case "Legs":
		p.Legs.Residues = append(p.Legs.Residues, residue)
	case "RightHand":
		p.RightHand.Residues = append(p.RightHand.Residues, residue)
	case "LeftHand":
		p.LeftHand.Residues = append(p.LeftHand.Residues, residue)
	}
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



// ---------------- Create a new person ----------------

func (w *World) createNewPerson(x, y int) *Person {
    person := NewPerson(w, x, y)
    return person
}
