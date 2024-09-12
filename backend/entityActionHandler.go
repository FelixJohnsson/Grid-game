package main

import "fmt"

type TaskType string

const (
	FindWater  TaskType = "Find water supply" // FindAndNoteWaterSupply
	DrinkWater TaskType = "Drink water" // DrinKWaterTask

	EatFood    TaskType = "Eat food"
	ClearAirway TaskType = "Clear airway" // ClearAirway > BodyFunction
	FixNose    TaskType = "Fix nose"	// FixNose > BodyFunction
	ReducePain TaskType = "Reduce pain"
	FindFood   TaskType = "Find food" // FindAndNoteFoodSupply
	HaveFood   TaskType = "Have food for storage"
	FindShelter TaskType = "Find shelter"
	MakeShelter TaskType = "Make shelter"
	ImproveDefense TaskType = "Improve defense"
	Talk TaskType = "Talk"
)

func (b *Brain) CheckIfCurrentMotorTaskIsDone(MotorCortexAction MotorCortexAction, ActionReason string) bool {
	return  b.MotorCortexCurrentTask.ActionReason == ActionReason && b.MotorCortexCurrentTask.Finished
}

func (b *Brain) ActionHandler() {
	// Take the action with the highest priority
	action := b.RankTasks()

	fmt.Println(Blue + "Action: ", action.Action)
	fmt.Println(Reset)

	// Perform the action
	switch action.Action {
	case "Find a water supply":
		b.CurrentTask = action
		b.FindAndNoteWaterSupply()
		return
	case "Drink water":
		b.CurrentTask = action
		b.DrinkWaterTask(action)
	case "Eat food":
		b.CurrentTask = action
		b.EatFoodTask(action)
	case "Clear airway":
		b.CurrentTask = action
		b.Owner.ClearAirway(action)
		return
	case "Fix nose":
		b.CurrentTask = action
		b.Owner.FixBrokenNose(action)
		return
	case "Reduce pain":
		b.CurrentTask = action
		//b.ReducePain()
		return

	case "Find a food supply":
		b.CurrentTask = action
		success := b.FindAndNoteFoodSupply()

		if success {
			b.RemoveActionFromActionList(action)
		} else {
			b.CurrentTask = TargetedAction{"Find a food supply", "", true, []BodyPartType{"Legs"}, 90}
			b.GoSearchFor("Food supply")
		}
		return
	case "Have food for storage":
		b.CurrentTask = action
		b.GetFoodForStorage(action)
	case "Find shelter":
		b.CurrentTask = action

		return
	case "Make shelter":
		fmt.Println(b.Owner.FullName + " is making a shelter.")
		b.CurrentTask = action
		return
	case "Improve defense":
		fmt.Println(b.Owner.FullName + " is improving defense.")
		b.CurrentTask = action
	case "Idle":
		b.CurrentTask = action
		fmt.Println(b.Owner.FullName + " is idle.")
		return

	}
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
	if !b.PhysiologicalNeeds.WayOfGettingWater {
		action := TargetedAction{"Find a water supply", "", false,[]BodyPartType{"Hands"}, 100}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
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


// ----------------- Tash requests ------------

// Receive a requested task from another person
func (b *Brain) ReceiveTaskRequest(requestedTask RequestedAction) bool {
    fmt.Println(b.Owner.FullName + " received a task request from " + requestedTask.From.FullName)
    
    hasRelationship := b.Owner.HasRelationship(requestedTask.From.FullName)

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
func (b *Brain) SendTaskRequest(to *Entity, taskType TaskType) {
    if b.Owner.IsTalking.IsActive {
        fmt.Println(b.Owner.FullName + " is already talking to someone.")
        return 
    }
    fmt.Println(b.Owner.FullName + " is sending a task request to " + to.FullName)
    task := RequestedAction{TargetedAction{taskType, to.FullName, true, make([]BodyPartType, 0), 10}, b.Owner}
    success := to.Brain.ReceiveTaskRequest(task)
    if success {
        fmt.Println(to.FullName + " accepted the task request.")
        //b.Owner.IsTalking = TargetedAction{"Hello " + to.FullName + ", how are you doing?", to.FullName, true, make([]BodyPartType, 0), 10}
        fmt.Println(b.Owner.FullName + " is talking to " + to.FullName)
    } else {
        fmt.Println(to.FullName + " declined the task request.")
    }
}