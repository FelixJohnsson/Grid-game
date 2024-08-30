package main

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

func NewPerson(worldAccessor WorldAccessor, x, y int) *Person {
	age := rand.Intn(63) + 2
	firstName := gofakeit.FirstName()
	familyName := gofakeit.LastName()
	gender := gofakeit.Gender()
	body := CreateNewBody()

	fmt.Println("Creating a new person", age, firstName)

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
		Occupation:       Unemployed,
		IsMoving:         TargetedAction{},
		IsTalking:        TargetedAction{},
		IsSitting:        TargetedAction{},
		IsEating:         TargetedAction{},
		IsSleeping:       TargetedAction{},
		Thinking:         "",
		WantsTo:          "",
		FeelingSafe: 	  0,
		FeelingScared:    0,
		Relationships:    []Relationship{},
		Personality:      "",
		Genes:            []string{},

		VisionRange:      5,
		Location:         Location{X: x, Y: y},
		WorldProvider:    worldAccessor,
		OnTileType:       0,
		Body:			  body,

		Strength:         1,
		Agility:          1,
		Intelligence:     1,
		Charisma:         1,
		Stamina:          1,

		CombatExperience: 1,
		CombatSkill:      1,
		CombatStyle:      "One handed",
	}

	person.Body.Head.Brain.Owner = person
	fmt.Printf("%s has been created\n", person.FullName)

	return person
}
// UpdateLocation updates the location of the person
func (p *Person) UpdateLocation(x, y int) {
	p.Location.X = x
	p.Location.Y = y
}

// Grab in the right hand
func (p *Person) GrabRight(item *Item) {
	if p.Body.RightArm.Hand.Items == nil {
		p.Body.RightArm.Hand.Items = []*Item{item}
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
	if p.Body.RightArm.Hand.Items != nil {
		p.Body.RightArm.Hand.Items = nil
	} else {
		fmt.Println("Right hand is empty")
	}
}

// Grab in the left hand
func (p *Person) GrabLeft(item *Item) {
	if p.Body.LeftArm.Hand.Items == nil {
		p.Body.LeftArm.Hand.Items = []*Item{item}
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
	if p.Body.LeftArm.Hand.Items != nil {
		p.Body.LeftArm.Hand.Items = nil
	} else {
		fmt.Println("Left hand is empty")
	}
}

// RemoveLimb removes a limb from the person
func (p *Person) RemoveLimb(limb LimbType) {
	fmt.Println(limb, "has been SEVERED!!!")

	switch limb {
	case "Head":
		p.Body.Head.Brain.turnOff()
		p.Body.Head = nil
		return
	case "RightHand":
		p.Body.RightArm.Hand = nil 
		return
	case "LeftHand":
		p.Body.LeftArm.Hand = nil
		return
	case "RightFoot":
		p.Body.RightLeg.Foot = nil
		p.IsIncapacitated = true
		return
	case "LeftFoot":
		p.Body.LeftLeg.Foot = nil
		p.IsIncapacitated = true
		return
	case "RightLeg":
		p.Body.RightLeg = nil
		p.IsIncapacitated = true
		return
	case "LeftLeg":
		p.Body.LeftLeg = nil
		p.IsIncapacitated = true
		return
	}
}

// AddResidue adds a residue to the limb
func (p *Person) AddResidue(limb LimbType, residue Residue) {

	// Check that the person has the limb
	if p.Body.Head == nil && limb == "Head" {
		fmt.Println("The person has no head to add residue to")
		return
	}
	switch limb {
	case "LeftFoot":
		// Should loop over the residues and add the new residue if it doesn't exist
		for _, r := range p.Body.LeftLeg.Foot.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.LeftLeg.Foot.Residues = append(p.Body.LeftLeg.Foot.Residues, residue)
			}
		}
	case "RightFoot":
		for _, r := range p.Body.RightLeg.Foot.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.RightLeg.Foot.Residues = append(p.Body.RightLeg.Foot.Residues, residue)
			}
		}
	case "Head":
		for _, r := range p.Body.Head.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.Head.Residues = append(p.Body.Head.Residues, residue)
			}
		}
	case "Torso":
		for _, r := range p.Body.Torso.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.Torso.Residues = append(p.Body.Torso.Residues, residue)
			}
		}
	case "RightLeg":
		for _, r := range p.Body.RightLeg.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.RightLeg.Residues = append(p.Body.RightLeg.Residues, residue)
			}
		}
	case "LeftLeg":
		for _, r := range p.Body.LeftLeg.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.LeftLeg.Residues = append(p.Body.LeftLeg.Residues, residue)
			}
		}
	case "RightHand":
		for _, r := range p.Body.RightArm.Hand.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.RightArm.Hand.Residues = append(p.Body.RightArm.Hand.Residues, residue)
			}
		}
	case "LeftHand":
		for _, r := range p.Body.LeftArm.Hand.Residues {
			if r.Name == residue.Name {
				r.Amount += residue.Amount
				return
			} else {
				p.Body.LeftArm.Hand.Residues = append(p.Body.LeftArm.Hand.Residues, residue)
			}
		}
	}
}

func (p *Person) GetVision() Vision {
    return p.WorldProvider.GetVision(p.Location.X, p.Location.Y, p.VisionRange)
}

func (p *Person) GetPersonByFullName(FullName string) *Person {
	return p.WorldProvider.GetPersonByFullName(FullName)
}

func (p *Person) addRelationship(person PersonInVision, relationship string, intensity int) {
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

func (p *Person) UpdateRelationship(fullName string, relationship string, intensity int) {
	for i, rel := range p.Relationships {
		if rel.WithPerson == fullName {
			p.Relationships[i].Relationship = relationship
			p.Relationships[i].Intensity = intensity
			break
		}
	}
}

// ---------------- Finding ----------------------------

// Find Wood - Find wood in the vision
func (p *Person) FindLumberTrees() []*Plant {
	vision := p.GetVision()

	lumberTreesInVision := []*Plant{}

	// Check if there are any Lumber trees in the vision, which are trees called "Oak Tree", for now.
	for _, plant := range vision.Plants {
		if plant.Name == "Oak Tree" {
			lumberTreesInVision = append(lumberTreesInVision, plant)
		}
	}

	if len(lumberTreesInVision) > 0 {
		return lumberTreesInVision
	} else {
		return nil
	}
}

// FindTheClosestPlant - Find the closest plant from a list of plants
func (p *Person) FindTheClosestPlant(plants []*Plant) *Plant {
	closestPlant := plants[0]
		for _, tree := range plants {
			if p.WorldProvider.CalculateDistance(p.Location.X, p.Location.Y, tree.Location.X, tree.Location.Y) < p.WorldProvider.CalculateDistance(p.Location.X, p.Location.Y, closestPlant.Location.X, closestPlant.Location.Y) {
				closestPlant = tree
			}
		}

		return closestPlant
}

// ---------------- Create a new person ----------------

func (w *World) createNewPerson(x, y int) *Person {
    person := NewPerson(w, x, y)
    return person
}
