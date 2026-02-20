package main

type TaskType string

const (
	FindWater  TaskType = "Find water supply"
	DrinkWater TaskType = "Drink water"
	HaveWater  TaskType = "Get water for storage"

	FindFood     TaskType = "Find food"
	EatFood      TaskType = "Eat food"
	HaveFood     TaskType = "Get food for storage"
	CraftFoodBox TaskType = "Craft food box"

	FindLumber TaskType = "Find lumber tree"
	FindSticks TaskType = "Find sticks"
	HaveLumber TaskType = "Get lumber for storage"
	ChopTree   TaskType = "Chop down tree"

	FindStone  TaskType = "Find stone"
	HaveStone  TaskType = "Get stone for storage"
	CraftStone TaskType = "Get stone"

	FindGrass TaskType = "Find grass"
	HaveGrass TaskType = "Get grass for storage"
	CutGrass  TaskType = "Cut down grass"

	CraftItem TaskType = "Craft item"

	ClearAirway TaskType = "Clear airway"
	FixNose     TaskType = "Fix nose"
	ReducePain  TaskType = "Reduce pain"

	GetWarm      TaskType = "Get warm"
	Excrete      TaskType = "Excrete"
	FindSafeArea TaskType = "Find a safe area"
	Rest         TaskType = "Rest"

	FindShelter        TaskType = "Find shelter"
	ClaimBase          TaskType = "Claim base"
	MakeShelter        TaskType = "Make shelter"
	StockpileResources TaskType = "Stockpile resources"
	BuildStorage       TaskType = "Build storage"

	Talk TaskType = "Talk"
	None TaskType = "Idle"
)

const (
	wantBreath             = "Be able to breath"
	wantRelievePain        = "Relieve pain"
	wantConsumeWater       = "Consume water"
	wantConsumeFood        = "Consume food"
	wantGetWarm            = "Get warm"
	wantExcrete            = "Excrete"
	wantSafeArea           = "Find a safe area"
	wantStockpileResources = "Stockpile resources"
	wantMakeShelter        = "Make shelter"
	wantRest               = "Rest"
)

type actionExecutor func(*Brain, TargetedAction)

type wantRule struct {
	want      string
	condition func(*Brain) bool
	onMatch   func(*Brain)
}

type planningRule struct {
	condition func(*Brain) bool
	build     func(*Brain) []TargetedAction
}

func handRequiredLimbs() []BodyPartType {
	return []BodyPartType{"Hands"}
}

func idleTask() TargetedAction {
	return TargetedAction{
		Action:       None,
		Target:       "Nothing",
		IsActive:     false,
		RequiresLimb: handRequiredLimbs(),
		Priority:     0,
	}
}

func newTask(action TaskType, target string, priority int) TargetedAction {
	return TargetedAction{
		Action:       action,
		Target:       target,
		IsActive:     false,
		RequiresLimb: handRequiredLimbs(),
		Priority:     priority,
	}
}

func noopAction(_ *Brain, _ TargetedAction) {}

var taskExecutors = map[TaskType]actionExecutor{
	FindWater: func(b *Brain, _ TargetedAction) {
		b.FindWaterSupply()
	},
	DrinkWater: func(b *Brain, action TargetedAction) {
		b.DrinkWaterTask(action)
	},
	HaveWater: noopAction,

	FindFood: func(b *Brain, _ TargetedAction) {
		b.FindFoodSupply()
	},
	EatFood: func(b *Brain, _ TargetedAction) {
		b.EatFoodTask()
	},
	HaveFood: func(b *Brain, action TargetedAction) {
		b.GetFoodForStorage(action)
	},
	CraftFoodBox: func(b *Brain, action TargetedAction) {
		b.CraftItem(action)
	},

	FindLumber: func(b *Brain, _ TargetedAction) {
		b.GetLumberTask()
	},
	HaveLumber: func(b *Brain, action TargetedAction) {
		b.MaintainStickStockTask(action)
	},
	ChopTree:   noopAction,
	FindSticks: noopAction,

	FindStone: func(b *Brain, _ TargetedAction) {
		b.FindStoneSupply()
	},
	HaveStone: func(b *Brain, action TargetedAction) {
		b.MaintainStoneStockTask(action)
	},
	CraftStone: func(b *Brain, action TargetedAction) {
		b.CraftItem(action)
	},

	FindGrass: func(b *Brain, _ TargetedAction) {
		b.FindGrassSupply()
	},
	HaveGrass: func(b *Brain, action TargetedAction) {
		b.MaintainGrassStockTask(action)
	},
	CutGrass: noopAction,

	CraftItem: func(b *Brain, action TargetedAction) {
		b.CraftItem(action)
	},

	ClearAirway: func(b *Brain, action TargetedAction) {
		b.Owner.ClearAirway(action)
	},
	FixNose: func(b *Brain, action TargetedAction) {
		b.Owner.FixBrokenNose(action)
	},
	ReducePain: noopAction,

	GetWarm:      noopAction,
	Excrete:      noopAction,
	FindSafeArea: noopAction,
	Rest:         noopAction,

	FindShelter: noopAction,
	ClaimBase: func(b *Brain, action TargetedAction) {
		b.ClaimBaseTask(action)
	},
	MakeShelter: func(b *Brain, action TargetedAction) {
		b.MakeShelterTask(action)
	},
	StockpileResources: func(b *Brain, action TargetedAction) {
		b.StockpileResourcesTask(action)
	},
	BuildStorage: noopAction,

	Talk: noopAction,
	None: noopAction,
}

func (b *Brain) ActionHandler() {
	action := b.RankTasks()
	b.CurrentTask = action

	if action.Action == None {
		b.CurrentTask.IsActive = false
		return
	}

	executor, ok := taskExecutors[action.Action]
	if !ok {
		b.CurrentTask = idleTask()
		return
	}

	b.CurrentTask.IsActive = true
	executor(b, action)
}

func (b *Brain) CheckIfCurrentMotorTaskIsDone(MotorCortexAction MotorCortexAction, ActionReason string) bool {
	return b.MotorCortexCurrentTask.ActionReason == ActionReason && b.MotorCortexCurrentTask.Finished
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

var homoSapiensWantRules = []wantRule{
	{
		want: wantBreath,
		condition: func(b *Brain) bool {
			return !b.CheckIfCanBreath()
		},
	},
	{
		want: wantRelievePain,
		condition: func(b *Brain) bool {
			return b.PhysiologicalNeeds.IsInPain
		},
	},
	{
		want: wantConsumeWater,
		condition: func(b *Brain) bool {
			return b.Owner.Body.Blood.Water < 40
		},
	},
	{
		want: wantConsumeFood,
		condition: func(b *Brain) bool {
			return b.Owner.Body.Blood.Glucose < 40
		},
	},
	{
		want: wantGetWarm,
		condition: func(b *Brain) bool {
			return !b.PhysiologicalNeeds.IsSufficientlyWarm
		},
	},
	{
		want: wantExcrete,
		condition: func(b *Brain) bool {
			return b.PhysiologicalNeeds.NeedToExcrete
		},
	},
	{
		want: wantSafeArea,
		condition: func(b *Brain) bool {
			return !b.PhysiologicalNeeds.IsInSafeArea
		},
	},
	{
		want: wantStockpileResources,
		condition: func(b *Brain) bool {
			if !b.Base.Claimed || !b.PhysiologicalNeeds.HasShelter {
				return false
			}
			b.ensureBaseStockpile()
			return b.Base.Stockpile[stickItemName] < baseStickMaintenanceTarget ||
				b.Base.Stockpile[stoneStockpileKey] < baseStoneStockTarget ||
				b.Base.Stockpile[grassStockpileKey] < baseGrassStockTarget
		},
	},
	{
		want: wantMakeShelter,
		condition: func(b *Brain) bool {
			return !b.PhysiologicalNeeds.HasShelter
		},
	},
	{
		want: wantRest,
		condition: func(b *Brain) bool {
			return b.PhysiologicalNeeds.Rested < 20
		},
	},
}

// HomoSapiensCalculateWant - Calculate the want of the person
func (b *Brain) HomoSapiensCalculateWant() {
	// Keep behavior consistent with previous switch: choose the first matching want.
	for _, rule := range homoSapiensWantRules {
		if !rule.condition(b) {
			continue
		}
		if rule.onMatch != nil {
			rule.onMatch(b)
		}
		if !b.CheckIfWantIsAlreadyInList(rule.want) {
			b.Owner.WantsTo = append(b.Owner.WantsTo, rule.want)
		}
		return
	}
}

func (b *Brain) buildBreathingTasks() []TargetedAction {
	tasks := make([]TargetedAction, 0, 3)
	if b.Owner.Body.Head.Mouth.IsObstructed {
		tasks = append(tasks, newTask(ClearAirway, "Mouth", 100))
	}
	if b.Owner.Body.Head.Nose.IsObstructed {
		tasks = append(tasks, newTask(ClearAirway, "Nose", 100))
	}
	if b.Owner.Body.Head.Nose.IsBroken {
		tasks = append(tasks, newTask(FixNose, "Nose", 100))
	}
	return tasks
}

func (b *Brain) buildFoodSupplyTasks() []TargetedAction {
	if !b.Base.Claimed {
		return nil
	}
	b.ensureBaseStockpile()
	if b.Base.Stockpile[foodStockpileKey] < baseFoodStockTarget {
		return []TargetedAction{newTask(HaveFood, "", 55)}
	}
	return nil
}

func (b *Brain) getPlanningRules() []planningRule {
	return []planningRule{
		{
			condition: func(b *Brain) bool { return !b.CheckIfCanBreath() },
			build: func(b *Brain) []TargetedAction {
				return b.buildBreathingTasks()
			},
		},
		{
			condition: func(b *Brain) bool { return !b.PhysiologicalNeeds.WayOfGettingWater },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(FindWater, "", 100)} },
		},
		{
			condition: func(b *Brain) bool { return b.Owner.Body.Blood.Water < 40 },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(DrinkWater, "", 99)} },
		},
		{
			condition: func(b *Brain) bool { return b.Owner.Body.Blood.Glucose < 40 },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(EatFood, "", 98)} },
		},
		{
			condition: func(b *Brain) bool { return b.PhysiologicalNeeds.IsInPain },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(ReducePain, "", 95)} },
		},
		{
			condition: func(b *Brain) bool { return !b.PhysiologicalNeeds.WayOfGettingFood },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(FindFood, "", 90)} },
		},
		{
			condition: func(b *Brain) bool { return !b.PhysiologicalNeeds.FoodSupply },
			build:     func(b *Brain) []TargetedAction { return b.buildFoodSupplyTasks() },
		},
		{
			condition: func(b *Brain) bool { return !b.PhysiologicalNeeds.IsSufficientlyWarm },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(GetWarm, "", 85)} },
		},
		{
			condition: func(b *Brain) bool { return b.PhysiologicalNeeds.NeedToExcrete },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(Excrete, "", 80)} },
		},
		{
			condition: func(b *Brain) bool { return !b.PhysiologicalNeeds.IsInSafeArea },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(FindSafeArea, "", 75)} },
		},
		{
			condition: func(b *Brain) bool { return !b.PhysiologicalNeeds.HasShelter && !b.Base.Claimed },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(ClaimBase, "", 92)} },
		},
		{
			condition: func(b *Brain) bool { return !b.PhysiologicalNeeds.HasShelter },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(MakeShelter, "", 91)} },
		},
		{
			condition: func(b *Brain) bool {
				return b.PhysiologicalNeeds.HasShelter &&
					b.Base.Claimed &&
					b.Base.Stockpile[stickItemName] < baseStickMaintenanceTarget &&
					b.Owner.Body.Blood.Water >= 40 &&
					b.Owner.Body.Blood.Glucose >= 40
			},
			build: func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(HaveLumber, "", 20)} },
		},
		{
			condition: func(b *Brain) bool {
				return b.PhysiologicalNeeds.HasShelter &&
					b.Base.Claimed &&
					b.Base.Stockpile[stoneStockpileKey] < baseStoneStockTarget &&
					b.Owner.Body.Blood.Water >= 40 &&
					b.Owner.Body.Blood.Glucose >= 40
			},
			build: func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(HaveStone, "", 19)} },
		},
		{
			condition: func(b *Brain) bool {
				return b.PhysiologicalNeeds.HasShelter &&
					b.Base.Claimed &&
					b.Base.Stockpile[grassStockpileKey] < baseGrassStockTarget &&
					b.Owner.Body.Blood.Water >= 40 &&
					b.Owner.Body.Blood.Glucose >= 40
			},
			build: func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(HaveGrass, "", 18)} },
		},
		{
			condition: func(b *Brain) bool {
				b.ensureBaseStockpile()
				recipe, ok := GetCraftingRecipe("Stone Axe")
				return ok &&
					b.PhysiologicalNeeds.HasShelter &&
					b.Base.Claimed &&
					!b.hasBaseOrOwnedItem("Stone Axe") &&
					b.hasStockpileIngredients(recipe.Ingredients)
			},
			build: func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(CraftStone, "Stone Axe", 18)} },
		},
		{
			condition: func(b *Brain) bool {
				b.ensureBaseStockpile()
				recipe, ok := GetCraftingRecipe("Food Box")
				return ok &&
					b.PhysiologicalNeeds.HasShelter &&
					b.Base.Claimed &&
					!b.hasBaseOrOwnedItem("Food Box") &&
					b.hasStockpileIngredients(recipe.Ingredients)
			},
			build: func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(CraftFoodBox, "Food Box", 17)} },
		},
		{
			condition: func(b *Brain) bool { return b.PhysiologicalNeeds.Rested < 20 },
			build:     func(_ *Brain) []TargetedAction { return []TargetedAction{newTask(Rest, "", 65)} },
		},
		{
			condition: func(b *Brain) bool {
				if !b.PhysiologicalNeeds.HasShelter || !b.Base.Claimed {
					return false
				}
				if b.Owner.Body.Blood.Water < 40 || b.Owner.Body.Blood.Glucose < 40 {
					return false
				}
				b.ensureBaseStockpile()
				return b.Base.Stockpile[stickItemName] < baseStickMaintenanceTarget ||
					b.Base.Stockpile[stoneStockpileKey] < baseStoneStockTarget ||
					b.Base.Stockpile[grassStockpileKey] < baseGrassStockTarget
			},
			build: func(_ *Brain) []TargetedAction {
				return []TargetedAction{newTask(StockpileResources, "", 60)}
			},
		},
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

// RemoveActionFromActionList - Remove an action from the action list
func (b *Brain) RemoveActionFromActionList(action TargetedAction) {
	for i, a := range b.ActionList {
		if a.Action == action.Action && a.Target == action.Target {
			b.ActionList = append(b.ActionList[:i], b.ActionList[i+1:]...)
			return
		}
	}
}

// ClearCurrentTask - Clear the current task
func (b *Brain) ClearCurrentTask() {
	b.CurrentTask = idleTask()
}

// AddTaskToActionList - Add a task to the action list
func (b *Brain) AddTaskToActionList(task TargetedAction) {
	if task.Action == None {
		return
	}
	if b.IsTaskInActionList(task) {
		return
	}

	filtered := make([]TargetedAction, 0, len(b.ActionList))
	for _, action := range b.ActionList {
		if action.Action != None {
			filtered = append(filtered, action)
		}
	}
	filtered = append(filtered, task)
	b.ActionList = filtered
}

// ClearActionList - Clear the action list
func (b *Brain) ClearActionList() {
	b.ActionList = make([]TargetedAction, 0)
}

// TranslateWantToTaskList - Translate current physiological state to a task list.
func (b *Brain) TranslateWantToTaskList() {
	b.ClearActionList()

	for _, rule := range b.getPlanningRules() {
		if !rule.condition(b) {
			continue
		}
		for _, task := range rule.build(b) {
			b.AddTaskToActionList(task)
		}
	}

	if len(b.ActionList) == 0 {
		b.ActionList = append(b.ActionList, idleTask())
	}
}

// RankTasks - Rank tasks based on priority
func (b *Brain) RankTasks() TargetedAction {
	if len(b.ActionList) == 0 {
		return idleTask()
	}

	highestPriorityAction := idleTask()
	highestPriority := highestPriorityAction.Priority

	for _, action := range b.ActionList {
		if action.Priority > highestPriority {
			highestPriority = action.Priority
			highestPriorityAction = action
		}
	}

	return highestPriorityAction
}

// ----------------- Task requests ------------

// ReceiveTaskRequest - Receive a requested task from another person
func (b *Brain) ReceiveTaskRequest(requestedTask RequestedAction) bool {
	hasRelationship := b.Owner.HasRelationship(requestedTask.From.FullName)

	// For now we only support talk requests.
	if !hasRelationship {
		return false
	}

	if requestedTask.Action == Talk && b.Owner.IsTalking.IsActive {
		return false
	}
	if requestedTask.Action == Talk && !b.Owner.IsTalking.IsActive {
		b.Owner.IsTalking = TargetedAction{
			Action:       Talk,
			Target:       requestedTask.From.FullName,
			IsActive:     true,
			RequiresLimb: make([]BodyPartType, 0),
			Priority:     10,
		}
		return true
	}

	return false
}

// SendTaskRequest - Send a task request to another person
func (b *Brain) SendTaskRequest(to *Entity, taskType TaskType) {
	if b.Owner.IsTalking.IsActive {
		return
	}
	task := RequestedAction{
		TargetedAction: TargetedAction{
			Action:       taskType,
			Target:       to.FullName,
			IsActive:     true,
			RequiresLimb: make([]BodyPartType, 0),
			Priority:     10,
		},
		From: b.Owner,
	}
	_ = to.Brain.ReceiveTaskRequest(task)
}
