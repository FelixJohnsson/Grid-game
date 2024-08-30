package main

import (
	"context"
)

// ----------------- World ------------------

type WorldState struct {
	Map 	 	[][]string;
  	Buildings	[]Building;
  	Persons	 	[]Person;
  	Resources 	Resources;
}

type TileType int

const (
	Grass TileType = iota
	Water
	Mountain
)
type Tile struct {
	Type     TileType  `json:"Type"`
	Building *Building `json:"Building,omitempty"`
	Persons  []*Person   `json:"Persons,omitempty"`
	Items    []*Item    `json:"Items,omitempty"`
	Plants   []*Plant    `json:"Plant,omitempty"`
	NutritionalValue int `json:"NutritionalValue,omitempty"`
}
type World struct {
	Tiles [][]Tile `json:"tiles"`
}

// ----------------- People -----------------

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
	Occupation       Jobs
	SkinColor        string
	Personality 	 string
	Genes            []string

	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction

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

// ----------------- Brain ------------------

type TargetedAction struct {
	Action string
	Target string
	IsActive bool
	RequiresLimb []string
}

type IsUnderAttack struct {
	Active bool
	From *Person
	Target LimbType
	ByLimb LimbType
}

type Memory struct {
	Event string
	Location Location
}
type Memories struct {
	ShortTermMemory []Memory
	LongTermMemory  []Memory
}
type RequestedAction struct {
	TargetedAction
	From *Person
}

type Brain struct {
	Owner  *Person
    Active bool
    Ctx    context.Context
    Cancel context.CancelFunc
	Actions []TargetedAction
    IsConscious bool
    IsAlive bool
    BrainDamage int
	IsUnderAttack IsUnderAttack
	Memories Memories
}
type Vision struct {
	Buildings []BuildingCleaned `json:"buildings"`
	Persons   []PersonInVision   `json:"persons"`
}

// ----------------- Body -------------------

type LimbStatus struct {
	BluntDamage int
	SharpDamage int
	IsBleeding bool
	IsBroken bool
	Residues []Residue
	CoveredWith []Wearable
	IsAttached bool
}

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

// ----------------- Plants -----------------
type PlantAction struct {
	Name string
	Target *Tile
	Priority int
}

type PlantLife struct {
	active bool
	ctx    context.Context
	cancel context.CancelFunc
	actions []PlantAction
}

type Nutrients struct {
	Calories int
	Carbs    int
	Protein  int
	Fat      int

	Vitamins int
	Minerals int
}

type Fruit struct {
	Name      string
	Taste     string
	Age       int
	RipeAge   int
	IsRipe    bool
	Nutrients []Nutrients
}

type PlantStage int

const (
	Seed PlantStage = iota
	Sprout
	Vegetative
	Flowering
	Fruiting
)
type Plant struct {
	Name          string
	Age           int
	Health        int
	IsAlive       bool
	ProducesFruit bool
	Fruit         []Fruit
	PlantStage    PlantStage
	PlantLife     *PlantLife
}

// ----------------- Items ------------------

type Wearable struct {
	Name string
	Material string
	Protection int
}

type Material struct {
	Name         string
	Type         string
	Hardness     int
	Weight       int
	Density      int
	Malleability int
}

type Residue struct {
	Name   string
	Amount int
}

type Item struct {
	Name      string
	Sharpness int
	Bluntness int
	Weight    int
	Material  []Material
	Residues  []Residue
}

// ----------------- Cleaned ------------------

type CleanedTile struct {
    Type     TileType         `json:"Type"`
    Building *BuildingCleaned `json:"Building,omitempty"`
    Persons  []PersonCleaned  `json:"Persons,omitempty"`
	Items    []*Item          `json:"Items,omitempty"`
	Plants   []*PlantCleaned  `json:"Plants,omitempty"`
}
type PlantCleaned struct {
	Name      string `json:"Name"`
	Age       int    `json:"Age"`
	Health    int    `json:"Health"`
	IsAlive   bool   `json:"IsAlive"`
	ProducesFruit bool   `json:"ProducesFruit"`
	Fruit    []Fruit `json:"Fruit"`
	PlantStage PlantStage    `json:"PlantStage"`
}
type BuildingCleaned struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Location Location `json:"location"`
}

type HeadCleaned struct {
	LimbStatus
}

type PersonCleaned struct {
	FirstName    string       `json:"FirstName"`
	FamilyName   string       `json:"FamilyName"`
	FullName     string       `json:"FullName"`
	Gender 	     string	      `json:"Gender"`
	Age 		 int          `json:"Age"`
	Title 		 string       `json:"Title"`

	Location     Location     `json:"Location"`

	Thinking 	 string       `json:"Thinking"`

	Head 		 HeadCleaned `json:"Head"`
	Torso 		 *LimbStatus  `json:"Torso"`
	RightArm 	 *Arm         `json:"RightArm"`
	LeftArm 	 *Arm         `json:"LeftArm"`
	RightLeg 	 *Leg         `json:"RightLeg"`
	LeftLeg 	 *Leg         `json:"LeftLeg"`

	Strength 	 int          `json:"Strength"`
	Agility 	int           `json:"Agility"`
	Intelligence int          `json:"Intelligence"`
	Charisma 	int           `json:"Charisma"`
	Stamina 	int           `json:"Stamina"`

	CombatExperience int      `json:"CombatExperience"`
	CombatSkill int           `json:"CombatSkill"`
	CombatStyle string        `json:"CombatStyle"`

	IsIncapacitated bool       `json:"IsIncapacitated"`

	Relationships []Relationship `json:"Relationships"`
}

type PersonInVision struct {
	FullName     string
	FirstName    string 
	FamilyName   string
	Gender 		 string
	Age 		 int
	Title 		 string
	Location     Location
	Body 		 *HumanBody
}
