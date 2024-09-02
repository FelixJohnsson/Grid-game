package main

// All available stone age materials
var materials = []Material{
	{"Grass", "Organic", 1, 1, 1, 10},
	{"Wood", "Organic", 2, 1, 1, 3},
	{"Stone", "Inorganic", 5, 3, 4, 2},
	{"Leather", "Organic", 1, 1, 1, 7},
	{"Bone", "Organic", 3, 2, 2, 4},
	{"Flint", "Inorganic", 6, 4, 5, 3},
	{"Obsidian", "Inorganic", 7, 5, 6, 3},
	{"Feathers", "Organic", 1, 1, 1, 6},
	{"Clay", "Inorganic", 4, 2, 3, 4},
	{"Fur", "Organic", 1, 1, 1, 5},
	{"Wet Clay", "Inorganic", 4, 2, 3, 10},
	{"Water", "Inorganic", 1, 1, 1, 1},
	{"Soil", "Inorganic", 1, 1, 1, 1},
}

// GetMaterialByName - Get a material by name
func GetMaterialByName(name string) Material {
	for _, material := range materials {
		if material.Name == name {
			return material
		}
	}
	return Material{}
}

// All available items
var items = []Item{
	// Weapons
	{"Wooden Spear", 5, 4, 3, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Staff", 1, 8, 2, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},

	// Tools
	{"Stone Axe", 6, 2, 5, []Material{materials[1], GetMaterialByName("Wood"), GetMaterialByName("Stone")}, make([]Residue, 0), Location{0, 0}},

	// Storage
	{"Food Box", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Box", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Crate", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Woven Grass Basket", 1, 1, 1, []Material{GetMaterialByName("Grass")}, make([]Residue, 0), Location{0, 0}},
}

var BuildingMaterials = []Item{
	{"Wood log", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
}

// CreateNewItem - Create a new item
func CreateNewItem(itemType string) *Item {
	switch itemType {
	// Items
	case "Wooden Spear":
		return &items[0]
	case "Wooden Staff":
		return &items[1]
	case "Stone Axe":
		return &items[2]
	case "Food Box":
		return &items[3]
	case "Wooden Box":
		return &items[4]
	case "Wooden Crate":
		return &items[5]
	case "Woven Grass Basket":
		return &items[6]

	// Building materials
	case "Wood log":
		return &BuildingMaterials[0]
	}

	return nil
}