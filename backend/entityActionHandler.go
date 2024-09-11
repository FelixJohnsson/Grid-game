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

func (b *Brain) DrinkWaterTask(TargetedAction TargetedAction) {
	if b.CheckIfCurrentMotorTaskIsDone(b.MotorCortexCurrentTask, "Drink water") {
		water := Liquid{"Water"}
		b.Owner.Drink(water)
		return
	}

	success := b.FindAndNoteWaterSupply()

	if success {
		water := b.GetWaterInVision()
		closestWater := b.FindClosestWaterSupply(water)
		b.MotorCortexCurrentTask = MotorCortexAction{"Drink water", "Walk", Location{closestWater.Location.X, closestWater.Location.Y}, false, false}
	} else {

	}
}

func (b *Brain) EatFoodTask(TargetedAction TargetedAction) {
	
}


func (b *Brain) ActionHandler() {
	// Take the action with the highest priority
	action := b.RankTasks()

	fmt.Println("Action: ", action.Action)

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