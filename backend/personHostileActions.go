package main

import (
	"fmt"
	"math/rand"
)

// ---------------- Hostile actions ----------------

// CalculateLegDamageGiven calculates the damage given by the leg
func (p *Person) CalculateLegDamageGiven(target *Person, targetLimb BodyPartType, withLimb *Leg) Damage {
	// Calculate the damage based on limb status, physical attributes and experience.
	damage := Damage{}

	if withLimb.IsBroken || withLimb == nil {
		return damage
	} else {
		damage.AmountBluntDamage = 1 + p.Strength
	}

	// Add a random factor to the damage
	damage.AmountBluntDamage += rand.Intn(3)

	// Calculate the defense of the target
	defense := target.CalculateDefense(targetLimb)

	// Apply the defense to the damage
	damage.AmountBluntDamage -= defense

	if damage.AmountBluntDamage < 0 {
		damage.AmountBluntDamage = 0
	}
	if damage.AmountSharpDamage < 0 {
		damage.AmountSharpDamage = 0
	}

	return damage
}



// CalculateArmDamageGiven calculates the damage given by the arm
func (p *Person) CalculateArmDamageGiven(target *Person, targetLimb BodyPartType, withLimb *LimbThatCanGrab) Damage {
		// Calculate the damage based on limb status, item in hand, physical attributes and experience.
		damage := Damage{}

		if withLimb.IsBroken || withLimb == nil {
			return damage
		} else if withLimb.Items == nil {
				damage.AmountBluntDamage = 1 + p.Strength
				damage.AmountSharpDamage = 1 + p.Strength
		} else {
			damage.AmountBluntDamage = withLimb.Items[0].Bluntness
			damage.AmountSharpDamage = withLimb.Items[0].Sharpness

			damage.AmountBluntDamage += p.Strength
			damage.AmountSharpDamage += p.Strength

			damage.AmountBluntDamage += p.CombatSkill + p.CombatExperience
			damage.AmountSharpDamage += p.CombatSkill + p.CombatExperience
		}
	// Add a random factor to the damage
	damage.AmountBluntDamage += rand.Intn(3)
	damage.AmountSharpDamage += rand.Intn(3)

	// Calculate the defense of the target
	defense := target.CalculateDefense(targetLimb)

	// Apply the defense to the damage
	damage.AmountBluntDamage -= defense
	damage.AmountSharpDamage -= defense

	if damage.AmountBluntDamage < 0 {
		damage.AmountBluntDamage = 0
	}
	if damage.AmountSharpDamage < 0 {
		damage.AmountSharpDamage = 0
	}

	return damage
}

func (p *Person) CalculateDefense(targetLimb BodyPartType) int {
	// Calculate the defense based on limb status, item in hand, physical attributes and experience.
	defense := 0

	defense += p.Strength + p.Agility
	defense += p.CombatSkill + p.CombatExperience

	return defense
}
	
// Flee is called when the person is feeling scared
func (p *Person) Flee(attacker *Person) {
	// Move away from the attacker
	
}
// AttackWithLeg - target is the person being attacked, targetLimb is the limb being attacked, withLimb (probably leg) is the limb that is attacking
func (p *Person) AttackWithLeg(target *Person, targetLimb BodyPartType, withLimb *Leg) Damage {
	if target == nil {
		fmt.Println("No target to attack")
		return Damage{}
	}

	damage := p.CalculateLegDamageGiven(target, targetLimb, withLimb)
	target.ReceivingApplyDamageTo(targetLimb, damage)

	return damage
}
		
// AttackWithArm - target is the person being attacked, targetLimb is the limb being attacked, withLimb (probably hand) is the limb that is attacking
func (p *Person) AttackWithArm(target *Person, targetLimb BodyPartType, withLimb *LimbThatCanGrab) Damage {
	if target == nil {
		fmt.Println("No target to attack")
		return Damage{}
	} 
	if !target.Brain.IsConscious {
		fmt.Println("The target is unconscious, decide what to do next")
		return Damage{}
	}

	damage := p.CalculateArmDamageGiven(target, targetLimb, withLimb)
	fmt.Println(p.FullName, "is attacking", target.FullName, "and causing", damage.AmountBluntDamage, "blunt damage and", damage.AmountSharpDamage, "sharp damage")
	target.ReceivingApplyDamageTo(targetLimb, damage)

	if target.Body.Head != nil {
		target.Brain.IsUnderAttack = IsUnderAttack{true, p, targetLimb, "RightHand"}
	}

	return damage
}

// Logic for being attacked, if the hands are broken, drop the items in the hands
func (p *Person) ReceivingApplyDamageTo(limb BodyPartType, damage Damage) {
	// TODO: Check if the limb is covered with a wearable that can protect the limb from the damage

	var bluntDamageUntilBroken = 50
	var bluntDamageUntilUnconscious = 75
	var bluntDamageUntilTorsoIncapacitated = 75
	var brainDamageUntilDead = 100

	var sharpDamageUntilBleeding = 10
	var sharpDamageUntilSevered = 50

	// Add blood residue to the limb and item
	bloodResidue := Residue{"Blood", 1}
	p.AddResidue(limb, bloodResidue)

	switch limb {
		case "Head":
			if p.Body.Head == nil {
				fmt.Println("The target has no head to attack")
				return
			}
		p.Body.Head.BluntDamage += damage.AmountBluntDamage
		p.Body.Head.SharpDamage += damage.AmountSharpDamage

		p.Brain.BrainDamage += damage.AmountBluntDamage

		if p.Body.Head.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("Head")
			bloodResidue := Residue{"Blood", 10}
			p.AddResidue("Torso", bloodResidue)
			return
		}
		if p.Body.Head.BluntDamage > bluntDamageUntilBroken {
			p.Body.Head.IsBroken = true
			if p.Body.Head.BluntDamage > bluntDamageUntilUnconscious {
				p.Brain.IsConscious = false
				p.IsIncapacitated = true
			}
			if p.Body.Head.BluntDamage >= brainDamageUntilDead && p.Brain.Active {
				p.Brain.turnOff()
			}
		}
		if p.Body.Head.SharpDamage > sharpDamageUntilBleeding {
			p.Body.Head.IsBleeding = true
		}
		return
	case "LeftFoot":
		if p.Body.LeftLeg.Foot == nil {
			fmt.Println("The target has no left foot to attack")
			return
		}
		p.Body.LeftLeg.Foot.BluntDamage += damage.AmountBluntDamage
		p.Body.LeftLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if p.Body.LeftLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			p.Body.LeftLeg.Foot.IsBroken = true
			p.IsIncapacitated = true
		} 
		if p.Body.LeftLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			p.Body.LeftLeg.Foot.IsBleeding = true
		}
		if p.Body.LeftLeg.Foot.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("LeftFoot")
		}
		return
	case "RightFoot":
		if p.Body.RightLeg.Foot == nil {
			fmt.Println("The target has no right foot to attack")
			return
		}
		p.Body.RightLeg.Foot.BluntDamage += damage.AmountBluntDamage
		p.Body.RightLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if p.Body.RightLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			p.Body.RightLeg.Foot.IsBroken = true
			p.IsIncapacitated = true
		}
		if p.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			p.Body.RightLeg.Foot.IsBleeding = true
		}
		if p.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("RightFoot")
		}
		return
	case "Torso":
		p.Body.Torso.BluntDamage += damage.AmountBluntDamage
		p.Body.Torso.SharpDamage += damage.AmountSharpDamage

		if p.Body.Torso.BluntDamage > bluntDamageUntilBroken {
			p.Body.Torso.IsBroken = true
			if p.Body.Torso.BluntDamage > bluntDamageUntilTorsoIncapacitated {
				p.IsIncapacitated = true
			}
		}
		if p.Body.Torso.SharpDamage > sharpDamageUntilBleeding {
			p.Body.Torso.IsBleeding = true
		}
		return
	case "RightLeg":
		if p.Body.RightLeg == nil {
			fmt.Println("The target has no right leg to attack")
			return
		}
		p.Body.RightLeg.Foot.BluntDamage += damage.AmountBluntDamage
		p.Body.RightLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if p.Body.RightLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			p.Body.RightLeg.Foot.IsBroken = true
			p.IsIncapacitated = true
		}
		if p.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			p.Body.RightLeg.Foot.IsBleeding = true
		}
		if p.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("RightLeg")
		}
		return 
	case "LeftLeg":
		if p.Body.LeftLeg == nil {
			fmt.Println("The target has no left leg to attack")
			return
		}
		p.Body.LeftLeg.Foot.BluntDamage += damage.AmountBluntDamage
		p.Body.LeftLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if p.Body.LeftLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			p.Body.LeftLeg.Foot.IsBroken = true
			p.IsIncapacitated = true
		}
		if p.Body.LeftLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			p.Body.LeftLeg.Foot.IsBleeding = true
		}
		if p.Body.LeftLeg.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("LeftLeg")
		}
		return 
	case "RightHand":
		if p.Body.RightArm.Hand == nil {
			fmt.Println("The target has no right hand to attack")
			return
		}
		p.Body.RightArm.Hand.BluntDamage += damage.AmountBluntDamage
		p.Body.RightArm.Hand.SharpDamage += damage.AmountSharpDamage

		if p.Body.RightArm.Hand.BluntDamage > bluntDamageUntilBroken {
			p.Body.RightArm.Hand.IsBroken = true
			p.Body.RightArm.Hand.Items = nil
		}
		if p.Body.RightArm.Hand.SharpDamage > sharpDamageUntilBleeding {
			p.Body.RightArm.Hand.IsBleeding = true
		}
		if p.Body.RightArm.Hand.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("RightHand")
		}
		return 
	case "LeftHand":
		if p.Body.LeftArm.Hand == nil {
			fmt.Println("The target has no left hand to attack")
			return
		}
		p.Body.LeftArm.Hand.BluntDamage += damage.AmountBluntDamage
		p.Body.LeftArm.Hand.SharpDamage += damage.AmountSharpDamage

		if p.Body.LeftArm.Hand.BluntDamage > bluntDamageUntilBroken {
			p.Body.LeftArm.Hand.IsBroken = true
			p.Body.LeftArm.Hand.Items = nil
		}
		if p.Body.LeftArm.Hand.SharpDamage > sharpDamageUntilBleeding {
			p.Body.LeftArm.Hand.IsBleeding = true
		}
		if p.Body.LeftArm.Hand.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("LeftHand")
		}
	}
		
}