package main

// All available stone age materials
var materials = []Material{
	{"Grass", "Organic", 1, 1, 1, 10},
	{"Wood", "Organic", 2, 1, 1, 3},
	{"Sticks", "Organic", 1, 1, 1, 1},
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
	{"Stone Spear", 7, 3, 3, []Material{GetMaterialByName("Stone"), GetMaterialByName("Sticks")}, make([]Residue, 0), Location{0, 0}},

	// Tools
	{"Stone Axe", 6, 2, 5, []Material{materials[1], GetMaterialByName("Wood"), GetMaterialByName("Stone")}, make([]Residue, 0), Location{0, 0}},
	{"Stone Knife", 8, 1, 1, []Material{GetMaterialByName("Stone"), GetMaterialByName("Sticks")}, make([]Residue, 0), Location{0, 0}},
	{"Stone Hammer", 1, 7, 4, []Material{GetMaterialByName("Stone"), GetMaterialByName("Sticks")}, make([]Residue, 0), Location{0, 0}},
	{"Stone Pickaxe", 4, 6, 5, []Material{GetMaterialByName("Stone"), GetMaterialByName("Sticks")}, make([]Residue, 0), Location{0, 0}},

	// Storage
	{"Food Box", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Box", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Crate", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Woven Grass Basket", 1, 1, 1, []Material{GetMaterialByName("Grass")}, make([]Residue, 0), Location{0, 0}},
}

var BuildingMaterials = []Item{
	{"Wood log", 1, 1, 1, []Material{GetMaterialByName("Wood")}, make([]Residue, 0), Location{0, 0}},
	{"Sticks", 1, 1, 1, []Material{GetMaterialByName("Sticks")}, make([]Residue, 0), Location{0, 0}},
	{"Stone", 0, 1, 2, []Material{GetMaterialByName("Stone")}, make([]Residue, 0), Location{0, 0}},
	{"Grass", 0, 0, 1, []Material{GetMaterialByName("Grass")}, make([]Residue, 0), Location{0, 0}},
}

func cloneItem(template Item) *Item {
	materialCopy := make([]Material, len(template.Material))
	copy(materialCopy, template.Material)

	residueCopy := make([]Residue, len(template.Residues))
	copy(residueCopy, template.Residues)

	return &Item{
		Name:      template.Name,
		Sharpness: template.Sharpness,
		Bluntness: template.Bluntness,
		Weight:    template.Weight,
		Material:  materialCopy,
		Residues:  residueCopy,
		Location:  template.Location,
	}
}

func createFromCatalog(catalog []Item, itemType string) *Item {
	for _, template := range catalog {
		if template.Name == itemType {
			return cloneItem(template)
		}
	}
	return nil
}

// CreateNewItem - Create a new item
func CreateNewItem(itemType string) *Item {
	if item := createFromCatalog(items, itemType); item != nil {
		return item
	}
	if item := createFromCatalog(BuildingMaterials, itemType); item != nil {
		return item
	}

	return nil
}
