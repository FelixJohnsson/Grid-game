package main

import (
	"context"
)

// ----------------- World ------------------

type Location struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type TileType int

const (
	Grass TileType = iota
	Water
	Mountain
)
type Tile struct {
	Type     TileType     `json:"Type"`
	Entity   *Entity      `json:"Ent,omitempty"`
	Items    []*Item      `json:"Items,omitempty"`
	Plant    *Plant       `json:"Plant,omitempty"`
	NutritionalValue int  `json:"NutritionalValue,omitempty"`
	Shelter  *Shelter     `json:"Shelter,omitempty"`
	Location Location     `json:"Location"`
	BipedalAnimal   *BipedalAnimal      `json:"BipedalAnimal,omitempty"`
	QuadrupedalAnimal   *QuadrupedalAnimal      `json:"QuadrupedalAnimal,omitempty"`
}
type World struct {
	Tiles [][]Tile `json:"tiles"`
	Width int      `json:"width"`
	Height int     `json:"height"`
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
	WithEntity string
	Relationship string
	Intensity int
}

type CognitiveMapEntity struct {
	FullName string
	SpeciesType SpeciesType
	IsAlive bool
}

type CognitiveMapPlant struct {
	Name          PlantType
	IsAlive       bool
	ProducesFruit bool
	PlantStage    PlantStage
}

type CognitiveMapTile struct {
	TileType TileType 
	Entity CognitiveMapEntity
	Plant CognitiveMapPlant
}

type CognitiveMap struct {
    KnownTiles map[Location]CognitiveMapTile
}

type Brain struct {
        Active  bool
        Ctx     context.Context
        Cancel  context.CancelFunc
        ActionList []TargetedAction
        CurrentTask  TargetedAction
        IsConscious bool
        OxygenLevel int
        PainLevel  int
        PainTolerance int
        IsAlive bool
		CanBreath bool
        BrainDamage int
        IsUnderAttack IsUnderAttack
        Memories   Memories
        MotorCortexCurrentTask   MotorCortexAction

        PhysiologicalNeeds  PhysiologicalNeeds
		CognitiveMap     CognitiveMap

		Owner *Entity
}

type Entity struct {
	Age              int
	Title 		     string
	FirstName        string
	FamilyName       string
	FullName 	     string
	Initials         string
	IsChild          bool
	Gender           string
	Occupation       Jobs
	SkinColor        string
	Entityality 	 string
	Genes            []string
	Species          SpeciesType

	OwnedItems	     []*Item
	OwnedBipedalAnimals	 []*BipedalAnimal
	OwnedQuadrupedalAnimals	 []*QuadrupedalAnimal

	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction
	IsBleeding       bool

	Thinking         string
	WantsTo          []string
	FeelingSafe 	 int
	FeelingScared	 int

	Body 		     *EntityBody
	Brain 		     *Brain

	Strength         int
	Agility          int
	Intelligence     int
	Charisma         int
	Stamina          int
	Curiosity        int

	CombatExperience int
	CombatSkill      int
	CombatStyle      string

	Relationships    []Relationship

	IsIncapacitated  bool
	VisionRange 	 int
	WorldProvider    WorldAccessor
	Location         Location
}

// ----------------- Brain ------------------

type TargetedAction struct {
	Action TaskType
	Target string
	IsActive bool
	RequiresLimb []BodyPartType
	Priority int
}

type MotorCortexAction struct {
	ActionReason string
	ActionType string
	TargetLocation Location
	IsActive bool
	Finished bool
}
type IsUnderAttack struct {
	Active bool
	From *Entity
	Target BodyPartType
	ByLimb BodyPartType
}

type Memory struct {
	Event string
	Details string
	Location Location
}
type Memories struct {
	ShortTermMemory []Memory
	LongTermMemory  []Memory
}
type RequestedAction struct {
	TargetedAction
	From *Entity
}

type PhysiologicalNeeds struct {
	Thirst int
	Hunger int
	FoodSupply bool
	WayOfGettingFood bool
	WayOfGettingWater bool
	CanBreath bool
	HasShelter bool
	IsSufficientlyWarm bool
	Rested int
	IsInPain bool
	NeedToExcrete bool
	IsInSafeArea bool
	IsCapableOfDefendingSelf bool
}

type Vision struct {
	Plants    []*Plant            `json:"Plants"`
	Entities  []EntityInVision    `json:"Entities"`
	Tiles     []Tile              `json:"Tile"`
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

type BodyPartType string

type BodyPart struct {
	Name string
	IsBleeding bool
	IsBroken bool
	IsObstructed bool
}

const (
	RightHand BodyPartType = "RightHand"
	LeftHand  BodyPartType = "LeftHand"
	RightFoot BodyPartType = "RightFoot"
	LeftFoot  BodyPartType = "LeftFoot"
	RightLeg  BodyPartType = "RightLeg"
	LeftLeg   BodyPartType = "LeftLeg"
	TheHead   BodyPartType = "Head"
	Torso     BodyPartType = "Torso"

	Mouth    BodyPartType = "Mouth"
	Nose     BodyPartType = "Nose"
	Eyes     BodyPartType = "Eyes"
	Ears     BodyPartType = "Ears"
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
	Eyes *BodyPart
	Ears *BodyPart
	Nose *BodyPart
	Mouth *BodyPart
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

type EntityBody struct {
	Head          *Head
	Torso         *LimbStatus

	RightFrontLeg *Leg
	LeftFrontLeg  *Leg
	RightBackLeg  *Leg
	LeftBackLeg   *Leg

	RightArm      *Arm
	LeftArm       *Arm

	RightLeg      *Leg
	LeftLeg       *Leg

	Wings         *LimbThatCanMove
	Tail          *LimbThatCanMove

}

// ----------------- Animals -----------------
type QuadrupedalAnimal struct {
	Age              int
	FullName 	     string
	IsChild          bool
	Gender           string
	Description      string
	Color            string
	Entityality 	 string
	Genes            []string
	Species		     string
	OwnedItems       []*Item

	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction

	Thinking         string
	WantsTo          []string
	FeelingSafe 	 int
	FeelingScared	 int

	Body 		     *EntityBody
	Brain 		     *Brain

	Strength         int
	Agility          int
	Intelligence     int
	Charisma         int
	Stamina          int

	CombatExperience int
	CombatSkill      int

	Relationships    []Relationship

	IsIncapacitated  bool
	VisionRange 	 int
	WorldProvider    WorldAccessor
	Location         Location
}

type BipedalAnimal struct {
	Age              int
	FullName 	     string
	IsChild          bool
	Gender           string
	Description      string
	Color            string
	Entityality 	 string
	Genes            []string
	Species		     string
	OwnedItems       []*Item

	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction

	Thinking         string
	WantsTo          []string
	FeelingSafe 	 int
	FeelingScared	 int

	Body 		     *EntityBody
	Brain 		     *Brain

	Strength         int
	Agility          int
	Intelligence     int
	Charisma         int
	Stamina          int

	CombatExperience int
	CombatSkill      int

	Relationships    []Relationship

	IsIncapacitated  bool
	VisionRange 	 int
	WorldProvider    WorldAccessor
	Location         Location
}

// ----------------- Food -------------------

type Food interface {
	GetName() string
	GetNutritionalValue() int
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
	Name        string
	RipeAge     int
	IsRipe      bool
	NutritionalValue int
}

func (f Fruit) GetName() string {
	return f.Name
}

func (f Fruit) GetNutritionalValue() int {
	return f.NutritionalValue
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
	Name          PlantType
	Age           int
	Health        int
	IsAlive       bool
	ProducesFruit bool
	Fruit         []Fruit
	PlantStage    PlantStage
	PlantLife     *PlantLife
	Location	  Location
}

// ----------------- Liquid ------------------

type Liquid struct {
	Name string
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
	Location Location
}

// ----------------- Cleaned ------------------

type CleanedTile struct {
    Type     TileType         `json:"Type"`
    Entity   *EntityCleaned   `json:"Entity,omitempty"`
	Items    []*Item          `json:"Items,omitempty"`
	Plant    *PlantCleaned    `json:"Plant,omitempty"`
	Animal   *AnimalCleaned   `json:"Animal,omitempty"`
	Shelter  *Shelter         `json:"Shelter,omitempty"`
}

type AnimalCleaned struct {
	FirstName    string       `json:"FirstName"`
	FamilyName   string       `json:"FamilyName"`
	FullName     string       `json:"FullName"`
	Gender 	     string	      `json:"Gender"`
	Age 		 int          `json:"Age"`
	Title 		 string       `json:"Title"`

	Location     Location     `json:"Location"`

	Thinking 	 string       `json:"Thinking"`

	Head 		 *Head        `json:"Head"`
	Torso 		 *LimbStatus  `json:"Torso"`
	RightFrontLeg *Leg        `json:"RightFrontLeg"`
	LeftFrontLeg  *Leg         `json:"LeftFrontLeg"`
	RightBackLeg  *Leg         `json:"RightBackLeg"`
	LeftBackLeg   *Leg         `json:"LeftBackLeg"`

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

	CurrentTask TargetedAction `json:"CurrentTask"`
}

type PlantCleaned struct {
	Name      PlantType `json:"Name"`
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

type EntityCleaned struct {
	FirstName    string       `json:"FirstName"`
	FamilyName   string       `json:"FamilyName"`
	FullName     string       `json:"FullName"`
	Gender 	     string	      `json:"Gender"`
	Age 		 int          `json:"Age"`
	Title 		 string       `json:"Title"`

	Location     Location     `json:"Location"`

	Thinking 	 string       `json:"Thinking"`

	Head 		 *Head `json:"Head"`
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

	CurrentTask TargetedAction `json:"CurrentTask"`
}

type EntityInVision struct {
	FullName     string
	FirstName    string 
	FamilyName   string
	Gender 		 string
	Age 		 int
	Title 		 string
	Location     Location
	Body 		 *EntityBody
}
