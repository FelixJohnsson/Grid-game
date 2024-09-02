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

type Building struct {
	Type      BuildingType
	Name      string
	State     BuildingState
	Location  Location
	Workers   []*Person
	Inventory []*Item
}

type Shelter struct {
	Owner       *Person
	Location    Location
	Inhabitants []*Person
	Inventory   []*Item
}

func NewBuilding(x, y int, buildingType BuildingType) *Building {
	building := &Building{
		Type:      buildingType,
		Name:      string(buildingType),
		State:     Idle,
		Location:  Location{X: x, Y: y},
		Workers:   []*Person{},
		Inventory: []*Item{},
	}
	return building
}

func NewShelter(x, y int, owner *Person) *Shelter {
	shelter := &Shelter{
		Owner:       owner,
		Location:    Location{X: x, Y: y},
		Inhabitants: []*Person{},
		Inventory:   []*Item{},
	}
	return shelter
}
