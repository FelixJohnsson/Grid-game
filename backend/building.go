package main

type BuildingType string

const (
	House     BuildingType = "House"
	WoodCabin BuildingType = "Lumberjack"
	Mine      BuildingType = "Mine"
	Farm      BuildingType = "Farm"
	Barracks  BuildingType = "Barracks"
)

// Enum for BuildingState
type BuildingState string

const (
	Idle             BuildingState = "Idle"
	Working          BuildingState = "Working"
	LackingResources BuildingState = "Lacking Resources"
	LackingWorkers   BuildingState = "Lacking Workers"
)

// Enum for BuildingDescription
type BuildingDescription string

const (
	HouseDescription      BuildingDescription = "The House provides your inhabitants with shelter."
	LumberjackDescription BuildingDescription = "The Lumberjack gathers wood for your village."
	MineDescription       BuildingDescription = "The Mine extracts valuable minerals."
	FarmDescription       BuildingDescription = "The Farm produces food for your villagers."
	BarracksDescription   BuildingDescription = "The Barracks trains your military units."
)

// Enum for BuildingIcon
type BuildingIcon string

const (
	HouseIcon      BuildingIcon = "H"
	LumberjackIcon BuildingIcon = "L"
	MineIcon       BuildingIcon = "M"
	FarmIcon       BuildingIcon = "F"
	BarracksIcon   BuildingIcon = "B"
)

// Enum for BuildingColor
type BuildingColor string

const (
	HouseColor      BuildingColor = "bg-blue-500"
	LumberjackColor BuildingColor = "bg-orange-900"
	MineColor       BuildingColor = "bg-gray-500"
	FarmColor       BuildingColor = "bg-yellow-500"
	BarracksColor   BuildingColor = "bg-red-500"
)

// Struct for Location (assuming it's the same as previously defined)
type Location struct {
	X int
	Y int
}

// Struct for Resources
type Resources struct {
	Inhabitants []Person
	Wood        int
	Food        int
	Stone       int
	Ores        int
	Money       int
}

// Struct for Building
type Building struct {
	Type        BuildingType
	Name        string
	State       BuildingState
	Description string
	Icon        string
	Color       string
	Location    Location
	Workers     []Person
	Resources   Resources
}

// NewBuilding creates a new building
func NewBuilding(t BuildingType, n string, l Location) Building {
	var b Building
	b.Type = t
	b.Name = n
	b.State = Idle
	b.Location = l

	switch t {
	case House:
		b.Description = string(HouseDescription)
		b.Icon = string(HouseIcon)
		b.Color = string(HouseColor)
	case WoodCabin:
		b.Description = string(LumberjackDescription)
		b.Icon = string(LumberjackIcon)
		b.Color = string(LumberjackColor)
	case Mine:
		b.Description = string(MineDescription)
		b.Icon = string(MineIcon)
		b.Color = string(MineColor)
	case Farm:
		b.Description = string(FarmDescription)
		b.Icon = string(FarmIcon)
		b.Color = string(FarmColor)
	case Barracks:
		b.Description = string(BarracksDescription)
		b.Icon = string(BarracksIcon)
		b.Color = string(BarracksColor)
	}

	return b
}

// AddWorker adds a person to the building's workers
func (b *Building) AddWorker(p Person) {
	b.Workers = append(b.Workers, p)
}

// RemoveWorker removes a person from the building's workers
func (b *Building) RemoveWorker(p Person) {
	for i, worker := range b.Workers {
		if worker.FullName == p.FullName {
			b.Workers = append(b.Workers[:i], b.Workers[i+1:]...)
			break
		}
	}
}

// UpdateState updates the building's state based on its resources and workers
func (b *Building) UpdateState() {
	if len(b.Workers) == 0 {
		b.State = LackingWorkers
	} else {
		b.State = Working
	}
}

// createNewBuilding creates a new building and adds it to the array
func createNewBuilding(t BuildingType, n string, l Location) Building {
	b := NewBuilding(t, n, l)
	addBuilding(b)
	return b
}