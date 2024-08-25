package main

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type action struct {
	name string
    target string
	priority int
}

type RequestedAction struct {
    action string
    fromPerson string
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

    fmt.Println("Brain is now active.")
    b.active = true

    go b.mainLoop()
}

func (b *Brain) addTask(task action) {
	b.actions = append(b.actions, task)
}

func (b *Brain) mainLoop() {
    for {
        select {
        case <-b.ctx.Done():
            fmt.Println("Brain is shutting down.")
            b.active = false
            return
        default:
            // Brain logic goes here
            fmt.Println("Brain is thinking...")
			obs := b.processInputs()
            b.makeDecisions(obs)

            // Sleep or yield for a bit to prevent CPU hogging
            time.Sleep(5000 * time.Millisecond)
        }
    }
}

func (b *Brain) processInputs() Vision {
    // This will probably have to be the WorldState struct but a smaller area
	// For now we could just decide if the person is in a friendly or hostile area
    
    // Get the vision of the person
    world, _ := loadWorldFromFile()
    vision := world.GetVision(b.owner.Location.X, b.owner.Location.Y, b.owner.VisionRange)

    return vision
}

func (b *Brain) makeDecisions(obs Vision) {
    // Loop through the observations and make decisions based on people
    for _, person := range obs.Persons {
        if person.FullName != b.owner.FullName {
            if b.owner.hasRelationship(person.FullName) {
                // Get the relationship intensity
                for _, relationship := range b.owner.Relationships {
                    if relationship.WithPerson == person.FullName {
                        relationship.Intensity++
                        if relationship.Intensity > 15 {
                            relationship.Relationship = "Aquantance"
                        } 
                        if relationship.Intensity > 40 {
                            relationship.Relationship = "Friend"
                        } 
                        if relationship.Intensity > 60 {
                            relationship.Relationship = "Close Friend"
                        } 
                        if relationship.Intensity > 80 {
                            relationship.Relationship = "Best Friend"
                        } 
                        if relationship.Intensity > 90 {
                            relationship.Relationship = "Family"
                        } 
                        if relationship.Intensity > 100 {
                            relationship.Relationship = "Soulmate"
                            relationship.Intensity = 100
                        }

                        b.owner.updateRelationship(person.FullName, relationship.Relationship, relationship.Intensity)
                    }
                }
            } else {
                b.owner.addRelationship(person, "Stranger", 0)
                if (b.owner.Personality == "Talkative") {
                    b.owner.addTask(action{"Talk", person.FullName, -2})
                    // Request action from the other person
                    

                    
                }
            }
        } else {
            fmt.Println("Found myself")
        }
    }
    fmt.Println(b.owner.Relationships)
}

// RequestTask requests a task from the brain.
func (b *Brain) RequestTask(requestedTask RequestedAction) bool {
    if !b.active {
        fmt.Println("Brain is inactive.")
        return false
    }

    // Check relationship with the person
    for _, relationship := range b.owner.Relationships {
        if relationship.WithPerson == requestedTask.fromPerson {
            if relationship.Relationship == "Stranger" {
                if b.owner.Personality == "Talkative" {
                    b.addTask(action{"Talk", requestedTask.fromPerson, -2})
                    return true
                }
                if b.owner.Personality == "Shy" {
                    return false
                }
                // Generate a random number between 1 and 10
                random := gofakeit.Number(1, 10)
                if random > 5 {
                    b.addTask(action{"Talk", requestedTask.fromPerson, 5})
                    return true
                }
            }
            if relationship.Relationship == "Aquantance" {
                if b.owner.Personality == "Talkative" {
                    b.addTask(action{"Talk", requestedTask.fromPerson, -2})
                    return true
                }
                if b.owner.Personality == "Shy" {
                    random := gofakeit.Number(1, 10)
                    if random > 8 {
                        b.addTask(action{"Talk", requestedTask.fromPerson, 10})
                        return true
                    }
                    return false
                }
                // Generate a random number between 1 and 10
                random := gofakeit.Number(1, 10)
                if random > 3 {
                    b.addTask(action{"Talk", requestedTask.fromPerson, 5})
                    return true
                }
            }
            if relationship.Relationship == "Friend" {
                if b.owner.Personality == "Talkative" {
                    b.addTask(action{"Talk", requestedTask.fromPerson, -2})
                    return true
                }
                if b.owner.Personality == "Shy" {
                    random := gofakeit.Number(1, 10)
                    if random > 5 {
                        b.addTask(action{"Talk", requestedTask.fromPerson, 7})
                        return true
                    }
                    return false
                }
                // Generate a random number between 1 and 10
                random := gofakeit.Number(1, 10)
                if random > 1 {
                    b.addTask(action{"Talk", requestedTask.fromPerson, 5})
                    return true
                }
            }
            if relationship.Relationship == "Close Friend" {
                if b.owner.Personality == "Talkative" {
                    b.addTask(action{"Talk", requestedTask.fromPerson, -10})
                    return true
                }
                b.addTask(action{"Talk", requestedTask.fromPerson, 0})
                return true
            }
            if relationship.Relationship == "Best Friend" {
                if b.owner.Personality == "Talkative" {
                    b.addTask(action{"Talk", requestedTask.fromPerson, -10})
                    return true
                }
                b.addTask(action{"Talk", requestedTask.fromPerson, -5})
                return true
            }
            if relationship.Relationship == "Family" {
                if b.owner.Personality == "Talkative" {
                    b.addTask(action{"Talk", requestedTask.fromPerson, -10})
                    return true
                }
                b.addTask(action{"Talk", requestedTask.fromPerson, -5})
                return true
            }
            if relationship.Relationship == "Soulmate" {
                if b.owner.Personality == "Talkative" {
                    b.addTask(action{"Talk", requestedTask.fromPerson, -15})
                    return true
                }
                b.addTask(action{"Talk", requestedTask.fromPerson, -10})
                return true
            }
        }
    }
    return false
}

func (b *Brain) performActions() {
    // Example: Perform the decided actions
    fmt.Println("Performing actions...")
}

func (b *Brain) turnOff() {
    if !b.active {
        fmt.Println("Brain is already inactive.")
        return
    }

    fmt.Println("Shutting down brain.")
    b.cancel()
    b.active = false
}