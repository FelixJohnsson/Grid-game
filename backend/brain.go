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
        MotorCortexCurrentTask : MotorCortexAction{"None", "Idle", Location{0, 0}, false, false},

        PhysiologicalNeeds: PhysiologicalNeeds{0, 0, false, false, false, true, false, true, 100, false, false, true, false},
    }
}

func NewHumanBrain() *HumanBrain {
    return &HumanBrain{
        Brain: *NewBrain(),
        Owner: nil,
    }
}

func NewAnimalBrain() *AnimalBrain {
    return &AnimalBrain{
        Brain: *NewBrain(),
        Owner: nil,
    }
}

// MotorCortex is a function that handles the motor cortex of the brain
func (b *HumanBrain) MotorCortex() {
    if b.MotorCortexCurrentTask.ActionType == "Idle" {
        return
    } 
    fmt.Println(b.Owner.FullName + "'s motor cortex is executing the task.")

    switch b.MotorCortexCurrentTask.ActionType {
    case "Walk":
        success := b.WalkOverPath(b.MotorCortexCurrentTask.TargetLocation.X, b.MotorCortexCurrentTask.TargetLocation.Y)
        if !success {
            fmt.Println(b.Owner.FullName + " is unable to walk to the location.")
        } else {
            fmt.Println(b.Owner.FullName + " has arrived at the location.")
            b.MotorCortexCurrentTask = MotorCortexAction{"Drink water", "Walk", Location{b.Owner.Location.X, b.Owner.Location.Y}, false, true}
        }
    case "Run":
        //b.RunOverPath(closestWater.Location.X, closestWater.Location.Y)
    }
}

// SensoryCortex is a function that handles the sensory cortex of the brain
func (b *HumanBrain) SensoryCortex() {
    
}

// MainLoop is the main loop of the brain that handles the person's actions
func (b *HumanBrain) MainLoop() {
    for {
        select {
        case <-b.Ctx.Done():
            b.Active = false
            return
        default:
            if !b.Active {
                return
            }

            if !b.IsConscious {
                fmt.Println(b.Owner.FullName + "'s brain is not conscious but still alive.")
                return
            }

            if b.IsUnderAttack.Active {
                b.IsUnderAttackHandler()
            }

            b.OxygenHandler()

            b.PainHandler()
            b.FoodHandler()
            b.ThirstHandler()

            b.ClearWants()
            b.HomoSapiensCalculateWant()
            b.TranslateWantToTaskList()

            b.performActions()
            go b.MotorCortex()

            // Sleep for 2 seconds
            time.Sleep(2000 * time.Millisecond)
        }
    }
}

// OxygenHandler is a function that handles the oxygen level of the person
func (b *HumanBrain) OxygenHandler() {
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
func (b *HumanBrain) IncreaseHungerLevel() {
    b.PhysiologicalNeeds.Hunger += 1
}

//DecreaseHungerLevel is a function that decreases the hunger level of the person
func (b *HumanBrain) DecreaseHungerLevel(amount int) {
    b.PhysiologicalNeeds.Hunger -= amount
    if b.PhysiologicalNeeds.Hunger < 0 {
        b.PhysiologicalNeeds.Hunger = 0
    }
    fmt.Println("After eating, Current hunger level: ", b.PhysiologicalNeeds.Hunger)
}

//IncreaseThirstLevel is a function that increases the thirst level of the person
func (b *HumanBrain) IncreaseThirstLevel() {
    b.PhysiologicalNeeds.Thirst += 1
}

//DecreaseThirstLevel is a function that decreases the thirst level of the person
func (b *HumanBrain) DecreaseThirstLevel(amount int) {
    b.PhysiologicalNeeds.Thirst -= amount
    if b.PhysiologicalNeeds.Thirst < 0 {
        b.PhysiologicalNeeds.Thirst = 0
    }
    fmt.Println("After drinking water, Current thirst level: ", b.PhysiologicalNeeds.Thirst)
}

// FoodHandler is a function that handles the food level of the person
func (b *HumanBrain) FoodHandler() {
    b.IncreaseHungerLevel()
}

// ThirstHandler is a function that handles the thirst level of the person
func (b *HumanBrain) ThirstHandler() {
    b.IncreaseThirstLevel()
}

// IsUnderAttackHandler is a function that handles the person being under attack
func (b *HumanBrain) IsUnderAttackHandler() {
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


// Helper function that goes through the observation list and returns a boolean if the person is there
func (v Vision) HasPerson(fullName string) bool {
    for _, person := range v.Persons {
        if person.FullName == fullName {
            return true
        }
    }
    return false
}