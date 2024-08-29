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
		WantsTo:          "",
		FeelingSafe: 	  0,
		FeelingScared:    0,
		Relationships:    []Relationship{},
		Personality:      "",
		Genes:            []string{},

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
		CombatStyle:      "One handed",
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
	if p.Body.RightArm.Hand.Items == nil {
		p.Body.RightArm.Hand.Items = []*Item{item}
		// If the item has residues, add them to the limb
		if item.Residues != nil {
			for _, residue := range item.Residues {
				p.AddResidue("RightHand", residue)
			}
		}
	} else {
		fmt.Println("Right hand is already holding something")
	}
}

// Drop from the right hand
func (p *Person) DropRight() {
	if p.Body.RightArm.Hand.Items != nil {
		p.Body.RightArm.Hand.Items = nil
	} else {
		fmt.Println("Right hand is empty")
	}
}

// Grab in the left hand
func (p *Person) GrabLeft(item *Item) {
	if p.Body.LeftArm.Hand.Items == nil {
		p.Body.LeftArm.Hand.Items = []*Item{item}
		// If the item has residues, add them to the limb
		if item.Residues != nil {
			for _, residue := range item.Residues {
				p.AddResidue("RightHand", residue)
			}
		}
	} else {
		fmt.Println("Left hand is already holding something")
	}
}

// Drop from the left hand
func (p *Person) DropLeft() {
	if p.Body.LeftArm.Hand.Items != nil {
		p.Body.LeftArm.Hand.Items = nil
	} else {
		fmt.Println("Left hand is empty")
	}
}

// RemoveLimb removes a limb from the person
func (p *Person) RemoveLimb(limb LimbType) {
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
func (p *Person) AddResidue(limb LimbType, residue Residue) {
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

func (p *Person) GetVision() Vision {
    return p.WorldProvider.GetVision(p.Location.X, p.Location.Y, p.VisionRange)
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

func (p *Person) removeRelationship(person Person) {
	for i, relationship := range p.Relationships {
		if relationship.WithPerson == person.FullName {
			p.Relationships = append(p.Relationships[:i], p.Relationships[i+1:]...)
			break
		}
	}
}

func (p *Person) updateRelationship(fullName string, relationship string, intensity int) {
	for i, rel := range p.Relationships {
		if rel.WithPerson == fullName {
			p.Relationships[i].Relationship = relationship
			p.Relationships[i].Intensity = intensity
			break
		}
	}
}

// ---------------- Hostile actions ----------------

func (p *Person) CalculateDamageGiven(target *Person, targetLimb LimbType, withLimb *LimbThatCanGrab) Damage {
		// Calculate the damage based on limb status, item in hand, physical attributes and experience.
		damage := Damage{}

		if withLimb.IsBroken {
			return damage
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

	return damage
}

func (p *Person) CalculateDefense (targetLimb LimbType) int {
	// Calculate the defense based on limb status, item in hand, physical attributes and experience.
	defense := 0

	defense += p.Strength + p.Agility
	defense += p.CombatSkill + p.CombatExperience

	return defense
}
		

func (p *Person) AttackWithWeapon(target *Person, targetLimb LimbType, withLimb *LimbThatCanGrab) Damage {
	if target == nil {
		fmt.Println("No target to attack")
		return Damage{}
	}

	damage := p.CalculateDamageGiven(target, targetLimb, withLimb)
	target.ReceivingApplyDamageTo(targetLimb, damage)

	return damage
	// TODO: This should also be sent to the brain's Being Attacked function
}

// Logic for being attacked, if the hands are broken, drop the items in the hands
func (p *Person) ReceivingApplyDamageTo(limb LimbType, damage Damage) {
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

		p.Body.Head.Brain.BrainDamage += damage.AmountBluntDamage

		if p.Body.Head.SharpDamage > sharpDamageUntilSevered {
			p.RemoveLimb("Head")
			bloodResidue := Residue{"Blood", 10}
			p.AddResidue("Torso", bloodResidue)
			return
		}
		if p.Body.Head.BluntDamage > bluntDamageUntilBroken {
			p.Body.Head.IsBroken = true
			if p.Body.Head.BluntDamage > bluntDamageUntilUnconscious {
				p.Body.Head.Brain.IsConscious = false
				p.IsIncapacitated = true
			}
			if p.Body.Head.BluntDamage >= brainDamageUntilDead && p.Body.Head.Brain.Active {
				p.Body.Head.Brain.turnOff()
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

// ---------------- Create a new person ----------------

func (w *World) createNewPerson(x, y int) *Person {
    person := NewPerson(w, x, y)
    return person
}
