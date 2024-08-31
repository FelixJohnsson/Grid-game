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
        IsConscious: true,
        OxygenLevel: 100,
        PainLevel: 0,
        PainTolerance: 100,
        IsAlive:    true,
        BrainDamage: 0,
        IsUnderAttack: IsUnderAttack{false, nil, "", ""}, 
        Memories: Memories{make([]Memory, 0), make([]Memory, 0)},

        PhysiologicalNeeds: PhysiologicalNeeds{0, 0, true, false, true, 100, false, false, false, false},
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
        b.Owner.Body.Head.Brain.turnOff()
        fmt.Println(b.Owner.FullName + "'s brain is shutting down due to lack of oxygen.")
        return
    }
}

// IsUnderAttackHandler is a function that handles the person being under attack
func (b *Brain) IsUnderAttackHandler() {
    if b.IsUnderAttack.Active && !b.IsUnderAttack.From.Body.Head.Brain.IsConscious {
        b.IsUnderAttack = IsUnderAttack{false, b.IsUnderAttack.From, "", ""}
        fmt.Println(b.Owner.FullName + " is no longer under attack because they're unconscious.")
        b.AddMemoryToShortTerm("Knocked out", b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)
    } else if b.IsUnderAttack.Active && b.IsUnderAttack.From.Body.Head.Brain.IsConscious {
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
        b.OxygenHandler()

        if !b.IsConscious{
            fmt.Println(b.Owner.FullName + "'s brain is not conscious but still alive.")
            return
        }

        if b.IsUnderAttack.Active {
            b.IsUnderAttackHandler()
        }
        
        if b.Active {
            b.CalculatePainLevel()
            b.PainHandler()

            obs := b.processInputs()
            b.makeDecisions(obs)
            b.CalculateWant()
            b.TranslateWantToTaskList()
            b.performActions()

            // Sleep for 2 seconds
            time.Sleep(2000 * time.Millisecond)
        }
        }
    }
}

func (b *Brain) processInputs() Vision {
    // Get the vision of the person
    obs := b.Owner.GetVision()

    // Check the tile type of the person
    b.checkTileType()

    // Check if the area is safe
    b.isAreaSafe(obs)

    return obs
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