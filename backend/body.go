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
	head.Brain = NewBrain()
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
	head := CreateNewHead()
	torso := CreateNewTorso()
	rightArm := CreateNewArm()
	leftArm := CreateNewArm()
	rightLeg := CreateNewLeg()
	leftLeg := CreateNewLeg()

	body := HumanBody{
		Head:     head,
		Torso:    torso,
		RightArm: rightArm,
		LeftArm:  leftArm,
		RightLeg: rightLeg,
		LeftLeg:  leftLeg,
	}
	return &body
}