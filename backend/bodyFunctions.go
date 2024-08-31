package main

import (
	"fmt"
	"math/rand"
)

// ---------------- Actions ----------------

// ClearAirway - Clear the airway of the person
func (b *Brain) ClearAirway(target string) {
    randomNumber := rand.Intn(100)
	fmt.Println("Trying to clear airway of the mouth")

    if target == "Mouth" && randomNumber < 2 {
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
func (b *Brain) FixBrokenNose(target string) {
    randomNumber := rand.Intn(100)
	fmt.Println("Trying to fix the nose")

    if randomNumber < 2 {
        b.Owner.Body.Head.Nose.IsBroken = false
        // Remove the action from the action list
        for i := len(b.ActionList) - 1; i >= 0; i-- {
            if b.ActionList[i].Action == "Fix nose" && b.ActionList[i].Target == "Nose" {
                b.ActionList = append(b.ActionList[:i], b.ActionList[i+1:]...)
            }
        }
        fmt.Println(b.Owner.FullName + " fixed the nose.")
        b.ApplyPain(101)
    }
}
