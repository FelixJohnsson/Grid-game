package main

func CreateNewHand() *LimbThatCanGrab {
	hand := LimbThatCanGrab{
		LimbStatus: LimbStatus{
			BluntDamage: 0,
			SharpDamage: 0,
			IsBleeding:  false,
			IsBroken:    false,
			Residues:    nil,
			CoveredWith: nil,
			IsAttached:  true,
		},
	}
	return &hand
}

func CreateNewArm() *Arm {
	arm := Arm{
		LimbThatCantGrab: LimbThatCantGrab{
			LimbStatus: LimbStatus{
				BluntDamage: 0,
				SharpDamage: 0,
				IsBleeding:  false,
				IsBroken:    false,
				Residues:    nil,
				CoveredWith: nil,
				IsAttached:  true,
			},
		},
	}
	arm.Hand = CreateNewHand()
	return &arm
}

func CreateNewFoot() *LimbThatCanMove {
	foot := LimbThatCanMove{
		LimbStatus: LimbStatus{
			BluntDamage: 0,
			SharpDamage: 0,
			IsBleeding:  false,
			IsBroken:    false,
			Residues:    nil,
			CoveredWith: nil,
			IsAttached:  true,
		},
	}
	return &foot
}

func CreateNewTail() *LimbThatCanMove {
	tail := LimbThatCanMove{
		LimbStatus: LimbStatus{
			BluntDamage: 0,
			SharpDamage: 0,
			IsBleeding:  false,
			IsBroken:    false,
			Residues:    nil,
			CoveredWith: nil,
			IsAttached:  true,
		},
	}
	return &tail
}

func CreateNewLeg() *Leg {
	leg := Leg{
		LimbThatCanMove: LimbThatCanMove{
			LimbStatus: LimbStatus{
				BluntDamage: 0,
				SharpDamage: 0,
				IsBleeding:  false,
				IsBroken:    false,
				Residues:    nil,
				CoveredWith: nil,
				IsAttached:  true,
			},
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
	head := Head{
		LimbStatus: LimbStatus{
			BluntDamage: 0,
			SharpDamage: 0,
			IsBleeding:  false,
			IsBroken:    false,
			Residues:    nil,
			CoveredWith: nil,
			IsAttached:  true,
		},
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

func CreateNewBody() *HumanBody {
	Head := CreateNewHead()
	Torso := CreateNewTorso()
	RightArm := CreateNewArm()
	LeftArm := CreateNewArm()
	RightLeg := CreateNewLeg()
	LeftLeg := CreateNewLeg()

	body := HumanBody{
		Head:     Head,
		Torso:    Torso,
		RightArm: RightArm,
		LeftArm:  LeftArm,
		RightLeg: RightLeg,
		LeftLeg:  LeftLeg,
	}
	return &body
}