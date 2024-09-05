package main

// Define a common Body interface
type Body interface {
	HasTail() bool
}

type QuadrupedalBody struct {
	Head          *Head
	Torso         *LimbStatus
	RightFrontLeg *Leg
	LeftFrontLeg  *Leg
	RightBackLeg  *Leg
	LeftBackLeg   *Leg
}

// Create a new human body and return as Body interface
func CreateBipedalBody() Body {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightArm := CreateNewArm()
	LeftArm := CreateNewArm()
	RightLeg := CreateNewLeg()
	LeftLeg := CreateNewLeg()

	body := &BipedalBody{
		Head:     Head,
		Torso:    Torso,
		RightArm: RightArm,
		LeftArm:  LeftArm,
		RightLeg: RightLeg,
		LeftLeg:  LeftLeg,
	}
	return body // Return as Body interface
}

func CreateBipedalWithTailBody() Body {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightArm := CreateNewArm()
	LeftArm := CreateNewArm()
	RightLeg := CreateNewLeg()
	LeftLeg := CreateNewLeg()
	Tail := CreateNewTail()

	body := &BipedalWithTailBody{
		Head:     Head,
		Torso:    Torso,
		RightArm: RightArm,
		LeftArm:  LeftArm,
		RightLeg: RightLeg,
		LeftLeg:  LeftLeg,
		Tail:     Tail,
	}
	return body // Return as Body interface
}

// Create a new horse body and return as Body interface
func CreateQuadrupedalBody() Body {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightFrontLeg := CreateNewLeg()
	LeftFrontLeg := CreateNewLeg()
	RightBackLeg := CreateNewLeg()
	LeftBackLeg := CreateNewLeg()

	body := &QuadrupedalBody{
		Head:          Head,
		Torso:         Torso,
		RightFrontLeg: RightFrontLeg,
		LeftFrontLeg:  LeftFrontLeg,
		RightBackLeg:  RightBackLeg,
		LeftBackLeg:   LeftBackLeg,
	}
	return body // Return as Body interface
}

// HasTail - Check if the body has a tail
func (Q *QuadrupedalBody) HasTail() bool {
	return true
}

// HasTail - Check if the body has a tail
func (B *BipedalBody) HasTail() bool {
	return false
}

// HasTail - Check if the body has a tail
func (B *BipedalWithTailBody) HasTail() bool {
	return true
}

