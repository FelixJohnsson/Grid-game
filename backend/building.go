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
	Workers   []*Entity
	Inventory []*Item
}

type Shelter struct {
	OwnerName   string
	Location    Location
	Inhabitants []*Entity
	Inventory   []*Item
}

func NewShelter(x, y int, owner *Entity) *Shelter {
	shelter := &Shelter{
		OwnerName:   owner.FullName,
		Location:    Location{X: x, Y: y},
		Inhabitants: []*Entity{},
		Inventory:   []*Item{},
	}
	return shelter
}
