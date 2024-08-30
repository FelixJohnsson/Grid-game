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
        Actions: []TargetedAction{
            {"Idle", "", false, make([]string, 0)},
        },
        IsConscious: true,
        IsAlive:    true,
        BrainDamage: 0,
        IsUnderAttack: IsUnderAttack{false, nil, "", ""}, 
        Memories: Memories{make([]Memory, 0), make([]Memory, 0)},
    }
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

    fmt.Println(b.Owner.FullName, "'s brain is shutting down")
    b.Cancel()
}

func (b *Brain) mainLoop() {
    fmt.Println(b.Owner.FullName + "'s brain is now ", b.Active)

    for {
        select {
        case <-b.Ctx.Done():
            b.Active = false
            return
        default:
        
        if !b.IsConscious{
            fmt.Println(b.Owner.FullName + "' brain is not conscious but still alive.")
            return
        } else {
            // Brain logic goes here

            fmt.Println(b.Owner.FullName + "'s brain is thinking...")

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
            b.performActions()
        }

            // Sleep or yield for a bit to prevent CPU hogging
            time.Sleep(2000 * time.Millisecond)
        }
    }
    
}

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
func (b *Brain) UnderAttack(attacker *Person, targettedLimb LimbType, attackersLimb LimbType) {
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
    // This will probably have to be the WorldState struct but a smaller area
	// For now we could just decide if the person is in a friendly or hostile area
    
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
            b.Owner.IsTalking = TargetedAction{"", "", false, make([]string, 0)}
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
            b.Owner.IsTalking = TargetedAction{"Bla bla bla ...", requestedTask.From.FullName, true, make([]string, 0)}
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
    task := RequestedAction{TargetedAction{taskType, to.FullName, true, make([]string, 0)}, b.Owner}
    success := to.Body.Head.Brain.ReceiveTaskRequest(task)
    if success {
        fmt.Println(to.FullName + " accepted the task request.")
        b.Owner.IsTalking = TargetedAction{"Hello " + to.FullName + ", how are you doing?", to.FullName, true, make([]string, 0)}
        fmt.Println(b.Owner.FullName + " is talking to " + to.FullName)
    } else {
        fmt.Println(to.FullName + " declined the task request.")
    }
}

func (b *Brain) performActions() {
    if b.Owner.IsTalking.IsActive {
        fmt.Println(b.Owner.FullName + " says " + b.Owner.IsTalking.Action + " to " + b.Owner.IsTalking.Target)
    }
}
