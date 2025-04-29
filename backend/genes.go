package main

// --------------------------- Example Genes ----------------------------

// Define some example genes that affect physical traits
var muscleFiberGene = &Gene{
	Name:       "MuscleFiber",
	StartCodon: 100,
	EndCodon:   102,
	Alleles: map[string]*Allele{
		"ATGCGTACG": {
			Sequence: "ATGCGTACG",
			Trait: &Trait{
				Name:  "Fast Twitch Muscles",
				Type:  PhysicalTrait,
				Value: 1.5,
				Modifiers: []*TraitModifier{
					{Attribute: "Strength", Modifier: 1.3},
					{Attribute: "Agility", Modifier: 1.2},
					{Attribute: "Stamina", Modifier: 0.8},
				},
			},
			Dominance:  0.7,
			Expression: 1.0,
		},
		"TTGAAATTT": {
			Sequence: "TTGAAATTT",
			Trait: &Trait{
				Name:  "Slow Twitch Muscles",
				Type:  PhysicalTrait,
				Value: 1.0,
				Modifiers: []*TraitModifier{
					{Attribute: "Strength", Modifier: 0.8},
					{Attribute: "Agility", Modifier: 0.9},
					{Attribute: "Stamina", Modifier: 1.3},
				},
			},
			Dominance:  0.3,
			Expression: 1.0,
		},
	},
}

var brainSizeGene = &Gene{
	Name:       "BrainSize",
	StartCodon: 200,
	EndCodon:   202,
	Alleles: map[string]*Allele{
		"CGTACGCGT": {
			Sequence: "CGTACGCGT",
			Trait: &Trait{
				Name:  "Large Brain",
				Type:  PhysicalTrait,
				Value: 1.3,
				Modifiers: []*TraitModifier{
					{Attribute: "Intelligence", Modifier: 1.4},
					{Attribute: "Charisma", Modifier: 1.2},
					{Attribute: "Stamina", Modifier: 0.9},
				},
			},
			Dominance:  0.6,
			Expression: 1.0,
		},
		"ATGCGTACG": {
			Sequence: "ATGCGTACG",
			Trait: &Trait{
				Name:  "Normal Brain",
				Type:  PhysicalTrait,
				Value: 1.0,
				Modifiers: []*TraitModifier{
					{Attribute: "Intelligence", Modifier: 1.0},
					{Attribute: "Charisma", Modifier: 1.0},
					{Attribute: "Stamina", Modifier: 1.0},
				},
			},
			Dominance:  0.4,
			Expression: 1.0,
		},
	},
}

var heartSizeGene = &Gene{
	Name:       "HeartSize",
	StartCodon: 300,
	EndCodon:   302,
	Alleles: map[string]*Allele{
		"TTGAAATTT": {
			Sequence: "TTGAAATTT",
			Trait: &Trait{
				Name:  "Large Heart",
				Type:  PhysicalTrait,
				Value: 1.4,
				Modifiers: []*TraitModifier{
					{Attribute: "Stamina", Modifier: 1.5},
					{Attribute: "Strength", Modifier: 1.2},
					{Attribute: "Agility", Modifier: 1.1},
				},
			},
			Dominance:  0.7,
			Expression: 1.0,
		},
		"CGTACGCGT": {
			Sequence: "CGTACGCGT",
			Trait: &Trait{
				Name:  "Normal Heart",
				Type:  PhysicalTrait,
				Value: 1.0,
				Modifiers: []*TraitModifier{
					{Attribute: "Stamina", Modifier: 1.0},
					{Attribute: "Strength", Modifier: 1.0},
					{Attribute: "Agility", Modifier: 1.0},
				},
			},
			Dominance:  0.3,
			Expression: 1.0,
		},
	},
}

// --------------------------- Helper Functions ----------------------------

// InitializeDNA creates a new DNA with some example genes
func InitializeDNA() *DNA {
	dna := NewDNA()
	
	// Add our example genes to the DNA
	dna.AddGene("MuscleFiber", "chr1", 100, 102, muscleFiberGene.Alleles)
	dna.AddGene("BrainSize", "chr2", 200, 202, brainSizeGene.Alleles)
	dna.AddGene("HeartSize", "chr3", 300, 302, heartSizeGene.Alleles)
	
	return dna
}