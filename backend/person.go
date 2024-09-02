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
		WantsTo:          make([]string, 0),
		FeelingSafe: 	  0,
		FeelingScared:    0,
		Relationships:    []Relationship{},
		Personality:      "",
		Genes:            []string{},

		OwnedItems:       []*Item{},

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
		CombatStyle:      "",
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
	p.Body.RightArm.Hand.Items = append(p.Body.RightArm.Hand.Items, item)
	p.OwnedItems = append(p.OwnedItems, item)
	// If the item has residues, add them to the limb
	if item.Residues != nil {
		for _, residue := range item.Residues {
			p.AddResidue("RightHand", residue)
		}
	}
}

// Grab in the left hand
func (p *Person) GrabLeft(item *Item) {
	p.Body.LeftArm.Hand.Items = append(p.Body.LeftArm.Hand.Items, item)
	p.OwnedItems = append(p.OwnedItems, item)
	// If the item has residues, add them to the limb
	if item.Residues != nil {
		for _, residue := range item.Residues {
			p.AddResidue("RightHand", residue)
		}
	}
}

// Drop from the right hand
func (p *Person) DropRight(item string) {
	for i, heldItem := range p.Body.RightArm.Hand.Items {
		if heldItem.Name == item {
			p.Body.RightArm.Hand.Items = append(p.Body.RightArm.Hand.Items[:i], p.Body.RightArm.Hand.Items[i+1:]...)
			p.WorldProvider.AddItem(p.Location.X, p.Location.Y, heldItem)
			p.Body.Head.Brain.AddMemoryToShortTerm("Dropped my " + item, p.FullName, p.Location)
			return
		}
	}
}

// Drop from the left hand
func (p *Person) DropLeft(item string) {
	for i, heldItem := range p.Body.LeftArm.Hand.Items {
		if heldItem.Name == item {
			p.Body.RightArm.Hand.Items = append(p.Body.RightArm.Hand.Items[:i], p.Body.RightArm.Hand.Items[i+1:]...)
			p.WorldProvider.AddItem(p.Location.X, p.Location.Y, heldItem)
			p.Body.Head.Brain.AddMemoryToShortTerm("Dropped my " + item, p.FullName, p.Location)
			return
		}
	}
}

// Drink - Consume a liquid
func (p *Person) Drink(liquid Liquid) {
    switch liquid.Name {
    case "Water":
        fmt.Println("Drinking water") 
        p.Body.Head.Brain.DecreaseThirstLevel(50)
    }
}

// Eat - Consume food
func (p *Person) Eat(food Food) {
	switch food.GetName() {
	case "Apple":
		fmt.Println("Eating an apple")
		p.Body.Head.Brain.DecreaseHungerLevel(food.GetNutritionalValue())
	}
}

// RemoveLimb removes a limb from the person
func (p *Person) RemoveLimb(limb BodyPartType) {
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
func (p *Person) AddResidue(limb BodyPartType, residue Residue) {

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

// WalkTo - Walk to a location - This assumes that the person can physically walk and that it's possible to walk to the location. This should be one tile away, so one step.
func (p *Person) WalkTo(x, y int) {
	requiredLimbs := []BodyPartType{"RightLeg", "LeftLeg"}
	p.IsMoving = TargetedAction{"Walk", string(x) + ", " + string(y), true, requiredLimbs, 10}
	p.WorldProvider.MovePerson(p, x, y)
}

// ---------------- Finding ----------------------------

// FindLumberTrees - Find wood in the vision
func (p *Person) FindLumberTrees() []*Plant {
	vision := p.WorldProvider.GetPlantsInVision(p.Location.X, p.Location.Y, p.VisionRange)

	lumberTreesInVision := []*Plant{}

	// Check if there are any Lumber trees in the vision, which are trees called "Oak Tree", for now.
	for _, plant := range vision {
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

// FindClosestGrass - Find the closest grass from a list
func (p *Person) FindClosestEmptyGrass(grass []Tile) Tile {
	closestGrass := grass[0]
	for _, tile := range grass {
		if p.IsTileEmpty(tile.Location.X, tile.Location.Y) {
			if p.WorldProvider.CalculateDistance(p.Location.X, p.Location.Y, tile.Location.X, tile.Location.Y) < p.WorldProvider.CalculateDistance(p.Location.X, p.Location.Y, closestGrass.Location.X, closestGrass.Location.Y) {
				closestGrass = tile
			}
		}
	}

	return closestGrass
}

// IsTileEmpty - Check if a tile is empty
func (p *Person) IsTileEmpty(x, y int) bool {
	tile := p.WorldProvider.GetTile(x, y)
	if tile.Shelter == nil && tile.Plant == nil {
		if tile.Person != nil && tile.Person.FullName != p.FullName {
			return false
		}
		return true
	}
	return false
}

//FindClosestWater - Find the closest water from a list of water
func (p *Person) FindClosestWaterSupply(water []Tile) Tile {
	closestWater := water[0]
	for _, tile := range water {
		if p.WorldProvider.CalculateDistance(p.Location.X, p.Location.Y, tile.Location.X, tile.Location.Y) < p.WorldProvider.CalculateDistance(p.Location.X, p.Location.Y, closestWater.Location.X, closestWater.Location.Y) {
			closestWater = tile
		}
	}

	return closestWater
}

// ---------------- Create a new person ----------------

func (w *World) createNewPerson(x, y int) *Person {
    person := NewPerson(w, x, y)
    return person
}
