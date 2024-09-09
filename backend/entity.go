package main

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

func NewAnimalEntity(worldAccessor WorldAccessor, species string, body *EntityBody, x, y int) *Entity {
	Age := rand.Intn(10) + 2
	FirstName := gofakeit.FirstName()
	FamilyName := gofakeit.LastName()
	Gender := gofakeit.Gender()

	animal := &Entity{
		Age:              Age,
		Title:            "",
		FirstName:        FirstName,
		FamilyName:       FamilyName,
		FullName:         FirstName + " " + FamilyName,
		Initials:         string(FirstName[0]) + string(FamilyName[0]),
		IsChild:          Age < 18,
		Gender:           Gender,
		Occupation:       Unemployed,
		IsMoving:         TargetedAction{},
		IsTalking:        TargetedAction{},
		IsSitting:        TargetedAction{},
		IsEating:         TargetedAction{},
		IsSleeping:       TargetedAction{},
		IsBleeding:       false,
		Thinking:         "",
		WantsTo:          make([]string, 0),
		FeelingSafe: 	  0,
		FeelingScared:    0,
		Relationships:    []Relationship{},
		Genes:            []string{},
		Species:          species,

		OwnedItems:       []*Item{},

		VisionRange:      5,
		Location:         Location{X: x, Y: y},
		WorldProvider:    worldAccessor,
		Body:			  body,
		Brain:            nil,

		Strength:         1,
		Agility:          1,
		Intelligence:     1,
		Charisma:         1,
		Stamina:          1,
		Curiosity:        25,

		CombatExperience: 1,
		CombatSkill:      1,
		CombatStyle:      "",
	}

	animal.Brain = NewBrain(animal)

	return animal
}

func NewPersonEntity(worldAccessor WorldAccessor, x, y int) *Entity {
	Age := rand.Intn(63) + 2
	FirstName := gofakeit.FirstName()
	FamilyName := gofakeit.LastName()
	Gender := gofakeit.Gender()
	Body := CreateBipedalBody()

	person := &Entity{
		Age:              Age,
		Title:            "",
		FirstName:        FirstName,
		FamilyName:       FamilyName,
		FullName:         FirstName + " " + FamilyName,
		Initials:         string(FirstName[0]) + string(FamilyName[0]),
		IsChild:          Age < 18,
		Gender:           Gender,
		Occupation:       Unemployed,
		IsMoving:         TargetedAction{},
		IsTalking:        TargetedAction{},
		IsSitting:        TargetedAction{},
		IsEating:         TargetedAction{},
		IsSleeping:       TargetedAction{},
		IsBleeding:       false,
		Thinking:         "",
		WantsTo:          make([]string, 0),
		FeelingSafe: 	  0,
		FeelingScared:    0,
		Relationships:    []Relationship{},
		Genes:            []string{},
		Species:          "Homo Sapiens",

		OwnedItems:       []*Item{},

		VisionRange:      5,
		Location:         Location{X: x, Y: y},
		WorldProvider:    worldAccessor,
		Body:			  Body,
		Brain:            nil,

		Strength:         1,
		Agility:          1,
		Intelligence:     1,
		Charisma:         1,
		Stamina:          1,
		Curiosity:        25,

		CombatExperience: 1,
		CombatSkill:      1,
		CombatStyle:      "",
	}

	person.Brain = NewBrain(person)

	person.Body = Body
	return person
}

func (e *Entity) GetPersonByFullName(FullName string) *Entity {
	return e.WorldProvider.GetPersonByFullName(FullName)
}

func (e *Entity) AddRelationship(entity EntityInVision, relationship string, intensity int) {
	e.Relationships = append(e.Relationships, Relationship{WithEntity: entity.FullName, Relationship: relationship, Intensity: intensity})
}

func (e *Entity) HasRelationship(fullName string) bool {
	for _, relationship := range e.Relationships {
		if relationship.WithEntity == fullName {
			return true
		}
		return false
	}
	return false
}

func (e *Entity) UpdateRelationship(fullName string, relationship string, intensity int) {
	for i, rel := range e.Relationships {
		if rel.WithEntity == fullName {
			e.Relationships[i].Relationship = relationship
			e.Relationships[i].Intensity = intensity
			break
		}
	}
}

// ---------------- Finding ----------------------------

// FindLumberTrees - Find wood in the vision
func (e *Entity) FindLumberTrees() []*Plant {
	vision := e.WorldProvider.GetPlantsInVision(e.Location.X, e.Location.Y, e.VisionRange)

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
func (e *Entity) FindTheClosestPlant(plants []*Plant) *Plant {
	closestPlant := plants[0]
		for _, tree := range plants {
			if e.WorldProvider.CalculateDistance(e.Location.X, e.Location.Y, tree.Location.X, tree.Location.Y) < e.WorldProvider.CalculateDistance(e.Location.X, e.Location.Y, closestPlant.Location.X, closestPlant.Location.Y) {
				closestPlant = tree
			}
		}

		return closestPlant
}

// FindClosestGrass - Find the closest grass from a list
func (e *Entity) FindClosestEmptyGrass(grass []Tile) Tile {
	closestGrass := grass[0]
	for _, tile := range grass {
		if e.WorldProvider.IsTileEmpty(tile.Location.X, tile.Location.Y) {
			if e.WorldProvider.CalculateDistance(e.Location.X, e.Location.Y, tile.Location.X, tile.Location.Y) < e.WorldProvider.CalculateDistance(e.Location.X, e.Location.Y, closestGrass.Location.X, closestGrass.Location.Y) {
				closestGrass = tile
			}
		}
	}

	return closestGrass
}




// ---------------- Create a new person ----------------

func (w *World) CreateNewPerson(x, y int) *Entity {
    person := NewPersonEntity(w, x, y)

	w.AddEntity(x, y, person)

    return person
}

func (w *World) CreateNewAnimalEntity(species string,x, y int) *Entity {
	switch species {
	case "Canis lupus":
		wolf := NewAnimalEntity(w, species, CreateQuadrupedalBody(), x, y)
		w.AddEntity(3, 3, wolf)
	default:
	}
	return nil
}
