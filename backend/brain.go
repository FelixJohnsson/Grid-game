package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

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

func (b *Brain) TranslateWantToTaskList() {
    // Translate the want to a list of tasks with priorities

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

// CheckIfCanBreah - Check if the person can breath
func (b *Brain) CheckIfCanBreath() bool {
    mouthCanBreath := b.Owner.Body.Head.Mouth != nil && !b.Owner.Body.Head.Mouth.IsObstructed
    noseCanBreath := b.Owner.Body.Head.Nose != nil && !b.Owner.Body.Head.Nose.IsObstructed && !b.Owner.Body.Head.Nose.IsBroken

    return mouthCanBreath || noseCanBreath
}

// Breath 
func (b *Brain) Breath() {
    b.OxygenLevel += 10
}

// ConsumeOxygen
func (b *Brain) ConsumeOxygen() {
    b.OxygenLevel -= 10
}

// CalculatePainLevel - Calculate the pain level of the person
func (b *Brain) CalculatePainLevel() {
    // Check the current situation of the person

    // Check if the person is in pain
    // We need to loop over the body parts and check if it's broken or bleeding

    if b.Owner.Body.Head != nil {
        if b.Owner.Body.Head.IsBroken {
            b.PainLevel += 5
        } 
        if b.Owner.Body.Head.IsBleeding {
            b.PainLevel += 2
        }
        if b.Owner.Body.Head.Ears != nil && b.Owner.Body.Head.Ears.IsBleeding || b.Owner.Body.Head.Ears.IsBroken {
            b.PainLevel += 1
        } 
        if b.Owner.Body.Head.Eyes != nil && b.Owner.Body.Head.Eyes.IsBleeding || b.Owner.Body.Head.Eyes.IsBroken {
            b.PainLevel += 5
        }
        if b.Owner.Body.Head.Nose != nil && b.Owner.Body.Head.Nose.IsBleeding || b.Owner.Body.Head.Nose.IsBroken {
            b.PainLevel += 2
        }
        if b.Owner.Body.Head.Mouth != nil  &&b.Owner.Body.Head.Mouth.IsBleeding || b.Owner.Body.Head.Mouth.IsBroken {
            b.PainLevel += 2
        }
    }
    if b.Owner.Body.Torso.IsBleeding || b.Owner.Body.Torso.IsBroken {
        b.PainLevel += 5
    }
    if b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.IsBleeding || b.Owner.Body.RightArm.IsBroken {
        b.PainLevel += 5
    }
    if b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.IsBleeding || b.Owner.Body.LeftArm.IsBroken {
        b.PainLevel += 5
    }
    if b.Owner.Body.RightLeg != nil && b.Owner.Body.RightLeg.IsBleeding || b.Owner.Body.RightLeg.IsBroken {
        b.PainLevel += 5
    }
    if b.Owner.Body.LeftLeg != nil && b.Owner.Body.LeftLeg.IsBleeding || b.Owner.Body.LeftLeg.IsBroken {
        b.PainLevel += 5
    }
    if b.Owner.Body.RightArm != nil && b.Owner.Body.RightArm.Hand != nil && b.Owner.Body.RightArm.Hand.IsBleeding || b.Owner.Body.RightArm.Hand.IsBroken {
        b.PainLevel += 2
    }
    if b.Owner.Body.LeftArm != nil && b.Owner.Body.LeftArm.Hand != nil && b.Owner.Body.LeftArm.Hand.IsBleeding || b.Owner.Body.LeftArm.Hand.IsBroken {
        b.PainLevel += 2
    }
    if b.Owner.Body.RightLeg != nil && b.Owner.Body.RightLeg.Foot != nil && b.Owner.Body.RightLeg.Foot.IsBleeding || b.Owner.Body.RightLeg.Foot.IsBroken {
        b.PainLevel += 2
    }
    if b.Owner.Body.LeftLeg != nil && b.Owner.Body.LeftLeg.Foot != nil && b.Owner.Body.LeftLeg.Foot.IsBleeding || b.Owner.Body.LeftLeg.Foot.IsBroken {
        b.PainLevel += 2
    }
}

func (b *Brain) LooseConsciousness() {
    // Should wait on another thread for 10 seconds before the person regains consciousness
    b.IsConscious = false
    go func() {
        time.Sleep(10 * time.Second)
        fmt.Println(b.Owner.FullName + " regained consciousness.")
        b.IsConscious = true
    }()
}

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

func (b *Brain) mainLoop() {
    fmt.Println(b.Owner.FullName + "'s brain is now ", b.Active)

    // This has to be broken up into "Bodily functions" and "Brain functions"
    for {
        select {
        case <-b.Ctx.Done():
            b.Active = false
            return
        default:
        if b.CheckIfCanBreath() {
            b.Breath()
        } else {
            fmt.Println(b.Owner.FullName + " is not able to breath.")
            b.CanBreath = false
        }
        b.ConsumeOxygen()
        b.CalculatePainLevel()

        if b.OxygenLevel <= 0 {
            b.Owner.Body.Head.Brain.turnOff()
            fmt.Println(b.Owner.FullName + "'s brain is shutting down due to lack of oxygen.")
            return
        }

        if b.PainLevel > b.PainTolerance {
            b.LooseConsciousness()
            fmt.Println(b.Owner.FullName + "'s brain lost consciousness due to pain.")
        }
        
        if !b.IsConscious{
            fmt.Println(b.Owner.FullName + "'s brain is not conscious but still alive.")
            return
        } else {
            // Brain logic goes here
            if  b.IsUnderAttack.Active && b.IsUnderAttack.From.Body.Head == nil {
                // The attacker is dead
                b.AddMemoryToLongTerm("Killed " + b.IsUnderAttack.From.FullName, b.IsUnderAttack.From.Location)

                b.IsUnderAttack = IsUnderAttack{false, nil, "", ""}
                fmt.Println(b.Owner.FullName + " is no longer under attack.")
            }

            // Check if the person is under attack
            if b.IsUnderAttack.Active{
                b.UnderAttack(b.IsUnderAttack.From, b.IsUnderAttack.Target, b.IsUnderAttack.ByLimb)
                b.Owner.UpdateRelationship(b.IsUnderAttack.From.FullName, "Enemy", 100)
                b.AddMemoryToLongTerm("Under Attack", b.IsUnderAttack.From.Location)
            }

            obs := b.processInputs()
            b.makeDecisions(obs)
            b.CalculateWant()
            b.TranslateWantToTaskList()
            b.performActions()
        }

            // Sleep or yield for a bit to prevent CPU hogging
            time.Sleep(2000 * time.Millisecond)
        }
    }
}

// TaskHandler is a function that handles the tasks that the brain receives


// AddMemoryToShortTerm adds a memory to the short term memory
func (b *Brain) AddMemoryToShortTerm(event string, location Location) {
    memory := Memory{event, location}
    b.Memories.ShortTermMemory = append(b.Memories.ShortTermMemory, memory)
}

// AddMemoryToLongTerm adds a memory to the long term memory
func (b *Brain) AddMemoryToLongTerm(event string, location Location) {
    memory := Memory{event, location}
    b.Memories.LongTermMemory = append(b.Memories.LongTermMemory, memory)
}

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

// Check what tile type the person is on
func (b *Brain) checkTileType() {
    // Check the tile type of the person
    b.Owner.OnTileType = b.Owner.WorldProvider.GetTileType(b.Owner.Location.X, b.Owner.Location.Y)
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

    fmt.Println(b.Owner.FullName + " is performing action: " + action.Action)

    // Perform the action
    switch action.Action {
    case "Idle":
        fmt.Println(b.Owner.FullName + " is idle.")
    case "Clear airway":
        b.ClearAirway(action.Target)
    case "Fix nose":
        b.FixNose(action.Target)
    }
}

// ClearAirway - Clear the airway of the person
func (b *Brain) ClearAirway(target string) {
    randomNumber := rand.Intn(100)

    if target == "Mouth" && randomNumber < 50 {
        b.Owner.Body.Head.Mouth.IsObstructed = false
        // Remove the action from the action list
        for i := len(b.ActionList) - 1; i >= 0; i-- {
            if b.ActionList[i].Action == "Clear airway" && b.ActionList[i].Target == "Mouth" {
                b.ActionList = append(b.ActionList[:i], b.ActionList[i+1:]...)
            }
        }
        fmt.Println(b.Owner.FullName + " cleared the airway of the mouth.")
    }
}

// FixNose - Fix the nose of the person
func (b *Brain) FixNose(target string) {
    randomNumber := rand.Intn(100)

    if randomNumber < 50 {
        b.Owner.Body.Head.Nose.IsBroken = false
        // Remove the action from the action list
        for i := len(b.ActionList) - 1; i >= 0; i-- {
            if b.ActionList[i].Action == "Fix nose" && b.ActionList[i].Target == "Nose" {
                b.ActionList = append(b.ActionList[:i], b.ActionList[i+1:]...)
            }
        }
        b.PainLevel += 120
        fmt.Println(b.Owner.FullName + " fixed the nose.")
    }
}

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