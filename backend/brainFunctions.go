package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ----------------- Brain Functions -----------------

// ----------------- Turn On and Off -----------------
func (b *Brain) turnOn() {
    if b.Active {
        fmt.Println("Brain is already active.")
        return
    }

    fmt.Println("Brain for: " + b.Owner.FullName + " is now active.")
    b.IsConscious = true
    b.Active = true

    go b.mainLoop()
}

func (b *Brain) turnOff() {
    if !b.Active {
        fmt.Println("Brain is already inactive.")
        return
    }

    fmt.Println(b.Owner.FullName + "'s brain is shutting down")
    b.Cancel()
}

// LooseConsciousness makes the person loose consciousness - duration is in seconds
func (b *Brain) LooseConsciousness(duration int) {
    b.IsConscious = false
    go func() {
        time.Sleep(time.Duration(duration) * time.Second)
        fmt.Println(b.Owner.FullName + " regained consciousness.")
        b.IsConscious = true
    }()
}

// Decide if the area is safe or not
func (b *Brain) isAreaSafe(obs Vision) {
    // Loop through the observations and make decisions based on relationships with the people

    collectiveIntensity := 0
    numberOfPeople := len(obs.Persons)

    for _, person := range obs.Persons {
        if b.Owner.hasRelationship(person.FullName) {
            for _, relationship := range b.Owner.Relationships {
                if relationship.WithPerson == person.FullName {
                    collectiveIntensity += relationship.Intensity
                }
            }
        }
    }

    // This is a pretty dumb way to determine if the area is safe or not, but it's a start
    b.Owner.FeelingSafe = collectiveIntensity/numberOfPeople
}

func (b *Brain) makeDecisions(obs Vision) {
    // Check if we're engaging in conversation with someone and if we are and we dont have that person in the observation, we should cancel the conversation
    if b.Owner.IsTalking.IsActive {
        if !obs.HasPerson(b.Owner.IsTalking.Target) {
            fmt.Println(b.Owner.FullName + " is no longer talking to " + b.Owner.IsTalking.Target)
            b.Owner.IsTalking = TargetedAction{"", "", false, make([]BodyPartType, 0), 10}
        }
    }
    // Loop through the observations and make decisions based on people
    for _, person := range obs.Persons {
        if person.FullName != b.Owner.FullName {
            if b.Owner.hasRelationship(person.FullName) {
                for _, relationship := range b.Owner.Relationships {
                    if relationship.WithPerson == person.FullName {
                        relationship.Intensity++
                        if relationship.Intensity > 3 { // This should be a constant and above 15 and below 40
                            relationship.Relationship = "Aquantance"
                            if !b.Owner.IsTalking.IsActive {
                                targetPerson := b.Owner.GetPersonByFullName(relationship.WithPerson)
                                if targetPerson != nil {
                                    b.SendTaskRequest(targetPerson, "Talk")
                                }
                            }
                        }
                        b.Owner.UpdateRelationship(person.FullName, relationship.Relationship, relationship.Intensity)
                    }
                }
            } else { // If the person does not have a relationship with the other person
                b.Owner.addRelationship(person, "Stranger", 0)
            }
        }
    }
}

// ----------------- Memory Functions -----------------

// AddMemoryToShortTerm adds a memory to the short term memory
func (b *Brain) AddMemoryToShortTerm(event string, details string, location Location) {
    memory := Memory{event, details, location}
    b.Memories.ShortTermMemory = append(b.Memories.ShortTermMemory, memory)
}

// AddMemoryToLongTerm adds a memory to the long term memory
func (b *Brain) AddMemoryToLongTerm(event string, details string, location Location) {
    memory := Memory{event, details, location}
    b.Memories.LongTermMemory = append(b.Memories.LongTermMemory, memory)
}


// ----------------- Pain -----------------------------

// PainHandler is a function that handles the pain level of the person
func (b *Brain) PainHandler() {
    if b.PainLevel > b.PainTolerance {
        // Get a random number for duration of unconsciousness
        durationInSeconds := rand.Intn(60)
        
        b.LooseConsciousness(durationInSeconds)
        fmt.Println(b.Owner.FullName + "'s brain lost consciousness due to pain." + " Duration: " + fmt.Sprint(durationInSeconds) + " seconds.")
    }
}

// ApplyPain - Apply pain to the person
func (b *Brain) ApplyPain(amount int) {
	b.PainLevel += amount
	b.PainHandler()
}

// CalculatePainLevel - Calculate the pain level of the person
func (b *Brain) CalculatePainLevel() {
    // We need to loop over the body parts and check if it's broken or bleeding

    if b.Owner.Body.Head != nil {
        if b.Owner.Body.Head.IsBroken {
            b.ApplyPain(5)
        } 
        if b.Owner.Body.Head.IsBleeding {
            b.ApplyPain(2)
        }
        if b.Owner.Body.Head.Ears != nil && b.Owner.Body.Head.Ears.IsBleeding || b.Owner.Body.Head.Ears.IsBroken {
           b.ApplyPain(1)
        } 
        if b.Owner.Body.Head.Eyes != nil && b.Owner.Body.Head.Eyes.IsBleeding {
           b.ApplyPain(5)
        }
        if b.Owner.Body.Head.Nose != nil && b.Owner.Body.Head.Nose.IsBleeding || b.Owner.Body.Head.Nose.IsBroken {
            b.ApplyPain(2)
        }
        if b.Owner.Body.Head.Mouth != nil  &&b.Owner.Body.Head.Mouth.IsBleeding || b.Owner.Body.Head.Mouth.IsBroken {
            b.ApplyPain(2)
        }
    }
    if b.Owner.Body.Torso.IsBleeding || b.Owner.Body.Torso.IsBroken {
        b.ApplyPain(5)
    }
    if b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.IsBleeding || b.Owner.Body.RightArm.IsBroken {
        b.ApplyPain(5)
    }
    if b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.IsBleeding || b.Owner.Body.LeftArm.IsBroken {
        b.ApplyPain(5)
    }
    if b.Owner.Body.RightLeg != nil && b.Owner.Body.RightLeg.IsBleeding || b.Owner.Body.RightLeg.IsBroken {
        b.ApplyPain(5)
    }
    if b.Owner.Body.LeftLeg != nil && b.Owner.Body.LeftLeg.IsBleeding || b.Owner.Body.LeftLeg.IsBroken {
        b.ApplyPain(5)
    }
    if b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.Hand != nil && b.Owner.Body.RightArm.Hand.IsBleeding || b.Owner.Body.RightArm.Hand.IsBroken {
        b.ApplyPain(2)
    }
    if b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.Hand != nil && b.Owner.Body.LeftArm.Hand.IsBleeding || b.Owner.Body.LeftArm.Hand.IsBroken {
        b.ApplyPain(2)
    }
    if b.Owner.Body.RightLeg != nil && b.Owner.Body.RightLeg.Foot != nil && b.Owner.Body.RightLeg.Foot.IsBleeding || b.Owner.Body.RightLeg.Foot.IsBroken {
        b.ApplyPain(2)
    }
    if b.Owner.Body.LeftLeg != nil && b.Owner.Body.LeftLeg.Foot != nil && b.Owner.Body.LeftLeg.Foot.IsBleeding || b.Owner.Body.LeftLeg.Foot.IsBroken {
        b.ApplyPain(2)
    }
}

// ----------------- Oxygen levels -----------------

// Breath 
func (b *Brain) Breath() {
    b.OxygenLevel += 10
}

// ConsumeOxygen
func (b *Brain) ConsumeOxygen() {
    b.OxygenLevel -= 10
}

// CheckIfCanBreah - Check if the person can breath
func (b *Brain) CheckIfCanBreath() bool {
    mouthCanBreath := b.Owner.Body.Head.Mouth != nil && !b.Owner.Body.Head.Mouth.IsObstructed
    noseCanBreath := b.Owner.Body.Head.Nose != nil && !b.Owner.Body.Head.Nose.IsObstructed && !b.Owner.Body.Head.Nose.IsBroken

    return mouthCanBreath || noseCanBreath
}


// ----------------- Wants ---------------------

// Pseudo code for a "Want" system
// For example, if the person is homeless, they want shelter
// So we have to construct a system that translates a want into a list of tasks with priorities
// The brain will then decide which task to do based on the priority and the current situation

// Pseudo code for a "Want to task" system

// Let's describe the Maslow's hierarchy of needs in terms of wants
// 1. Physiological needs
// 2. Safety needs
// 3. Love and belonging
// 4. Self-esteem
// 5. Self-actualization

// 1. Physiological needs
/* type PhysiologicalNeeds struct {
	IsBreathing bool
	IsInPain bool

	Thirst int
	Hunger int
	IsSufficientlyWarm bool

	NeedToExcrete bool

	IsInSafeArea bool
	IsCapableOfDefendingSelf bool

	HasShelter bool
	Rested int
} */

// CalculateWant - Calculate the want of the person
func (b *Brain) CalculateWant() {
    // Check the current situation of the person
    // We can be more specific when we have a "Dopamin" system in place

    // Check if breathing
    if !b.CheckIfCanBreath() {
        b.Owner.WantsTo = "Be able to breath"
    }

    // Check if in pain
    if b.PhysiologicalNeeds.IsInPain {
        b.Owner.WantsTo = "Relieve pain"
    }

    // Check thirst
    if b.PhysiologicalNeeds.Thirst > 30 {
        b.Owner.WantsTo = "Consume water"
    }

    // Check hunger
    if b.PhysiologicalNeeds.Hunger > 30 {
        b.Owner.WantsTo = "Consume food"
    }
    
    // Check IsSufficientlyWarm
    if !b.PhysiologicalNeeds.IsSufficientlyWarm {
        b.Owner.WantsTo = "Get warm"
    }

    // Check NeedToExcrete
    if b.PhysiologicalNeeds.NeedToExcrete {
        b.Owner.WantsTo = "Excrete"
    }

    // Check IsInSafeArea
    if !b.PhysiologicalNeeds.IsInSafeArea {
        b.Owner.WantsTo = "Find a safe area"
    }

    // Check IsCapableOfDefendingSelf
    if !b.PhysiologicalNeeds.IsCapableOfDefendingSelf {
        b.Owner.WantsTo = "Improve defense"
    }

    // Check HasShelter
    if !b.PhysiologicalNeeds.HasShelter {
        b.Owner.WantsTo = "Make shelter"
    }

    // Check Rested
    if b.PhysiologicalNeeds.Rested < 20 {
        b.Owner.WantsTo = "Rest"
    }
}

// Translate the want to a list of tasks with priorities
func (b *Brain) TranslateWantToTaskList() {

    // Check if the person is breathing
    if !b.CheckIfCanBreath() {

        if b.Owner.Body.Head.Mouth == nil {
            return
        }
        if b.Owner.Body.Head.Mouth.IsObstructed {
            b.ActionList = append(b.ActionList, TargetedAction{"Clear airway", "Mouth", false,[]BodyPartType{"Hands"}, 100})
        }
        if b.Owner.Body.Head.Nose == nil {
            return
        }
        if b.Owner.Body.Head.Nose.IsObstructed {
            b.ActionList = append(b.ActionList, TargetedAction{"Clear airway", "Nose", false,[]BodyPartType{"Hands"}, 100})
        }
        if b.Owner.Body.Head.Nose.IsBroken {
            b.ActionList = append(b.ActionList, TargetedAction{"Fix nose", "Nose", false,[]BodyPartType{"Hands"}, 100})
        }
    }
}

// ----------------- Tasks ---------------------

// Rank the tasks based on priority
func (b *Brain) RankTasks() TargetedAction {
    highestPriority := 0
    highestPriorityAction := TargetedAction{"Idle", "Nothing", false, []BodyPartType{"Hands"}, 0}

    for _, action := range b.ActionList {
        if action.Priority > highestPriority {
            highestPriority = action.Priority
            highestPriorityAction = action
        }
    }

    return highestPriorityAction
}

func (b *Brain) performActions() {
    // Take the action with the highest priority
    action := b.RankTasks()

    // Perform the action
    switch action.Action {
    case "Idle":
        fmt.Println(b.Owner.FullName + " is idle.")
    case "Clear airway":
        b.ClearAirway(action.Target)
    case "Fix nose":
		// Decide if the person wants to fix the broken nose
        b.FixBrokenNose(action.Target)
    }
}


// Receive a requested task from another person
func (b *Brain) ReceiveTaskRequest(requestedTask RequestedAction) bool {
    fmt.Println(b.Owner.FullName + " received a task request from " + requestedTask.From.FullName)
    // Check the relationship between the two people
    hasRelationship := b.Owner.hasRelationship(requestedTask.From.FullName)

    // For now we will just accept the task
    if hasRelationship {
        if requestedTask.Action == "Talk" && b.Owner.IsTalking.IsActive {
            fmt.Println(b.Owner.FullName + " is already talking to someone.")
            return false
        } else if requestedTask.Action == "Talk" && !b.Owner.IsTalking.IsActive {
            b.Owner.IsTalking = TargetedAction{"Bla bla bla ...", requestedTask.From.FullName, true, make([]BodyPartType, 0), 10}
            fmt.Println(b.Owner.FullName + " accepted the task request from " + requestedTask.From.FullName)
            fmt.Println(b.Owner.FullName + " is talking to " + requestedTask.From.FullName)
            return true
        }
    } else {
        fmt.Println(b.Owner.FullName + " denied the task request from " + requestedTask.From.FullName + " because they are strangers.")
        return false
    }
    return false
}

// Send a task request to another person
func (b *Brain) SendTaskRequest(to *Person, taskType string) {
    if b.Owner.IsTalking.IsActive {
        fmt.Println(b.Owner.FullName + " is already talking to someone.")
        return 
    }
    fmt.Println(b.Owner.FullName + " is sending a task request to " + to.FullName)
    task := RequestedAction{TargetedAction{taskType, to.FullName, true, make([]BodyPartType, 0), 10}, b.Owner}
    success := to.Body.Head.Brain.ReceiveTaskRequest(task)
    if success {
        fmt.Println(to.FullName + " accepted the task request.")
        b.Owner.IsTalking = TargetedAction{"Hello " + to.FullName + ", how are you doing?", to.FullName, true, make([]BodyPartType, 0), 10}
        fmt.Println(b.Owner.FullName + " is talking to " + to.FullName)
    } else {
        fmt.Println(to.FullName + " declined the task request.")
    }
}

// ----------------- Safety ---------------------

// UnderAttack is called when the person is being attacked
func (b *Brain) UnderAttack(attacker *Person, targettedLimb BodyPartType, attackersLimb BodyPartType) {
	// Decide between fight or flight

	// Check if arms or hands are broken, if so, attack with legs
	if !b.Owner.Body.RightArm.IsBroken {
		b.Owner.AttackWithArm(attacker, targettedLimb, b.Owner.Body.RightArm.Hand)
	} else if !b.Owner.Body.LeftArm.IsBroken {
		b.Owner.AttackWithArm(attacker, targettedLimb, b.Owner.Body.LeftArm.Hand)
	} else if !b.Owner.Body.RightLeg.IsBroken {
		b.Owner.AttackWithLeg(attacker, targettedLimb, b.Owner.Body.RightLeg)
	} else if !b.Owner.Body.LeftLeg.IsBroken {
		b.Owner.AttackWithLeg(attacker, targettedLimb, b.Owner.Body.LeftLeg)
	} else {
		//
	}
}


// ----------------- Pathfinding ---------------------

// Decide a path to the target location - Check first if it's physically possible to walk.
func (b *Brain) DecidePathTo(x, y int) []*Node {
    path := b.AStar(b.Owner.Location.X, b.Owner.Location.Y, x, y)
    return path
}

// WalkOverPath - Walk over the path that was decided
func (b *Brain) WalkOverPath(x, y int) {
    path := b.DecidePathTo(x, y)
    if path == nil {
        fmt.Println(b.Owner.FullName + " could not find a path to the location.")
        return
    }
    for _, node := range path {
        b.Owner.WalkTo(node.X, node.Y)
    }
}

// ----------------- Tiles --------------------------

// Check what tile type the person is on
func (b *Brain) checkTileType() {
    // Check the tile type of the person
    b.Owner.OnTileType = b.Owner.WorldProvider.GetTileType(b.Owner.Location.X, b.Owner.Location.Y)
}