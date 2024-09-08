package main

func CreateNewLimbStatus() LimbStatus {
	limbStatus := LimbStatus{
		BluntDamage: 0,
		SharpDamage: 0,
		IsBleeding:  false,
		IsBroken:    false,
		Residues:    nil,
		CoveredWith: nil,
		IsAttached:  true,
	}
	return limbStatus
}

func CreateNewHand() *LimbThatCanGrab {
	limbStatus := CreateNewLimbStatus()
	hand := LimbThatCanGrab{
		LimbStatus: limbStatus,
	}
	return &hand
}

func CreateNewArm() *Arm {
	limbStatus := CreateNewLimbStatus()
	arm := Arm{
		LimbThatCantGrab: LimbThatCantGrab{
			LimbStatus: limbStatus,
		},
	}
	arm.Hand = CreateNewHand()
	return &arm
}

func CreateNewFoot() *LimbThatCanMove {
	limbStatus := CreateNewLimbStatus()
	foot := LimbThatCanMove{
		LimbStatus: limbStatus,
	}
	return &foot
}

func CreateNewTail() *LimbThatCanMove {
	limbStatus := CreateNewLimbStatus()
	tail := LimbThatCanMove{
		LimbStatus: limbStatus,
	}
	return &tail
}

func CreateNewLeg() *Leg {
	limbStatus := CreateNewLimbStatus()
	leg := Leg{
		LimbThatCanMove: LimbThatCanMove{
			LimbStatus: limbStatus,
		},
	}
	leg.Foot = CreateNewFoot()
	return &leg
}

func CreateNewBodyPart(name string) *BodyPart {
	bodyPart := BodyPart{
		Name:         name,
		IsBleeding:   false,
		IsBroken:     false,
		IsObstructed: false,
	}
	return &bodyPart
}

func CreateNewHead() *Head {
		limbStatus := CreateNewLimbStatus()
	head := Head{
		LimbStatus: limbStatus,
	}
	head.Eyes = CreateNewBodyPart("Eyes")
	head.Ears = CreateNewBodyPart("Ears")
	head.Nose = CreateNewBodyPart("Nose")
	head.Mouth = CreateNewBodyPart("Mouth")

	return &head
}

func CreateNewTorso() *LimbStatus {
	torso := LimbStatus{
		BluntDamage: 0,
		SharpDamage: 0,
		IsBleeding:  false,
		IsBroken:    false,
		Residues:    nil,
		CoveredWith: nil,
		IsAttached:  true,
	}
	return &torso
}

// Create a new human body and return as Body interface
func CreateBipedalBody() *EntityBody {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightArm := CreateNewArm()
	LeftArm := CreateNewArm()
	RightLeg := CreateNewLeg()
	LeftLeg := CreateNewLeg()

	body := &EntityBody{
		Head:     Head,
		Torso:    Torso,
		RightArm: RightArm,
		LeftArm:  LeftArm,
		RightLeg: RightLeg,
		LeftLeg:  LeftLeg,
	}
	return body // Return as Body interface
}

func CreateBipedalWithTailBody() *EntityBody {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightArm := CreateNewArm()
	LeftArm := CreateNewArm()
	RightLeg := CreateNewLeg()
	LeftLeg := CreateNewLeg()
	Tail := CreateNewTail()

	body := &EntityBody{
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
func CreateQuadrupedalBody() *EntityBody {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightFrontLeg := CreateNewLeg()
	LeftFrontLeg := CreateNewLeg()
	RightBackLeg := CreateNewLeg()
	LeftBackLeg := CreateNewLeg()
	Tail := CreateNewTail()

	body := &EntityBody{
		Head:          Head,
		Torso:         Torso,
		RightFrontLeg: RightFrontLeg,
		LeftFrontLeg:  LeftFrontLeg,
		RightBackLeg:  RightBackLeg,
		LeftBackLeg:   LeftBackLeg,
		Tail:          Tail,
	}
	return body // Return as Body interface
}
