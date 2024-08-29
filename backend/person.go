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

// Available actions
var actions = []string {
	// World actions
	"Move",
	"Talk",

	"Grab",
	"Drop",

	"Sit",
	"Hold",
	"Eat",
	"Sleep",
	"Work",
	"Throw",
	"Build",
	"Dig",
	"Plant",
	"Harvest",
	"Chop",
	"Mine",

	"Open",
	"Close",
	"Enter",
	"Exit",
	"Use",

	// Hostile actions
	"Attack",
	"Steal",
	"Destroy",

	// Friendly actions
	"Help",
	"Gift",
	"Protect",
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

	RightHand        LimbThatCanHold
	LeftHand 	     LimbThatCanHold
	Back 		     LimbStatus
	LeftFoot 	     LimbStatus
	RightFoot 	     LimbStatus
	Head 		     LimbStatus
	Torso 		     LimbStatus
	Legs 		     LimbStatus

	Strength         int
	Agility          int
	Intelligence     int
	Charisma         int
	Stamina          int

	CombatExperience int
	CombatSkill      int
	CombatStyle      string

	Relationships    []Relationship

	Brain			 Brain
	IsConscious	     bool
	IsIncapacitated  bool 
	BrainDamage 	 int
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
	brain := NewBrain()

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

		Brain:            *brain,
		IsConscious:      false,
		BrainDamage:      0,
		VisionRange:      5,
		Location:         Location{X: x, Y: y},
		WorldProvider:    worldAccessor,
		OnTileType:       0,
		RightHand:        LimbThatCanHold{},
		LeftHand:         LimbThatCanHold{},
		Back:             LimbStatus{},
		LeftFoot:         LimbStatus{},
		RightFoot:        LimbStatus{},
		Head:             LimbStatus{},
		Torso:            LimbStatus{},
		Legs:             LimbStatus{},

		Strength:         1,
		Agility:          1,
		Intelligence:     1,
		Charisma:         1,
		Stamina:          1,

		CombatExperience: 1,
		CombatSkill:      1,
		CombatStyle:      "One handed",
	}

	person.Brain.owner = person
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
	if p.RightHand.Items == nil {
		p.RightHand.Items = []*Item{item}
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
	if p.RightHand.Items != nil {
		p.RightHand.Items = nil
	} else {
		fmt.Println("Right hand is empty")
	}
}

// Grab in the left hand
func (p *Person) GrabLeft(item *Item) {
	if p.LeftHand.Items == nil {
		p.LeftHand.Items = []*Item{item}
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
	if p.LeftHand.Items != nil {
		p.LeftHand.Items = nil
	} else {
		fmt.Println("Left hand is empty")
	}
}

// AddResidue adds a residue to the limb
func (p *Person) AddResidue(limb string, residue Residue) {
	switch limb {
	case "Back":
		p.Back.Residues = append(p.Back.Residues, residue)
	case "LeftFoot":
		p.LeftFoot.Residues = append(p.LeftFoot.Residues, residue)
	case "RightFoot":
		p.RightFoot.Residues = append(p.RightFoot.Residues, residue)
	case "Head":
		p.Head.Residues = append(p.Head.Residues, residue)
	case "Torso":
		p.Torso.Residues = append(p.Torso.Residues, residue)
	case "Legs":
		p.Legs.Residues = append(p.Legs.Residues, residue)
	case "RightHand":
		p.RightHand.Residues = append(p.RightHand.Residues, residue)
	case "LeftHand":
		p.LeftHand.Residues = append(p.LeftHand.Residues, residue)
	}
}

func (p *Person) GetVision() Vision {
    return p.WorldProvider.GetVision(p.Location.X, p.Location.Y, p.VisionRange)
}

func (p *Person) GetPersonByFullName(FullName string) *Person {
	return p.WorldProvider.GetPersonByFullName(FullName)
}

func (p *Person) addEmployer(building *Building) {
	if p.IsChild {
		return
	}
	p.IsWorkingAt = building
	building.Workers = append(building.Workers, *p)
	switch building.Type {
	case "Lumberjack":
		p.Occupation = Lumberjack
	case "Mine":
		p.Occupation = Miner
	case "Farm":
		p.Occupation = Farmer
	default:
		p.Occupation = Unemployed
	}
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

		if p.RightHand.IsBroken {
			return damage
		} else {
			damage.AmountBluntDamage = p.RightHand.Items[0].Bluntness
			damage.AmountSharpDamage = p.RightHand.Items[0].Sharpness

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

	switch limb {
		case "Head":
		p.Head.BluntDamage += damage.AmountBluntDamage
		p.Head.SharpDamage += damage.AmountSharpDamage

		if damage.AmountBluntDamage > 0 {
			p.Brain.owner.BrainDamage += damage.AmountBluntDamage
		}

		if p.Head.BluntDamage > 50 {
			p.Head.IsBroken = true
			if p.Head.BluntDamage > 75 {
				p.IsConscious = false
				p.IsIncapacitated = true
			}
			if p.Head.BluntDamage > 100 && p.Brain.active {
				p.Brain.turnOff()
			}
		}
		if p.Head.SharpDamage > 20 {
			p.Head.IsBleeding = true
		}
	case "Back":
		p.Back.BluntDamage += damage.AmountBluntDamage
		p.Back.SharpDamage += damage.AmountSharpDamage
		// Check if the limb is broken
		if p.Back.BluntDamage > 50 {
			p.Back.IsBroken = true
			p.IsIncapacitated = true
		} else if p.Back.SharpDamage > 20 {
			p.Back.IsBleeding = true
		}
	case "LeftFoot":
		p.LeftFoot.BluntDamage += damage.AmountBluntDamage
		p.LeftFoot.SharpDamage += damage.AmountSharpDamage

		if p.LeftFoot.BluntDamage > 50 {
			p.LeftFoot.IsBroken = true
			p.IsIncapacitated = true
		} else if p.LeftFoot.SharpDamage > 20 {
			p.LeftFoot.IsBleeding = true
		}
	case "RightFoot":
		p.RightFoot.BluntDamage += damage.AmountBluntDamage
		p.RightFoot.SharpDamage += damage.AmountSharpDamage

		if p.RightFoot.BluntDamage > 50 {
			p.RightFoot.IsBroken = true
			p.IsIncapacitated = true
		}
		if p.RightFoot.SharpDamage > 20 {
			p.RightFoot.IsBleeding = true
		}
	case "Torso":
		p.Torso.BluntDamage += damage.AmountBluntDamage
		p.Torso.SharpDamage += damage.AmountSharpDamage

		if p.Torso.BluntDamage > 50 {
			p.Torso.IsBroken = true
			if p.Torso.BluntDamage > 75 {
				p.IsIncapacitated = true
			}
		}
		if p.Torso.SharpDamage > 20 {
			p.Torso.IsBleeding = true
		}
	case "Legs":
		p.Legs.BluntDamage += damage.AmountBluntDamage
		p.Legs.SharpDamage += damage.AmountSharpDamage

		if p.Legs.BluntDamage > 50 {
			p.Legs.IsBroken = true
			p.IsIncapacitated = true
		}
		if p.Legs.SharpDamage > 20 {
			p.Legs.IsBleeding = true
		}
	case "RightHand":
		p.RightHand.BluntDamage += damage.AmountBluntDamage
		p.RightHand.SharpDamage += damage.AmountSharpDamage

		if p.RightHand.BluntDamage > 50 {
			p.RightHand.IsBroken = true
			p.IsIncapacitated = true
			p.RightHand.Items = nil
		}
		if p.RightHand.SharpDamage > 20 {
			p.RightHand.IsBleeding = true
		}
	case "LeftHand":
		p.LeftHand.BluntDamage += damage.AmountBluntDamage
		p.LeftHand.SharpDamage += damage.AmountSharpDamage

		if p.LeftHand.BluntDamage > 50 {
			p.LeftHand.IsBroken = true
			p.IsIncapacitated = true
			p.LeftHand.Items = nil
		}
		if p.LeftHand.SharpDamage > 20 {
			p.LeftHand.IsBleeding = true
		}
	}
		
}



// ---------------- Create a new person ----------------

func (w *World) createNewPerson(x, y int) *Person {
    person := NewPerson(w, x, y)
    return person
}
