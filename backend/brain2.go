package main

import (
	"context"
)

type Brain2 struct {
    Active bool
    Ctx    context.Context
    Cancel context.CancelFunc
	ActionList []TargetedAction
	CurrentTask TargetedAction
    IsConscious bool
	CanBreath bool
	OxygenLevel int
	PainLevel int
	PainTolerance int
    IsAlive bool
    BrainDamage int
	IsUnderAttack IsUnderAttack
	Memories Memories
	MotorCortexCurrentTask MotorCortexAction

	PhysiologicalNeeds PhysiologicalNeeds

	Owner *Entity
}

type EntitySharedInfo struct {
	Name string
	Age int
	Gender string
	Species string
	Height int
	Weight int

	Genes []string

	IsMoving         TargetedAction
	IsTalking        TargetedAction
	IsSitting        TargetedAction
	IsEating         TargetedAction
	IsSleeping       TargetedAction

	Thinking         string
	WantsTo          []string

	Strength         int
	Agility          int
	Intelligence     int
	Stamina          int

	CombatExperience int
	CombatSkill      int
	Predator bool

	Relationships    []Relationship

	IsIncapacitated  bool
	VisionRange 	 int
	WorldProvider    WorldAccessor
	Location         Location

	Brain *Brain
}

type Entity struct {
	EntitySharedInfo

	Body  *Body
	Brain *Brain
}

// NewBrain creates a new Brain and assigns an owner to it.
func NewBrain2(Owner *Entity) *Brain {
    ctx, cancel := context.WithCancel(context.Background())
    return &Brain{
		Owner: Owner,
        Active:  false,
        Ctx:     ctx,
        Cancel:  cancel,
        ActionList: []TargetedAction{
            {"Idle", "", false, make([]BodyPartType, 0), 0},
        },
        CurrentTask: TargetedAction{"Idle", "", false, make([]BodyPartType, 0), 0},
        IsConscious: true,
        OxygenLevel: 100,
        PainLevel: 0,
        PainTolerance: 100,
        IsAlive:    true,
        BrainDamage: 0,
        IsUnderAttack: IsUnderAttack{false, nil, "", ""}, 
        Memories: Memories{make([]Memory, 0), make([]Memory, 0)},
        MotorCortexCurrentTask : MotorCortexAction{"None", "Idle", Location{0, 0}, false, false},

        PhysiologicalNeeds: PhysiologicalNeeds{0, 0, false, false, false, true, false, true, 100, false, false, true, false},
    }
}

