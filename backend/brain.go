package main

import (
	"context"
	"fmt"
)

// NewBrain creates a new Brain and assigns an owner to it.
func NewBrain(Owner *Entity) *Brain {
    ctx, cancel := context.WithCancel(context.Background())
    return &Brain{
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
        CanBreath: true,
        BrainDamage: 0,
        IsUnderAttack: IsUnderAttack{false, nil, "", ""}, 
        Memories: Memories{make([]Memory, 0), make([]Memory, 0)},
        MotorCortexCurrentTask : MotorCortexAction{"None", "Idle", Location{0, 0}, false, false},

        PhysiologicalNeeds: PhysiologicalNeeds{0, 0, false, false, false, true, false, true, 100, false, false, true, false},

        Owner: Owner,
    }
}

// OxygenHandler is a function that handles the oxygen level of the person
func (b *Brain) OxygenHandler() {
    if b.CheckIfCanBreath() {
        b.Breath()
    } else {
        fmt.Println(b.Owner.FullName + " is not able to breath.")
        b.CanBreath = false
    }
    b.ConsumeOxygen()

    if b.OxygenLevel <= 0 {
        b.turnOff()
        fmt.Println(b.Owner.FullName + "'s brain is shutting down due to lack of oxygen.")
        return
    }
}

// IncreaseHungerLevel is a function that increases the hunger level of the person
func (b *Brain) IncreaseHungerLevel() {
    b.PhysiologicalNeeds.Hunger += 1
}

//DecreaseHungerLevel is a function that decreases the hunger level of the person
func (b *Brain) DecreaseHungerLevel(amount int) {
    b.PhysiologicalNeeds.Hunger -= amount
    if b.PhysiologicalNeeds.Hunger < 0 {
        b.PhysiologicalNeeds.Hunger = 0
    }
    fmt.Println("After eating, Current hunger level: ", b.PhysiologicalNeeds.Hunger)
}

//IncreaseThirstLevel is a function that increases the thirst level of the person
func (b *Brain) IncreaseThirstLevel() {
    b.PhysiologicalNeeds.Thirst += 1
}

//DecreaseThirstLevel is a function that decreases the thirst level of the person
func (b *Brain) DecreaseThirstLevel(amount int) {
    b.PhysiologicalNeeds.Thirst -= amount
    if b.PhysiologicalNeeds.Thirst < 0 {
        b.PhysiologicalNeeds.Thirst = 0
    }
    fmt.Println("After drinking water, Current thirst level: ", b.PhysiologicalNeeds.Thirst)
}

// FoodHandler is a function that handles the food level of the person
func (b *Brain) FoodHandler() {
    b.IncreaseHungerLevel()
}

// ThirstHandler is a function that handles the thirst level of the person
func (b *Brain) ThirstHandler() {
    b.IncreaseThirstLevel()
}

// IsUnderAttackHandler is a function that handles the person being under attack
func (b *Brain) IsUnderAttackHandler() {
    if b.IsUnderAttack.Active && !b.IsUnderAttack.From.Brain.IsConscious {
        b.IsUnderAttack = IsUnderAttack{false, b.IsUnderAttack.From, "", ""}
        fmt.Println(b.Owner.FullName + " is no longer under attack because attacker is unconscious.")
        b.AddMemoryToShortTerm("Knocked out", b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)
    } else if b.IsUnderAttack.Active {
        b.UnderAttack(b.IsUnderAttack.From, b.IsUnderAttack.Target, b.IsUnderAttack.ByLimb)
        b.Owner.UpdateRelationship(b.IsUnderAttack.From.FullName, "Enemy", 100)
        b.AddMemoryToShortTerm("Under attack", b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)
        b.AddMemoryToLongTerm("Under attack", b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)
    }
}
