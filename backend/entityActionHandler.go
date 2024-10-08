package main

import "fmt"

type TaskType string

const (
	FindWater  TaskType = "Find water supply"
	DrinkWater TaskType = "Drink water"
	HaveWater  TaskType = "Get water for storage"

	FindFood   TaskType = "Find food"
	EatFood    TaskType = "Eat food"	
	HaveFood   TaskType = "Get food for storage"
	Hunt       TaskType = "Hunt"

	FindLumber TaskType = "Find lumber tree"
	HaveLumber TaskType = "Get lumber for storage"
	ChopTree   TaskType = "Chop down tree"

	FindStone TaskType = "Find stone"
	HaveStone TaskType = "Get stone for storage"
	CraftStone TaskType = "Get stone"

	FindGrass TaskType = "Find grass"
	HaveGrass TaskType = "Get grass for storage"
	CutGrass TaskType = "Cut down grass"

	CraftItem TaskType = "Craft item"

	ClearAirway TaskType = "Clear airway"
	FixNose    TaskType = "Fix nose"
	ReducePain TaskType = "Reduce pain"

	FindShelter TaskType = "Find shelter"
	MakeShelter TaskType = "Make shelter"
	ImproveDefense TaskType = "Improve defense"

	Talk TaskType = "Talk"
	None TaskType = "Idle"
)

func (b *Brain) ActionHandler() {
	// Take the action with the highest priority
	action := b.RankTasks()

	fmt.Println(Blue + string(b.Owner.Species) + "- Action:", action.Action)
	fmt.Println(Reset)

	// Perform the action
	switch action.Action {
	// ----------------- Water ---------------
	case FindWater:
		b.CurrentTask = action
		b.FindWaterSupply()
		return
	case DrinkWater:
		b.CurrentTask = action
		b.DrinkWaterTask(action)
		return
	case HaveWater:
		b.CurrentTask = action
		//b.GetWaterForStorage(action)
		return
	
	// ----------------- Food -----------------
	case FindFood:
		b.CurrentTask = action
		if b.Owner.Predator && !b.Owner.Herbivore{
			fmt.Println(Red + "I am a predator, so I will hunt food." + Reset)
			b.PredatorFindFood()
		} else {
			fmt.Println(Green + "I will search for vegetables or fruits." + Reset)
			b.FindFoodSupply()
		}
		return
	case EatFood:
		b.CurrentTask = action
		if b.Owner.Predator && !b.Owner.Herbivore{
			fmt.Println(Red + "I am a predator, so I will hunt food." + Reset)
			b.PredatorFindFood()
		} else {
			fmt.Println(Green + "I will search for vegetables or fruits." + Reset)
			b.EatFoodTask()
		}
	case HaveFood:
		b.CurrentTask = action
		b.GetFoodForStorage(action)
		return

	// ----------------- Lumber ---------------
	case FindLumber:
		b.GetLumberTask()
		b.CurrentTask = action
		return
	case HaveLumber:
		b.CurrentTask = action
		//b.GetLumberForStorage(action)
		return
	case ChopTree:
		b.CurrentTask = action
		//b.ChopTree(action)
		return

	// ----------------- Stone ---------------
	case FindStone:
		b.CurrentTask = action
		return
	case HaveStone:
		b.CurrentTask = action
		//b.GetStoneForStorage(action)
		return
	case CraftStone:
		b.CurrentTask = action
		//b.CraftStone(action)
		return

	// ----------------- Grass ---------------
	case FindGrass:
		b.CurrentTask = action
		return
	case HaveGrass:
		b.CurrentTask = action
		//b.GetGrassForStorage(action)
		return
	case CutGrass:
		b.CurrentTask = action
		//b.CutGrass(action)
		return

	// ----------------- Craft ---------------
	case CraftItem:
		b.CurrentTask = action
		b.CraftItem(action)
		return
	

	// ----------------- Physical -------------
	case ClearAirway:
		b.CurrentTask = action
		b.Owner.ClearAirway(action)
		return
	case FixNose:
		b.CurrentTask = action
		b.Owner.FixBrokenNose(action)
		return
	case ReducePain:
		b.CurrentTask = action
		//b.ReducePain(action)
		return

	// ----------------- Shelter --------------
	case FindShelter:
		b.CurrentTask = action
		//b.FindShelter(action)
		return
	case MakeShelter:
		b.CurrentTask = action
		return

	// ----------------- Social ---------------
	case Talk:
		b.CurrentTask = action
		//b.Talk(action)
		return

	// ----------------- Improvements ---------
	case ImproveDefense:
		b.CurrentTask = action

	
	// ----------------- Misc -----------------
	case None:
		b.CurrentTask = action
		return
	}
}

func (b *Brain) CheckIfCurrentMotorTaskIsDone(MotorCortexAction MotorCortexAction, ActionReason string) bool {
	return  b.MotorCortexCurrentTask.ActionReason == ActionReason && b.MotorCortexCurrentTask.Finished
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
            action := TargetedAction{ClearAirway, "Mouth", false,[]BodyPartType{"Hands"}, 100}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
        }
        if b.Owner.Body.Head.Nose.IsObstructed {
            action := TargetedAction{ClearAirway, "Nose", false,[]BodyPartType{"Hands"}, 100}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
        }
        if b.Owner.Body.Head.Nose.IsBroken {
            action := TargetedAction{FixNose, "Nose", false,[]BodyPartType{"Hands"}, 100}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
        }
    }
	if !b.PhysiologicalNeeds.WayOfGettingWater {
		action := TargetedAction{FindWater, "", false,[]BodyPartType{"Hands"}, 100}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
    if b.PhysiologicalNeeds.Thirst > 30 {
		action := TargetedAction{DrinkWater, "", false,[]BodyPartType{"Hands"}, 99}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if b.PhysiologicalNeeds.Hunger > 30 {
		action := TargetedAction{EatFood, "", false,[]BodyPartType{"Hands"}, 98}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if b.PhysiologicalNeeds.IsInPain {
		action := TargetedAction{ReducePain, "", false,[]BodyPartType{"Hands"}, 95}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}	
	}

	if !b.PhysiologicalNeeds.WayOfGettingFood {
		action := TargetedAction{FindFood, "", false, []BodyPartType{"Hands"}, 90}
		if !b.IsTaskInActionList(action) {
			b.AddTaskToActionList(action)
		}
	}
	if !b.PhysiologicalNeeds.FoodSupply {
		if b.FindInOwnedItems("Food Box") == nil && b.FindInOwnedItems("Wood log") == nil {
			action := TargetedAction{FindLumber, "", false, []BodyPartType{"Hands"}, 90}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
		} else if b.FindInOwnedItems("Food Box") == nil && b.FindInOwnedItems("Wood log") != nil {
			action := TargetedAction{"Craft food box", "", false, []BodyPartType{"Hands"}, 90}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
		} else if b.FindInOwnedItems("Food Box") != nil {
			action := TargetedAction{"Get food for storage", "", false, []BodyPartType{"Hands"}, 90}
			if !b.IsTaskInActionList(action) {
				b.AddTaskToActionList(action)
			}
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
    
    hasRelationship := b.Owner.HasRelationship(requestedTask.From.FullName)

    // For now we will just accept the task
    if hasRelationship {
        if requestedTask.Action == "Talk" && b.Owner.IsTalking.IsActive {
            return false
        } else if requestedTask.Action == "Talk" && !b.Owner.IsTalking.IsActive {
            b.Owner.IsTalking = TargetedAction{"Bla bla bla ...", requestedTask.From.FullName, true, make([]BodyPartType, 0), 10}
            return true
        }
    } else {
        return false
    }
    return false
}

// Send a task request to another person
func (b *Brain) SendTaskRequest(to *Entity, taskType TaskType) {
    if b.Owner.IsTalking.IsActive {
        return 
    }
    task := RequestedAction{TargetedAction{taskType, to.FullName, true, make([]BodyPartType, 0), 10}, b.Owner}
    success := to.Brain.ReceiveTaskRequest(task)
    if success {
        //b.Owner.IsTalking = TargetedAction{"Hello " + to.FullName + ", how are you doing?", to.FullName, true, make([]BodyPartType, 0), 10}
    } else {
    }
}