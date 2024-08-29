package main

import (
	"context"
	"fmt"
	"time"
)

type action struct {
	name string
    target string
	priority int
}
type RequestedAction struct {
    ActionType string
    Action     string
    FromPerson string

}
type Brain struct {
	owner  *Person
    active bool
    ctx    context.Context
    cancel context.CancelFunc
	actions []action
}

// NewBrain creates a new Brain and assigns an owner to it.
func NewBrain() *Brain {
    ctx, cancel := context.WithCancel(context.Background())
    return &Brain{
        active:  false,
        ctx:     ctx,
        cancel:  cancel,
        actions: []action{
            {"Idle", "", 1},
        },
    }
}

func (b *Brain) turnOn() {
    if b.active {
        fmt.Println("Brain is already active.")
        return
    }

    fmt.Println("Brain for: " + b.owner.FullName + " is now active.")
    b.owner.IsConscious = true
    b.active = true

    go b.mainLoop()
}

func (b *Brain) turnOff() {
    if !b.active {
        fmt.Println("Brain is already inactive.")
        return
    }

    fmt.Println(b.owner.FullName, "'s brain is shutting down")
    b.cancel()
    b.active = false
}

func (b *Brain) mainLoop() {
    fmt.Println(b.owner.FullName + "'s brain is now ", b.active)

    for {
        select {
        case <-b.ctx.Done():
            b.active = false
            return
        default:
        
        if !b.owner.IsConscious{
            fmt.Println(b.owner.FullName + "' brain is not conscious but still alive.")
            return
        } else {
            // Brain logic goes here

            fmt.Println(b.owner.FullName + "'s brain is thinking...")

            obs := b.processInputs()
            b.makeDecisions(obs)
            b.performActions()
        }

            // Sleep or yield for a bit to prevent CPU hogging
            time.Sleep(2000 * time.Millisecond)
        }
    }
    
}

func (b *Brain) processInputs() Vision {
    // This will probably have to be the WorldState struct but a smaller area
	// For now we could just decide if the person is in a friendly or hostile area
    
    // Get the vision of the person
    obs := b.owner.GetVision()

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
    b.owner.OnTileType = b.owner.WorldProvider.GetTileType(b.owner.Location.X, b.owner.Location.Y)
}

// Decide if the area is safe or not
func (b *Brain) isAreaSafe(obs Vision) {
    // Loop through the observations and make decisions based on relationships with the people

    collectiveIntensity := 0
    numberOfPeople := len(obs.Persons)

    for _, person := range obs.Persons {
        if b.owner.hasRelationship(person.FullName) {
            for _, relationship := range b.owner.Relationships {
                if relationship.WithPerson == person.FullName {
                    collectiveIntensity += relationship.Intensity
                }
            }
        }
    }

    // This is a pretty dumb way to determine if the area is safe or not, but it's a start
    b.owner.FeelingSafe = collectiveIntensity/numberOfPeople
}

func (b *Brain) makeDecisions(obs Vision) {
    // Check if we're engaging in conversation with someone and if we are and we dont have that person in the observation, we should cancel the conversation
    if b.owner.IsTalking.IsActive {
        if !obs.HasPerson(b.owner.IsTalking.Target) {
            fmt.Println(b.owner.FullName + " is no longer talking to " + b.owner.IsTalking.Target)
            b.owner.IsTalking = TargetedAction{"", "", false, make([]string, 0)}
        }
    }
    // Loop through the observations and make decisions based on people
    for _, person := range obs.Persons {
        if person.FullName != b.owner.FullName {
            if b.owner.hasRelationship(person.FullName) {
                for _, relationship := range b.owner.Relationships {
                    if relationship.WithPerson == person.FullName {
                        relationship.Intensity++
                        if relationship.Intensity > 3 { // This should be a constant and above 15 and below 40
                            relationship.Relationship = "Aquantance"
                            if !b.owner.IsTalking.IsActive {
                                targetPerson := b.owner.GetPersonByFullName(relationship.WithPerson)
                                b.SendTaskRequest(targetPerson, "Talk")
                            }
                        }
                        b.owner.updateRelationship(person.FullName, relationship.Relationship, relationship.Intensity)
                    }
                }
            } else { // If the person does not have a relationship with the other person
                b.owner.addRelationship(person, "Stranger", 0)
            }
        }
    }
}

// Receive a requested task from another person
func (b *Brain) ReceiveTaskRequest(requestedTask RequestedAction, from *Person) bool {
    fmt.Println(b.owner.FullName + " received a task request from " + from.FullName)
    // Check the relationship between the two people
    hasRelationship := b.owner.hasRelationship(from.FullName)

    // For now we will just accept the task
    if hasRelationship {
        if requestedTask.ActionType == "Talk" && b.owner.IsTalking.IsActive {
            fmt.Println(b.owner.FullName + " is already talking to someone.")
            return false
        } else if requestedTask.ActionType == "Talk" && !b.owner.IsTalking.IsActive {
            b.owner.IsTalking = TargetedAction{"Bla bla bla ...", from.FullName, true, make([]string, 0)}
            fmt.Println(b.owner.FullName + " accepted the task request from " + from.FullName)
            fmt.Println(b.owner.FullName + " is talking to " + from.FullName)
            return true
        }
    } else {
        fmt.Println(b.owner.FullName + " denied the task request from " + from.FullName + " because they are strangers.")
        return false
    }
    return false
}

// Send a task request to another person
func (b *Brain) SendTaskRequest(to *Person, taskType string) {
    if b.owner.IsTalking.IsActive {
        fmt.Println(b.owner.FullName + " is already talking to someone.")
        return 
    }
    fmt.Println(b.owner.FullName + " is sending a task request to " + to.FullName)
    success := to.Brain.ReceiveTaskRequest(RequestedAction{taskType, "Hello!", b.owner.FullName}, b.owner)
    if success {
        fmt.Println(to.FullName + " accepted the task request.")
        b.owner.IsTalking = TargetedAction{"Hello " + to.FullName + ", how are you doing?", to.FullName, true, make([]string, 0)}
        fmt.Println(b.owner.FullName + " is talking to " + to.FullName)
    } else {
        fmt.Println(to.FullName + " declined the task request.")
    }
}

func (b *Brain) performActions() {
    if b.owner.IsTalking.IsActive {
        fmt.Println(b.owner.FullName + " says " + b.owner.IsTalking.Action + " to " + b.owner.IsTalking.Target)
    }
}
