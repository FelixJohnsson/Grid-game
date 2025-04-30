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

func CreateNewHeart() *Heart {
	heart := Heart{
		Vessel: 0,
		IsPumping: false,
	}
	return &heart
}

func CreateNewLungs() *Lungs {
	lungs := Lungs{
		Vessel: 0,
		IsBreathing: false,
		Oxygen: 0,
	}
	return &lungs
}

func CreateNewStomach() *Stomach {
	stomach := Stomach{
		Vessel: 0,
		IsDigesting: false,
		Contains: []Food{},
	}
	return &stomach
}

func CreateNewKidneys() *Kidneys {
	kidneys := Kidneys{
		Vessel: 0,
		IsFiltering: false,
		Toxins: []string{},
	}
	return &kidneys
}


func CreateNewTorso() *Torso {
	limbStatus := CreateNewLimbStatus()
	torso := Torso{
		LimbStatus: limbStatus,
	}
	torso.Heart = CreateNewHeart()
	torso.Lungs = CreateNewLungs()
	torso.Stomach = CreateNewStomach()
	torso.Kidneys = CreateNewKidneys()
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
		Blood: Blood{
			Oxygen: 100,
			Glucose: 100,
			Toxins: 0,
			Water: 100,
			Type: "O-negative",
			Hormones: Hormones{
				Adrenaline: 1,   // Very low
				Cortisol: 15,     // Moderate
				Dopamine: 25,     // Moderate
				Epinephrine: 1,   // Very low
				Endorphin: 10,    // Low-moderate
				Serotonin: 30,    // High
			},
			Amount: 100,
		},
	}
	return body
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
		Blood: Blood{
			Oxygen: 100,
			Glucose: 100,
			Toxins: 0,
			Water: 100,
			Type: "O-negative",
			Hormones: Hormones{
				Adrenaline: 1,   // Very low
				Cortisol: 15,     // Moderate
				Dopamine: 25,     // Moderate
				Epinephrine: 1,   // Very low
				Endorphin: 10,    // Low-moderate
				Serotonin: 30,    // High
			},
			Amount: 100,
		},
	}
	return body
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
	return body
}
