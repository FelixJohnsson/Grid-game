package main

import (
	"context"
	"fmt"
	"time"
)

// NewBrain creates a new Brain and assigns an owner to it.
func NewBrain() *Brain {
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
        BrainDamage: 0,
        IsUnderAttack: IsUnderAttack{false, nil, "", ""}, 
        Memories: Memories{make([]Memory, 0), make([]Memory, 0)},

        PhysiologicalNeeds: PhysiologicalNeeds{0, 0, false, false, false, true, false, true, 100, false, false, false, false},
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

//IncreaseThirstLevel is a function that increases the thirst level of the person
func (b *Brain) IncreaseThirstLevel() {
    b.PhysiologicalNeeds.Thirst += 1
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
    if b.IsUnderAttack.Active && !b.IsUnderAttack.From.Body.Head.Brain.IsConscious {
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

// MainLoop is the main loop of the brain that handles the person's actions
func (b *Brain) mainLoop() {
    for {
        select {
        case <-b.Ctx.Done():
            b.Active = false
            return
        default:

        if !b.Active {
            return
        }
        b.OxygenHandler()

        if !b.IsConscious && b.Active {
            fmt.Println(b.Owner.FullName + "'s brain is not conscious but still alive.")
            return
        }

        if b.IsUnderAttack.Active {
            b.IsUnderAttackHandler()
        }
        
        if b.Active && b.IsConscious {
            b.CalculatePainLevel()
            b.PainHandler()
            b.FoodHandler()
            b.ThirstHandler()


            b.CalculateWant()
            b.TranslateWantToTaskList()
            fmt.Println(b.ActionList)
            if !b.CurrentTask.IsActive {
                b.performActions()
            }

            // Sleep for 2 seconds
            time.Sleep(2000 * time.Millisecond)
        }
        }
    }
}

// Helper function that goes through the observation list and returns a boolean if the person is there
func (v Vision) HasPerson(fullName string) bool {
    for _, person := range v.Persons {
        if person.FullName == fullName {
            return true
        }
    }
    return false
}