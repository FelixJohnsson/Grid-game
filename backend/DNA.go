package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// --------------------------- Structures ----------------------------

// DNA represents the complete genetic code of an entity
type DNA struct {
	Chromosomes []*Chromosome
	Genes       map[string]*Gene
}

// Chromosome represents a single chromosome in the DNA
type Chromosome struct {
	ID       string
	Sequence []string // Each element is a codon (3-letter string)
	Genes    []*Gene  // Genes located on this chromosome
}

// Gene represents a specific gene that codes for traits
type Gene struct {
	Name        string
	StartCodon  int
	EndCodon    int
	Chromosome  *Chromosome
	Alleles     map[string]*Allele
	Expression  float64 // Gene expression level (0-1)
	Regulators  []*Gene // Genes that regulate this gene's expression
}

// Allele represents a specific variant of a gene
type Allele struct {
	Sequence    string
	Trait       *Trait
	Dominance   float64 // How dominant this allele is (0-1)
	Expression  float64 // How strongly this allele is expressed (0-1)
}

// Trait represents a physical or behavioral characteristic
type Trait struct {
	Name        string
	Type        TraitType
	Value       float64
	Modifiers   []*TraitModifier
}

// TraitType represents the category of a trait
type TraitType string

const (
	PhysicalTrait  TraitType = "Physical"
	BehavioralTrait TraitType = "Behavioral"
	MetabolicTrait TraitType = "Metabolic"
)

// TraitModifier represents how a trait affects entity attributes
type TraitModifier struct {
	Attribute string  // e.g., "Strength", "Intelligence"
	Modifier  float64 // Multiplier for the attribute
}

// -------------------------- DNA Functions --------------------------

// NewDNA creates a new DNA structure with random chromosomes
func NewDNA() *DNA {
	dna := &DNA{
		Chromosomes: make([]*Chromosome, 0),
		Genes:       make(map[string]*Gene),
	}
	
	// Create 23 pairs of chromosomes (like humans)
	for i := 0; i < 23; i++ {
		dna.Chromosomes = append(dna.Chromosomes, NewChromosome(fmt.Sprintf("chr%d", i+1)))
		dna.Chromosomes = append(dna.Chromosomes, NewChromosome(fmt.Sprintf("chr%d", i+1)))
	}
	
	return dna
}

// NewChromosome creates a new chromosome with random sequence
func NewChromosome(id string) *Chromosome {
	length := 1000 // Length of chromosome in codons
	sequence := make([]string, length)
	for i := range sequence {
		sequence[i] = MakeRandomCodon()
	}
	return &Chromosome{
		ID:       id,
		Sequence: sequence,
		Genes:    make([]*Gene, 0),
	}
}

// MakeRandomCodon generates a random 3-letter codon
func MakeRandomCodon() string {
	bases := []rune{'A', 'T', 'C', 'G'}
	return string([]rune{bases[rand.Intn(4)], bases[rand.Intn(4)], bases[rand.Intn(4)]})
}

// AddGene adds a new gene to the DNA
func (d *DNA) AddGene(name string, chromosomeID string, startCodon, endCodon int, alleles map[string]*Allele) *Gene {
	chromosome := d.findChromosome(chromosomeID)
	if chromosome == nil {
		return nil
	}
	
	gene := &Gene{
		Name:       name,
		StartCodon: startCodon,
		EndCodon:   endCodon,
		Chromosome: chromosome,
		Alleles:    alleles,
		Expression: 1.0,
	}
	
	chromosome.Genes = append(chromosome.Genes, gene)
	d.Genes[name] = gene
	return gene
}

// findChromosome finds a chromosome by ID
func (d *DNA) findChromosome(id string) *Chromosome {
	for _, chr := range d.Chromosomes {
		if chr.ID == id {
			return chr
		}
	}
	return nil
}

// ExpressGene determines which allele is expressed for a gene
func (d *DNA) ExpressGene(geneName string) *Allele {
	gene := d.Genes[geneName]
	if gene == nil {
		return nil
	}
	
	// Get the sequence for this gene
	sequence := strings.Join(gene.Chromosome.Sequence[gene.StartCodon:gene.EndCodon+1], "")
	
	// Find matching allele
	for _, allele := range gene.Alleles {
		if allele.Sequence == sequence {
			return allele
		}
	}
	
	return nil
}

// Mutate randomly changes a codon in the DNA
func (d *DNA) Mutate() {
	// Select random chromosome
	chr := d.Chromosomes[rand.Intn(len(d.Chromosomes))]
	
	// Select random position
	pos := rand.Intn(len(chr.Sequence))
	
	// Mutate codon
	chr.Sequence[pos] = MakeRandomCodon()
}

// Crossover performs genetic recombination between two DNA strands
func (d *DNA) Crossover(other *DNA) *DNA {
	child := NewDNA()
	
	for i := 0; i < len(d.Chromosomes); i += 2 {
		// Randomly choose which parent's chromosome to inherit
		if rand.Float64() < 0.5 {
			child.Chromosomes[i] = d.Chromosomes[i]
			child.Chromosomes[i+1] = other.Chromosomes[i+1]
		} else {
			child.Chromosomes[i] = other.Chromosomes[i]
			child.Chromosomes[i+1] = d.Chromosomes[i+1]
		}
	}
	
	return child
}

// ApplyTraits applies all expressed traits to an entity
func (d *DNA) ApplyTraits(entity *Entity) {
	for _, gene := range d.Genes {
		allele := d.ExpressGene(gene.Name)
		if allele != nil && allele.Trait != nil {
			for _, modifier := range allele.Trait.Modifiers {
				switch modifier.Attribute {
				case "Strength":
					entity.Strength = int(float64(entity.Strength) * modifier.Modifier)
				case "Agility":
					entity.Agility = int(float64(entity.Agility) * modifier.Modifier)
				case "Intelligence":
					entity.Intelligence = int(float64(entity.Intelligence) * modifier.Modifier)
				case "Charisma":
					entity.Charisma = int(float64(entity.Charisma) * modifier.Modifier)
				case "Stamina":
					entity.Stamina = int(float64(entity.Stamina) * modifier.Modifier)
				}
			}
		}
	}
}
