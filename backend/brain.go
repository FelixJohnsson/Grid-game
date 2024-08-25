package main

import (
	"context"
	"fmt"
	"time"
)

type action struct {
	name string
	priority int
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
            {"Idle", 1},
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
			b.processInputs([]string{"input1", "input2", "input3"})

            // Sleep or yield for a bit to prevent CPU hogging
            time.Sleep(1000 * time.Millisecond)
        }
    }
}

func (b *Brain) processInputs(inputs []string) {
    // This will probably have to be the WorldState struct but a smaller area
	// For now we could just decide if the person is in a friendly or hostile area
    fmt.Println(b.owner.Name + " is processing inputs...")
}

func (b *Brain) makeDecisions() {
    // Example: Decide what to do based on inputs and state
    fmt.Println("Making decisions...")
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
