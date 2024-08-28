package main

type Material struct {
	Name string
	Type string
	Hardness int
	Weight int
	Density int
	Malleability int
}

type Residue struct {
	Name string
	Amount int
}

type Item struct {
	Name string
	Sharpness int
	Bluntness int
	Weight int
	Material []Material
	Residues []Residue
}

type Inventory struct {
	Items []Item
}

// All available stone age materials
var materials = []Material{
	{"Wood", "Organic", 2, 1, 1, 1},
	{"Stone", "Inorganic", 5, 3, 4, 2},
	{"Leather", "Organic", 1, 1, 1, 1},
	{"Bone", "Organic", 3, 2, 2, 4},
	{"Flint", "Inorganic", 6, 4, 5, 3},
	{"Obsidian", "Inorganic", 7, 5, 6, 3},
	{"Feathers", "Organic", 1, 1, 1, 1},
	{"Clay", "Inorganic", 4, 2, 3, 1},
	{"Fur", "Organic", 1, 1, 1, 1},
	{"Wet Clay", "Inorganic", 4, 2, 3, 10},
	{"Water", "Inorganic", 1, 1, 1, 1},
	{"Soil", "Inorganic", 1, 1, 1, 1},
}

// All available items
var items = []Item{
	{"Wooden Spear", 8, 4, 5, []Material{materials[0]}, make([]Residue, 0)},
}

