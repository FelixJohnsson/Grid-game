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
            b.processInputs()

            // Sleep or yield for a bit to prevent CPU hogging
            time.Sleep(1000 * time.Millisecond)
        }
    }
}

func (b *Brain) processInputs() {
    // Example: Check surroundings, inputs, etc.
    fmt.Println("Processing inputs...")
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
