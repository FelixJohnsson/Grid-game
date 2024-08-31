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

// All available items
var items = []Item{
	// Weapons
	{"Wooden Spear", 5, 4, 3, []Material{materials[0]}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Staff", 1, 8, 2, []Material{materials[0]}, make([]Residue, 0), Location{0, 0}},
	{"Stone Axe", 6, 2, 5, []Material{materials[1], materials[0], materials[2]}, make([]Residue, 0), Location{0, 0}},

	// Tools

	// Storage
	{"Wooden Basket", 1, 1, 1, []Material{materials[0]}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Box", 1, 1, 1, []Material{materials[0]}, make([]Residue, 0), Location{0, 0}},
	{"Wooden Crate", 1, 1, 1, []Material{materials[0]}, make([]Residue, 0), Location{0, 0}},
	{"Woven Grass Basket", 1, 1, 1, []Material{materials[0]}, make([]Residue, 0), Location{0, 0}},
}
