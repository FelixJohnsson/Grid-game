package main

import (
	"fmt"
	"math/rand"
)

// ---------------- Hostile actions ----------------

// CalculateLegDamageGiven calculates the damage given by the leg
func (e *Entity) CalculateLegDamageGiven(target *Entity, targetLimb BodyPartType, withLimb *Leg) Damage {
	// Calculate the damage based on limb status, physical attributes and experience.
	damage := Damage{}

	if withLimb.IsBroken || withLimb == nil {
		return damage
	} else {
		damage.AmountBluntDamage = 1 + e.Strength
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
func (e *Entity) CalculateArmDamageGiven(target *Entity, targetLimb BodyPartType, withLimb *LimbThatCanGrab) Damage {
		// Calculate the damage based on limb status, item in hand, physical attributes and experience.
		damage := Damage{}

		if withLimb.IsBroken || withLimb == nil {
			return damage
		} else if withLimb.Items == nil {
				damage.AmountBluntDamage = 1 + e.Strength
				damage.AmountSharpDamage = 1 + e.Strength
		} else {
			damage.AmountBluntDamage = withLimb.Items[0].Bluntness
			damage.AmountSharpDamage = withLimb.Items[0].Sharpness

			damage.AmountBluntDamage += e.Strength
			damage.AmountSharpDamage += e.Strength

			damage.AmountBluntDamage += e.CombatSkill + e.CombatExperience
			damage.AmountSharpDamage += e.CombatSkill + e.CombatExperience
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

func (e *Entity) CalculateDefense(targetLimb BodyPartType) int {
	// Calculate the defense based on limb status, item in hand, physical attributes and experience.
	defense := 0

	defense += e.Strength + e.Agility
	defense += e.CombatSkill + e.CombatExperience

	return defense
}
	
// Flee is called when the Entity is feeling scared
func (e *Entity) Flee(attacker *Entity) {
	// Move away from the attacker
	
}
// AttackWithLeg - target is the Entity being attacked, targetLimb is the limb being attacked, withLimb (probably leg) is the limb that is attacking
func (e *Entity) AttackWithLeg(target *Entity, targetLimb BodyPartType, withLimb *Leg) Damage {
	if target == nil {
		fmt.Println("No target to attack")
		return Damage{}
	}

	damage := e.CalculateLegDamageGiven(target, targetLimb, withLimb)
	target.ReceivingApplyDamageTo(targetLimb, damage)

	return damage
}
		
// AttackWithArm - target is the Entity being attacked, targetLimb is the limb being attacked, withLimb (probably hand) is the limb that is attacking
func (e *Entity) AttackWithArm(target *Entity, targetLimb BodyPartType, withLimb *LimbThatCanGrab) Damage {
	if target == nil {
		fmt.Println("No target to attack")
		return Damage{}
	} 
	if !target.Brain.IsConscious {
		fmt.Println("The target is unconscious, decide what to do next")
		return Damage{}
	}

	damage := e.CalculateArmDamageGiven(target, targetLimb, withLimb)
	fmt.Println(e.FullName, "is attacking", target.FullName, "and causing", damage.AmountBluntDamage, "blunt damage and", damage.AmountSharpDamage, "sharp damage")
	target.ReceivingApplyDamageTo(targetLimb, damage)

	if target.Body.Head != nil {
		target.Brain.IsUnderAttack = IsUnderAttack{true, e, targetLimb, "RightHand"}
	}

	return damage
}

// Logic for being attacked, if the hands are broken, drop the items in the hands
func (e *Entity) ReceivingApplyDamageTo(limb BodyPartType, damage Damage) {
	// TODO: Check if the limb is covered with a wearable that can protect the limb from the damage

	var bluntDamageUntilBroken = 50
	var bluntDamageUntilUnconscious = 75
	var bluntDamageUntilTorsoIncapacitated = 75
	var brainDamageUntilDead = 100

	var sharpDamageUntilBleeding = 10
	var sharpDamageUntilSevered = 50

	// Add blood residue to the limb and item
	bloodResidue := Residue{"Blood", 1}
	e.AddResidue(limb, bloodResidue)

	switch limb {
		case "Head":
			if e.Body.Head == nil {
				fmt.Println("The target has no head to attack")
				return
			}
		e.Body.Head.BluntDamage += damage.AmountBluntDamage
		e.Body.Head.SharpDamage += damage.AmountSharpDamage

		e.Brain.BrainDamage += damage.AmountBluntDamage

		if e.Body.Head.SharpDamage > sharpDamageUntilSevered {
			e.RemoveLimb("Head")
			bloodResidue := Residue{"Blood", 10}
			e.AddResidue("Torso", bloodResidue)
			return
		}
		if e.Body.Head.BluntDamage > bluntDamageUntilBroken {
			e.Body.Head.IsBroken = true
			if e.Body.Head.BluntDamage > bluntDamageUntilUnconscious {
				e.Brain.IsConscious = false
				e.IsIncapacitated = true
			}
			if e.Body.Head.BluntDamage >= brainDamageUntilDead && e.Brain.Active {
				e.Brain.turnOff()
			}
		}
		if e.Body.Head.SharpDamage > sharpDamageUntilBleeding {
			e.Body.Head.IsBleeding = true
		}
		return
	case "LeftFoot":
		if e.Body.LeftLeg.Foot == nil {
			fmt.Println("The target has no left foot to attack")
			return
		}
		e.Body.LeftLeg.Foot.BluntDamage += damage.AmountBluntDamage
		e.Body.LeftLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if e.Body.LeftLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			e.Body.LeftLeg.Foot.IsBroken = true
			e.IsIncapacitated = true
		} 
		if e.Body.LeftLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			e.Body.LeftLeg.Foot.IsBleeding = true
		}
		if e.Body.LeftLeg.Foot.SharpDamage > sharpDamageUntilSevered {
			e.RemoveLimb("LeftFoot")
		}
		return
	case "RightFoot":
		if e.Body.RightLeg.Foot == nil {
			fmt.Println("The target has no right foot to attack")
			return
		}
		e.Body.RightLeg.Foot.BluntDamage += damage.AmountBluntDamage
		e.Body.RightLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if e.Body.RightLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			e.Body.RightLeg.Foot.IsBroken = true
			e.IsIncapacitated = true
		}
		if e.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			e.Body.RightLeg.Foot.IsBleeding = true
		}
		if e.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilSevered {
			e.RemoveLimb("RightFoot")
		}
		return
	case "Torso":
		e.Body.Torso.BluntDamage += damage.AmountBluntDamage
		e.Body.Torso.SharpDamage += damage.AmountSharpDamage

		if e.Body.Torso.BluntDamage > bluntDamageUntilBroken {
			e.Body.Torso.IsBroken = true
			if e.Body.Torso.BluntDamage > bluntDamageUntilTorsoIncapacitated {
				e.IsIncapacitated = true
			}
		}
		if e.Body.Torso.SharpDamage > sharpDamageUntilBleeding {
			e.Body.Torso.IsBleeding = true
		}
		return
	case "RightLeg":
		if e.Body.RightLeg == nil {
			fmt.Println("The target has no right leg to attack")
			return
		}
		e.Body.RightLeg.Foot.BluntDamage += damage.AmountBluntDamage
		e.Body.RightLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if e.Body.RightLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			e.Body.RightLeg.Foot.IsBroken = true
			e.IsIncapacitated = true
		}
		if e.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			e.Body.RightLeg.Foot.IsBleeding = true
		}
		if e.Body.RightLeg.Foot.SharpDamage > sharpDamageUntilSevered {
			e.RemoveLimb("RightLeg")
		}
		return 
	case "LeftLeg":
		if e.Body.LeftLeg == nil {
			fmt.Println("The target has no left leg to attack")
			return
		}
		e.Body.LeftLeg.Foot.BluntDamage += damage.AmountBluntDamage
		e.Body.LeftLeg.Foot.SharpDamage += damage.AmountSharpDamage

		if e.Body.LeftLeg.Foot.BluntDamage > bluntDamageUntilBroken {
			e.Body.LeftLeg.Foot.IsBroken = true
			e.IsIncapacitated = true
		}
		if e.Body.LeftLeg.Foot.SharpDamage > sharpDamageUntilBleeding {
			e.Body.LeftLeg.Foot.IsBleeding = true
		}
		if e.Body.LeftLeg.SharpDamage > sharpDamageUntilSevered {
			e.RemoveLimb("LeftLeg")
		}
		return 
	case "RightHand":
		if e.Body.RightArm.Hand == nil {
			fmt.Println("The target has no right hand to attack")
			return
		}
		e.Body.RightArm.Hand.BluntDamage += damage.AmountBluntDamage
		e.Body.RightArm.Hand.SharpDamage += damage.AmountSharpDamage

		if e.Body.RightArm.Hand.BluntDamage > bluntDamageUntilBroken {
			e.Body.RightArm.Hand.IsBroken = true
			e.Body.RightArm.Hand.Items = nil
		}
		if e.Body.RightArm.Hand.SharpDamage > sharpDamageUntilBleeding {
			e.Body.RightArm.Hand.IsBleeding = true
		}
		if e.Body.RightArm.Hand.SharpDamage > sharpDamageUntilSevered {
			e.RemoveLimb("RightHand")
		}
		return 
	case "LeftHand":
		if e.Body.LeftArm.Hand == nil {
			fmt.Println("The target has no left hand to attack")
			return
		}
		e.Body.LeftArm.Hand.BluntDamage += damage.AmountBluntDamage
		e.Body.LeftArm.Hand.SharpDamage += damage.AmountSharpDamage

		if e.Body.LeftArm.Hand.BluntDamage > bluntDamageUntilBroken {
			e.Body.LeftArm.Hand.IsBroken = true
			e.Body.LeftArm.Hand.Items = nil
		}
		if e.Body.LeftArm.Hand.SharpDamage > sharpDamageUntilBleeding {
			e.Body.LeftArm.Hand.IsBleeding = true
		}
		if e.Body.LeftArm.Hand.SharpDamage > sharpDamageUntilSevered {
			e.RemoveLimb("LeftHand")
		}
	}
		
}