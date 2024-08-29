package main

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

type WorldState struct {
	Map 	 	[][]string;
  	Buildings	[]Building;
  	Persons	 	[]Person;
  	Resources 	Resources;
}

// Enum for Jobs
type Jobs string

const (
	Farmer     Jobs = "Farmer"
	Miner      Jobs = "Miner"
	Lumberjack Jobs = "Lumberjack"
	Builder    Jobs = "Builder"
	Soldier    Jobs = "Soldier"
	Unemployed Jobs = "Unemployed"
)
type Relationship struct {
	WithPerson string
	Relationship string
	Intensity int
}

type TargetedAction struct {
	Action string
	Target string
	IsActive bool
	RequiresLimb []string
}

type Wearable struct {
	Name string
	Material string
	Protection int
}

// Body status
type LimbStatus struct {
	BluntDamage int
	SharpDamage int
	IsBleeding bool
	IsBroken bool
	Residues []Residue
	CoveredWith []Wearable
	IsAttached bool
}

type LimbThatCanHold struct {
	LimbStatus
	Items []*Item
	WeightOfItems int
}

type Damage struct {
	AmountBluntDamage int
	AmountSharpDamage int
}

type Head struct {
	LimbStatus
	Brain *Brain
}

type LimbThatCanGrab struct {
	LimbStatus
	Items []*Item
	WeightOfItems int
}

type LimbThatCantGrab struct {
	LimbStatus
}

type LimbThatCanMove struct {
	LimbStatus
}

type Leg struct {
	LimbThatCanMove
	Foot *LimbThatCanMove
}

type Arm struct {
	LimbThatCantGrab
	Hand *LimbThatCanGrab
}

type HumanBody struct {
	Head *Head
	Torso *LimbStatus
	RightArm *Arm
	LeftArm *Arm
	RightLeg *Leg
	LeftLeg *Leg
}

type Person struct {
	Age              int
	Title 		     string
	FirstName        string
	FamilyName       string
	FullName 	     string
	Initials         string
	IsChild          bool
	Gender           string
	Description      string
	Icon             string
	Occupation       Jobs
	IsWorkingAt      *Building
	Color            string
	Personality 	 string
	Genes            []string

	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsHolding        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction
	IsWorking        TargetedAction

	Thinking         string
	WantsTo          string
	FeelingSafe 	 int
	FeelingScared	 int

	Body 		     *HumanBody

	Strength         int
	Agility          int
	Intelligence     int
	Charisma         int
	Stamina          int

	CombatExperience int
	CombatSkill      int
	CombatStyle      string

	Relationships    []Relationship

	IsIncapacitated  bool
	VisionRange 	 int
	WorldProvider    WorldAccessor
	Location         Location
	OnTileType 	     TileType
}

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
		Icon:             "P",
		Occupation:       Unemployed,
		IsWorkingAt:      nil,
		Color:            "",
		IsMoving:         TargetedAction{},
		IsTalking:        TargetedAction{},
		IsSitting:        TargetedAction{},
		IsHolding:        TargetedAction{},
		IsEating:         TargetedAction{},
		IsSleeping:       TargetedAction{},
		IsWorking:        TargetedAction{},
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

	person.Body.Head.Brain.owner = person
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
				p.IsHolding.IsActive = true
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
	case "RightHand":
		p.Body.RightArm.Hand = nil 
		return
	case "LeftHand":
		p.Body.LeftArm.Hand = nil
		return
	case "RightFoot":
		p.Body.RightLeg.Foot = nil
		return
	case "LeftFoot":
		p.Body.LeftLeg.Foot = nil
		return
	case "RightLeg":
		p.Body.RightLeg = nil
		return
	case "LeftLeg":
		p.Body.LeftLeg = nil
		return
	case "Head":
		p.Body.Head = nil
		return
	}
}

// A type for all the limbs instead of "string"
type LimbType string

const (
	RightHand LimbType = "RightHand"
	LeftHand  LimbType = "LeftHand"
	RightFoot LimbType = "RightFoot"
	LeftFoot  LimbType = "LeftFoot"
	RightLeg  LimbType = "RightLeg"
	LeftLeg   LimbType = "LeftLeg"
	TheHead   LimbType = "Head"
	Torso     LimbType = "Torso"
)

// AddResidue adds a residue to the limb
func (p *Person) AddResidue(limb LimbType, residue Residue) {
	switch limb {
	case "LeftFoot":
		p.Body.LeftLeg.Foot.Residues = append(p.Body.LeftLeg.Foot.Residues, residue)
	case "RightFoot":
		p.Body.RightLeg.Foot.Residues = append(p.Body.RightLeg.Foot.Residues, residue)
	case "Head":
		p.Body.Head.Residues = append(p.Body.Head.Residues, residue)
	case "Torso":
		p.Body.Torso.Residues = append(p.Body.Torso.Residues, residue)
	case "RightLeg":
		p.Body.RightLeg.Residues = append(p.Body.RightLeg.Residues, residue)
	case "LeftLeg":
		p.Body.LeftLeg.Residues = append(p.Body.LeftLeg.Residues, residue)
	case "RightHand":
		p.Body.RightArm.Hand.Residues = append(p.Body.RightArm.Hand.Residues, residue)
	case "LeftHand":
		p.Body.LeftArm.Hand.Residues = append(p.Body.LeftArm.Hand.Residues, residue)
	}
}

func (p *Person) GetVision() Vision {
    return p.WorldProvider.GetVision(p.Location.X, p.Location.Y, p.VisionRange)
}

func (p *Person) GetPersonByFullName(FullName string) *Person {
	return p.WorldProvider.GetPersonByFullName(FullName)
}

func (p *Person) addRelationship(person PersonCleaned, relationship string, intensity int) {
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

func (p *Person) CalculateDamageGiven (target *Person, targetLimb string) Damage {
		// Calculate the damage based on limb status, item in hand, physical attributes and experience.
		damage := Damage{}

		if p.Body.RightArm.Hand.IsBroken {
			return damage
		} else {
			damage.AmountBluntDamage = p.Body.RightArm.Hand.Items[0].Bluntness
			damage.AmountSharpDamage = p.Body.RightArm.Hand.Items[0].Sharpness

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

func (p *Person) CalculateDefense (targetLimb string) int {
	// Calculate the defense based on limb status, item in hand, physical attributes and experience.
	defense := 0

	defense += p.Strength + p.Agility
	defense += p.CombatSkill + p.CombatExperience

	return defense
}
		

func (p *Person) Attack(target *Person, targetLimb string) {
	if target == nil {
		fmt.Println("No target to attack")
		return
	}
	// Logic for attacking
	// 1. Calculate the damage based on limb status, item in hand, physical attributes and experience.
	// 2. Apply the damage to the target's limb depending on the attack style and target's combat experience, skill and physical attributes.
	// 3. If the target limb breaks, the target drops the item that was in the limb, if the limb was holding an item.
	// 4. If the target limb starts bleeding, apply the bleeding effect to the target.
	// 5. The attack can either be blunt or sharp, depending on the item in the attacker's hand.
	// 6. The attack can be blocked by the target if the attack isnt a surprise attack.

	// 1. Calculate the damage
	damage := p.CalculateDamageGiven(target, targetLimb)

	// 2. Apply the damage
	target.ApplyDamageTo(targetLimb, damage)

	// TODO: This should also be sent to the brain's Being Attacked function
}

// Logic for being attacked, if the hands are broken, drop the items in the hands
func (p *Person) ApplyDamageTo(limb string, damage Damage) {
	// Apply the damage to the limb

	// TODO: Check if the limb is covered with a wearable that can protect the limb from the damage
	// TODO: Sharp damage can sever the limb if the sharp damage is high enough, then the limb should be removed from the person

	var bluntDamageUntilBroken = 50
	var bluntDamageUntilUnconscious = 75
	var bluntDamageUntilTorsoIncapacitated = 75
	var brainDamageUntilDead = 100

	var sharpDamageUntilBleeding = 10
	var sharpDamageUntilSevered = 50

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
			p.Body.Head.Brain.IsConscious = false
			p.IsIncapacitated = true
			p.Body.Head.Brain.turnOff()
			p.RemoveLimb("Head")
			return
		}
		if p.Body.Head.BluntDamage > bluntDamageUntilBroken {
			p.Body.Head.IsBroken = true
			if p.Body.Head.BluntDamage > bluntDamageUntilUnconscious {
				p.Body.Head.Brain.IsConscious = false
				p.IsIncapacitated = true
			}
			if p.Body.Head.BluntDamage >= brainDamageUntilDead && p.Body.Head.Brain.active {
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
