package main

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

//CreateNewQuadrupedalBody - Returns a new QuadrupedalBody
func CreateNewQuadrupedalBody(Species string) *AnimalBody {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightFrontLeg := CreateNewLeg()
	LeftFrontLeg := CreateNewLeg()
	RightBackLeg := CreateNewLeg()
	LeftBackLeg := CreateNewLeg()
	Tail := CreateNewTail()

	Body := AnimalBody{
		Head:          Head,
		Torso:        Torso,
		RightFrontLeg: RightFrontLeg,
		LeftFrontLeg: LeftFrontLeg,
		RightBackLeg: RightBackLeg,
		LeftBackLeg:  LeftBackLeg,

		Tail:          Tail,
	}
	return &Body
}
		

// CreateNewQuadrupedalBody - Returns a new QuadrupedalBody
func CreateNewAnimal(WorldAccessor WorldAccessor, Species string, x, y int) *Animal {
	Age := rand.Intn(10) + 2
	Gender := gofakeit.Gender()	
	Body := CreateNewQuadrupedalBody(Species)

	Animal := Animal{
		Age:              Age,
		FullName:         Species,
		IsChild:          Age < 18,
		Gender:           Gender,
		Description:      "",
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
		Species:          Species,

		OwnedItems:       []*Item{},

		VisionRange:      8,
		Location:         Location{X: x, Y: y},
		WorldProvider:    WorldAccessor,
		Body:			  Body,

		Strength:         1,
		Agility:          1,
		Intelligence:     1,
		Charisma:         1,
		Stamina:          1,	

		CombatExperience: 1,
		CombatSkill:      1,
	}
	Animal.Body = Body
	return &Animal

}

// CreateNewAnimal - Takes in type of animal and name of animal and returns a new animal
func (w *World) CreateNewAnimalByType(animalType string, x, y int) *Animal {
	switch animalType {
	case "Wolf":
		return CreateNewAnimal(w, "Wolf", x, y)
	default:
	}
	return nil
}
