package main

import (
	"context"
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

func (b *Brain) HeartHandler() {
    if b.Owner.Body.Torso.Heart.IsPumping {
        b.Owner.HeartBeat()
    } else {
        b.KillEntity("Heart stopped")
        return
    }
}

func (b *Brain) LungsHandler() {
	if b.CheckIfCanBreath() && b.Owner.Body.Blood.Oxygen < 100 {
        b.Owner.Breath()
    } else {
		b.AddHormone("Adrenaline", 70)
		b.AddHormone("Cortisol", 70)
    }
}

func (b *Brain) StomachHandler() {
    if b.Owner.Body.Torso.Stomach.IsDigesting && len(b.Owner.Body.Torso.Stomach.Contains) > 0 {
        b.Owner.Digest()
    }
}

func (b *Brain) KidneysHandler() {
    if b.Owner.Body.Torso.Kidneys.IsFiltering && b.Owner.Body.Blood.Toxins > 0 {
        b.Owner.Filter()
    } else {
        b.KillEntity("Kidneys stopped")
        return
    }
}

func (b *Brain) KillEntity(reason string){
    b.Owner.IsIncapacitated = true
    b.Owner.Brain.turnOff(reason)
    b.Owner.Brain.Active = false
    b.Owner.StopHeart()
    LogDeath(reason, b.Owner.FullName)
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


// ----------------- Hormones ------------------

func (b *Brain) AddHormone(hormone string, amount int) {
    switch hormone {
    case "Adrenaline":
        b.Owner.Body.Blood.Hormones.Adrenaline += amount
        if b.Owner.Body.Blood.Hormones.Adrenaline > 100 {
            b.Owner.Body.Blood.Hormones.Adrenaline = 100
        }
    case "Cortisol":
        b.Owner.Body.Blood.Hormones.Cortisol += amount
        if b.Owner.Body.Blood.Hormones.Cortisol > 100 {
            b.Owner.Body.Blood.Hormones.Cortisol = 100
        }
    case "Dopamine":
        b.Owner.Body.Blood.Hormones.Dopamine += amount
        if b.Owner.Body.Blood.Hormones.Dopamine > 100 {
            b.Owner.Body.Blood.Hormones.Dopamine = 100
        }
    case "Epinephrine":
        b.Owner.Body.Blood.Hormones.Epinephrine += amount
        if b.Owner.Body.Blood.Hormones.Epinephrine > 100 {
            b.Owner.Body.Blood.Hormones.Epinephrine = 100
        }
    case "Endorphin":
        b.Owner.Body.Blood.Hormones.Endorphin += amount
        if b.Owner.Body.Blood.Hormones.Endorphin > 100 {
            b.Owner.Body.Blood.Hormones.Endorphin = 100
        }
    case "Serotonin":
        b.Owner.Body.Blood.Hormones.Serotonin += amount
        if b.Owner.Body.Blood.Hormones.Serotonin > 100 {
            b.Owner.Body.Blood.Hormones.Serotonin = 100
        }
    }
}

func (b *Brain) RemoveHormone(hormone string, amount int) {
    switch hormone {
    case "Adrenaline":
        b.Owner.Body.Blood.Hormones.Adrenaline -= amount
        if b.Owner.Body.Blood.Hormones.Adrenaline < 0 {
            b.Owner.Body.Blood.Hormones.Adrenaline = 0
        }
    case "Cortisol":
        b.Owner.Body.Blood.Hormones.Cortisol -= amount
        if b.Owner.Body.Blood.Hormones.Cortisol < 0 {
            b.Owner.Body.Blood.Hormones.Cortisol = 0
        }
    case "Dopamine":
        b.Owner.Body.Blood.Hormones.Dopamine -= amount
        if b.Owner.Body.Blood.Hormones.Dopamine < 0 {
            b.Owner.Body.Blood.Hormones.Dopamine = 0
        }
    case "Epinephrine":
        b.Owner.Body.Blood.Hormones.Epinephrine += amount
        if b.Owner.Body.Blood.Hormones.Epinephrine < 0 {
            b.Owner.Body.Blood.Hormones.Epinephrine = 0
        }
    case "Endorphin":
        b.Owner.Body.Blood.Hormones.Endorphin -= amount
        if b.Owner.Body.Blood.Hormones.Endorphin < 0 {
            b.Owner.Body.Blood.Hormones.Endorphin = 0
        }
    case "Serotonin":
        b.Owner.Body.Blood.Hormones.Serotonin -= amount
        if b.Owner.Body.Blood.Hormones.Serotonin < 0 {
            b.Owner.Body.Blood.Hormones.Serotonin = 0
        }
    }
}