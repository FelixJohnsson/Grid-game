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

    b.IsConscious = true
    b.Active = true
    switch b.Owner.Species {
    case "Wolf":
        go b.WolfMainLoop()
    case "Homo Sapiens":
        go b.HomoSapiensMainLoop()
    }

}

func (b *Brain) turnOff() {
    if !b.Active {
        fmt.Println("Brain is already inactive.")
        return
    }

    fmt.Println(b.Owner.FullName + "'s brain is shutting down")
	b.Active = false
	b.IsConscious = false
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
    b.CalculatePainLevel()

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

// CheckIfWantIsAlreadyInList - Check if the want is already in the list
func (b *Brain) CheckIfWantIsAlreadyInList(want string) bool {
    for _, w := range b.Owner.WantsTo {
        if w == want {
            return true
        }
    }
    return false
}

// ClearWants - Clear the wants of the person
func (b *Brain) ClearWants() {
    b.Owner.WantsTo = make([]string, 0)
}

// HomoSapiensCalculateWant - Calculate the want of the person
func (b *Brain) HomoSapiensCalculateWant() {
    switch {
    case !b.CheckIfCanBreath() && !b.CheckIfWantIsAlreadyInList("Be able to breath"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Be able to breath")
    case b.PhysiologicalNeeds.IsInPain && !b.CheckIfWantIsAlreadyInList("Relieve pain"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Relieve pain")
    case b.PhysiologicalNeeds.Thirst > 30 && !b.CheckIfWantIsAlreadyInList("Consume water"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Consume water")
    case b.PhysiologicalNeeds.Hunger > 30 && !b.CheckIfWantIsAlreadyInList("Consume food"):
        fmt.Println(b.Owner.FullName + " wants to consume food.", b.PhysiologicalNeeds.Hunger)
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Consume food")
    case !b.PhysiologicalNeeds.IsSufficientlyWarm && !b.CheckIfWantIsAlreadyInList("Get warm"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Get warm")
    case b.PhysiologicalNeeds.NeedToExcrete && !b.CheckIfWantIsAlreadyInList("Excrete"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Excrete")
    case !b.PhysiologicalNeeds.IsInSafeArea && !b.CheckIfWantIsAlreadyInList("Find a safe area"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Find a safe area")
    case !b.PhysiologicalNeeds.IsCapableOfDefendingSelf && !b.CheckIfWantIsAlreadyInList("Improve defense"):
        if b.Owner.CombatSkill > 30 {
            b.PhysiologicalNeeds.IsCapableOfDefendingSelf = true
        }
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Improve defense")
    case !b.PhysiologicalNeeds.HasShelter && !b.CheckIfWantIsAlreadyInList("Make shelter"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Make shelter")
    case b.PhysiologicalNeeds.Rested < 20 && !b.CheckIfWantIsAlreadyInList("Rest"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Rest")
    }
}

// AnimalCalculateWant - Calculate the want of the person
func (b *Brain) AnimalCalculateWant() {
    switch {
    case !b.CheckIfCanBreath() && !b.CheckIfWantIsAlreadyInList("Be able to breath"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Be able to breath")
    case b.PhysiologicalNeeds.IsInPain && !b.CheckIfWantIsAlreadyInList("Relieve pain"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Relieve pain")
    case b.PhysiologicalNeeds.Thirst > 30 && !b.CheckIfWantIsAlreadyInList("Consume water"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Consume water")
    case b.PhysiologicalNeeds.Hunger > 30 && !b.CheckIfWantIsAlreadyInList("Consume food"):
        fmt.Println(b.Owner.FullName + " wants to consume food.", b.PhysiologicalNeeds.Hunger)
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Consume food")
    case !b.PhysiologicalNeeds.IsSufficientlyWarm && !b.CheckIfWantIsAlreadyInList("Get warm"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Get warm")
    case b.PhysiologicalNeeds.NeedToExcrete && !b.CheckIfWantIsAlreadyInList("Excrete"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Excrete")
    case !b.PhysiologicalNeeds.IsInSafeArea && !b.CheckIfWantIsAlreadyInList("Find a safe area"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Find a safe area")
    case !b.PhysiologicalNeeds.HasShelter && !b.CheckIfWantIsAlreadyInList("Find shelter"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Make shelter")
    case b.PhysiologicalNeeds.Rested < 20 && !b.CheckIfWantIsAlreadyInList("Rest"):
        b.Owner.WantsTo = append(b.Owner.WantsTo, "Rest")
    }
}

// IsTaskInActionList - Check if the task is already in the action list
func (b *Brain) IsTaskInActionList(task TargetedAction) bool {
	for _, action := range b.ActionList {
		if action.Action == task.Action && action.Target == task.Target {
			return true
		}
	}
	return false
}

//RemoveActionFromActionList - Remove an action from the action list
func (b *Brain) RemoveActionFromActionList(action TargetedAction) {
	for i, a := range b.ActionList {
		if a.Action == action.Action && a.Target == action.Target {
			b.ActionList = append(b.ActionList[:i], b.ActionList[i+1:]...)
		}
	}
}

// ClearCurrentTask - Clear the current task
func (b *Brain) ClearCurrentTask() {
    b.CurrentTask = TargetedAction{"Idle", "Nothing", false, []BodyPartType{"Hands"}, 0}
}

// AddTaskToActionList - Add a task to the action list
func (b *Brain) AddTaskToActionList(task TargetedAction) {
	for i, action := range b.ActionList {
		if action.Action == "Idle" {
			b.ActionList = append(b.ActionList[:i], b.ActionList[i+1:]...)
		}
	}
	b.ActionList = append(b.ActionList, task)
}

// ClearActionList - Clear the action list
func (b *Brain) ClearActionList() {
    b.ActionList = make([]TargetedAction, 0)
}

// Translate the want to a list of tasks with priorities
func (b *Brain) TranslateWantToTaskList() {
    b.ClearActionList()
    if !b.CheckIfCanBreath() {
        if b.Owner.Body.Head.Mouth.IsObstructed {
            action := TargetedAction{"Clear airway", "Mouth", false,[]BodyPartType{"Hands"}, 100}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
        }
        if b.Owner.Body.Head.Nose.IsObstructed {
            action := TargetedAction{"Clear airway", "Nose", false,[]BodyPartType{"Hands"}, 100}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
        }
        if b.Owner.Body.Head.Nose.IsBroken {
            action := TargetedAction{"Fix nose", "Nose", false,[]BodyPartType{"Hands"}, 100}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
        }
    }
    if b.PhysiologicalNeeds.Thirst > 30 {
		action := TargetedAction{"Drink water", "", false,[]BodyPartType{"Hands"}, 99}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if b.PhysiologicalNeeds.Hunger > 30 {
		action := TargetedAction{"Eat food", "", false,[]BodyPartType{"Hands"}, 98}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if b.PhysiologicalNeeds.IsInPain {
		action := TargetedAction{"Reduce pain", "", false,[]BodyPartType{"Hands"}, 95}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}	
	}
	if !b.PhysiologicalNeeds.WayOfGettingWater {
		action := TargetedAction{"Find a water supply", "", false,[]BodyPartType{"Hands"}, 90}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if !b.PhysiologicalNeeds.WayOfGettingFood {
		action := TargetedAction{"Find a food supply", "", false,[]BodyPartType{"Hands"}, 90}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if !b.PhysiologicalNeeds.FoodSupply {
		action := TargetedAction{"Have food for storage", "", false,[]BodyPartType{"Hands"}, 90}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}

	if !b.PhysiologicalNeeds.IsSufficientlyWarm {
		action := TargetedAction{"Get warm", "", false,[]BodyPartType{"Hands"}, 85}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if b.PhysiologicalNeeds.NeedToExcrete {
		action := TargetedAction{"Excrete", "", false,[]BodyPartType{"Hands"}, 80}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if !b.PhysiologicalNeeds.IsInSafeArea {
		action := TargetedAction{"Find a safe area", "", false,[]BodyPartType{"Hands"}, 75}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}

	if !b.PhysiologicalNeeds.HasShelter {
		action := TargetedAction{"Make shelter", "", false,[]BodyPartType{"Hands"}, 70}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
    if b.PhysiologicalNeeds.Rested < 20 {
		action := TargetedAction{"Rest", "", false,[]BodyPartType{"Hands"}, 65}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
    if !b.PhysiologicalNeeds.IsCapableOfDefendingSelf {
		action := TargetedAction{"Improve defense", "", false,[]BodyPartType{"Hands"}, 60}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
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

    fmt.Println(b.Owner.FullName + " is performing the action: " + action.Action)

    // Perform the action
    switch action.Action {
    case "Drink water":
        b.CurrentTask = action
        if b.MotorCortexCurrentTask.ActionReason != "Drink water" && !b.MotorCortexCurrentTask.IsActive {
            b.MotorCortexFindDrinkWater()
            return
        } else if b.MotorCortexCurrentTask.ActionReason == "Drink water" && !b.MotorCortexCurrentTask.Finished && b.MotorCortexCurrentTask.IsActive {
            return
        }

        if b.Owner.isStandingOnWater() && b.MotorCortexCurrentTask.ActionReason == "Drink water" && b.MotorCortexCurrentTask.Finished && !b.MotorCortexCurrentTask.IsActive {
            water := Liquid{"Water"}
            b.Owner.Drink(water)
            b.AddMemoryToLongTerm("Found water supply", "Water", b.Owner.Location)
            b.PhysiologicalNeeds.WayOfGettingWater = true
            b.MotorCortexCurrentTask = MotorCortexAction{"None", "Idle", Location{0, 0}, false, false}
        }

    case "Eat food":
        b.CurrentTask = action
        b.EatFruit()
        b.ClearCurrentTask()
	case "Clear airway":
        b.CurrentTask = action
        b.ClearAirway(action)
        b.ClearCurrentTask()
		return
    case "Fix nose":
        b.CurrentTask = action
        b.FixBrokenNose(action)
        b.ClearCurrentTask()
		return
	case "Reduce pain":
        b.CurrentTask = action
		//b.ReducePain()
        b.ClearCurrentTask()
		return
	case "Find a water supply":
        b.CurrentTask = action
		b.FindWaterSupply(action)
        b.ClearCurrentTask()
		return
	case "Find a food supply":
        b.CurrentTask = action
		b.FindFoodSupply(action)
        b.ClearCurrentTask()
		return
	case "Have food for storage":
        b.CurrentTask = action
		b.GetFoodForStorage(action)
        b.ClearCurrentTask()
    case "Find shelter":
        b.CurrentTask = action
        b.FindAndNote("Shelter")
        b.ClearCurrentTask()
        return
    case "Make shelter":
        fmt.Println(b.Owner.FullName + " is making a shelter.")
        b.CurrentTask = action
        b.MakeShelter(action)
        b.ClearCurrentTask()
        return
    case "Improve defense":
        fmt.Println(b.Owner.FullName + " is improving defense.")
        b.CurrentTask = action
        b.ImproveDefense()
    case "Idle":
        b.CurrentTask = action
        fmt.Println(b.Owner.FullName + " is idle.")
		return

    }
}


// Receive a requested task from another person
func (b *Brain) ReceiveTaskRequest(requestedTask RequestedAction) bool {
    fmt.Println(b.Owner.FullName + " received a task request from " + requestedTask.From.FullName)
    
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
    success := to.Brain.ReceiveTaskRequest(task)
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
func (b *Brain) WalkOverPath(x, y int) bool {
    path := b.DecidePathTo(x, y)
    if path == nil {
        fmt.Println(b.Owner.FullName + " could not find a path to the location.")
        return false
    }
    fmt.Println(b.Owner.FullName + " is walking to ", x, y)
    for _, node := range path {
		// Wait for a half second before walking to the next node
		time.Sleep(500 * time.Millisecond)
        b.Owner.WalkTo(node.X, node.Y)
    }
	b.MotorCortexCurrentTask = MotorCortexAction{"None", "Idle", Location{0, 0}, false, true}
    return true
}

// ----------------- Items -------------------------

// FindInOwnedItems - Find an item in the owned items
func (b *Brain) FindInOwnedItems(itemName string) *Item {
    for _, item := range b.Owner.OwnedItems {
        if item.Name == itemName {
            return item
        }
    }
    return nil
}

// HasItemEquippedInRight - Check if the person has an item equipped in right hand
func (b *Brain) HasItemEquippedInRight(itemName string) bool {
    for _, item := range b.Owner.Body.RightArm.Hand.Items {
        if item.Name == itemName {
            return true
        }
    }
    return false
}

// HasItemEquippedInLeft - Check if the person has an item equipped in left hand
func (b *Brain) HasItemEquippedInLeft(itemName string) bool {
    for _, item := range b.Owner.Body.LeftArm.Hand.Items {
        if item.Name == itemName {
            return true
        }
    }
    return false
}