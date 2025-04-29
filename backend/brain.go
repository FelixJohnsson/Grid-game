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
        CognitiveMap: CognitiveMap{make(map[Location]CognitiveMapTile)},

        Owner: Owner,
    }
}

// OxygenHandler is a function that handles the oxygen level of the person
func (b *Brain) OxygenHandler() {
    if b.CheckIfCanBreath() {
        b.Breath()
    } else {
        b.CanBreath = false
    }
    b.ConsumeOxygen()

    if b.OxygenLevel <= 0 {
        b.turnOff("Oxygen level is 0")
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
}

// FoodHandler is a function that handles the food level of the person
func (b *Brain) FoodHandler() {
    previousHunger := b.PhysiologicalNeeds.Hunger
    b.IncreaseHungerLevel()
    currentHunger := b.PhysiologicalNeeds.Hunger

    // Log messages at important thresholds
    if GlobalInfoPanel != nil {
        // Only log when crossing thresholds to avoid spam
        if previousHunger < 30 && currentHunger >= 30 {
            LogFood(fmt.Sprintf("%s is getting hungry (%d%%)", b.Owner.FullName, currentHunger), b.Owner.FullName)
        } else if previousHunger < 50 && currentHunger >= 50 {
            LogFood(fmt.Sprintf("%s is very hungry (%d%%)", b.Owner.FullName, currentHunger), b.Owner.FullName)
        } else if previousHunger < 70 && currentHunger >= 70 {
            LogFood(fmt.Sprintf("%s is starving (%d%%)", b.Owner.FullName, currentHunger), b.Owner.FullName)
        } else if previousHunger < 90 && currentHunger >= 90 {
            LogFood(fmt.Sprintf("%s is near death from starvation (%d%%)", b.Owner.FullName, currentHunger), b.Owner.FullName)
        }
    }

    if b.PhysiologicalNeeds.Hunger >= 100 {
        if GlobalInfoPanel != nil {
            LogDeath(fmt.Sprintf("%s has died from starvation", b.Owner.FullName), b.Owner.FullName)
        }
        b.KillEntity("Starved")
    }
}

// ThirstHandler is a function that handles the thirst level of the person
func (b *Brain) ThirstHandler() {
    previousThirst := b.PhysiologicalNeeds.Thirst
    b.IncreaseThirstLevel()
    currentThirst := b.PhysiologicalNeeds.Thirst
    
    // Log messages at important thresholds
    if GlobalInfoPanel != nil {
        // Only log when crossing thresholds to avoid spam
        if previousThirst < 30 && currentThirst >= 30 {
            LogWater(fmt.Sprintf("%s is getting thirsty (%d%%)", b.Owner.FullName, currentThirst), b.Owner.FullName)
        } else if previousThirst < 50 && currentThirst >= 50 {
            LogWater(fmt.Sprintf("%s is very thirsty (%d%%)", b.Owner.FullName, currentThirst), b.Owner.FullName)
        } else if previousThirst < 70 && currentThirst >= 70 {
            LogWater(fmt.Sprintf("%s is dehydrated (%d%%)", b.Owner.FullName, currentThirst), b.Owner.FullName)
        } else if previousThirst < 90 && currentThirst >= 90 {
            LogWater(fmt.Sprintf("%s is near death from dehydration (%d%%)", b.Owner.FullName, currentThirst), b.Owner.FullName)
        }
    }
    
    if b.PhysiologicalNeeds.Thirst >= 100 {
        if GlobalInfoPanel != nil {
            LogDeath(fmt.Sprintf("%s has died from dehydration", b.Owner.FullName), b.Owner.FullName)
        }
        b.KillEntity("Dehydrated")
    }
}

func (b *Brain) KillEntity(reason string){
    b.Owner.IsIncapacitated = true
    b.Owner.Brain.turnOff(reason)
}


// IsUnderAttackHandler is a function that handles the person being under attack
func (b *Brain) IsUnderAttackHandler() {
    if b.IsUnderAttack.Active && !b.IsUnderAttack.From.Brain.IsConscious {
        b.IsUnderAttack = IsUnderAttack{false, b.IsUnderAttack.From, "", ""}
        b.AddMemoryToShortTerm("Knocked out", b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)
    } else if b.IsUnderAttack.Active {
        b.UnderAttack(b.IsUnderAttack.From, b.IsUnderAttack.Target, b.IsUnderAttack.ByLimb)
        b.Owner.UpdateRelationship(b.IsUnderAttack.From.FullName, "Enemy", 100)
        b.AddMemoryToShortTerm("Under attack", b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)
        b.AddMemoryToLongTerm("Under attack", b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)
    }
}
