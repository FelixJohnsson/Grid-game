package main

import "testing"

func TestProcessSocialInteractionsFriendlyPersonalityIncreasesRelationship(t *testing.T) {
	world := NewWorld(12, 12)
	brain := newVisionTestBrain(t, world, 2, 2)
	other := newVisionTestBrain(t, world, 3, 2)

	brain.Owner.Species = Human
	other.Owner.Species = Human
	other.Owner.Personalities = []Personality{Friendly}

	obs := world.GetVision(brain.Owner.Location.X, brain.Owner.Location.Y, brain.Owner.VisionRange)
	brain.ProcessSocialInteractions(obs)

	intensity, ok := brain.Owner.GetRelationshipIntensity(other.Owner.FullName)
	if !ok {
		t.Fatalf("expected relationship to be created")
	}
	if intensity != 1 {
		t.Fatalf("expected relationship intensity to increase by +1, got %d", intensity)
	}
}

func TestProcessSocialInteractionsHostilePersonalityDecreasesRelationship(t *testing.T) {
	world := NewWorld(12, 12)
	brain := newVisionTestBrain(t, world, 2, 2)
	other := newVisionTestBrain(t, world, 3, 2)

	brain.Owner.Species = Human
	other.Owner.Species = Human
	other.Owner.Personalities = []Personality{Hostile}

	obs := world.GetVision(brain.Owner.Location.X, brain.Owner.Location.Y, brain.Owner.VisionRange)
	brain.ProcessSocialInteractions(obs)

	intensity, ok := brain.Owner.GetRelationshipIntensity(other.Owner.FullName)
	if !ok {
		t.Fatalf("expected relationship to be created")
	}
	if intensity != -1 {
		t.Fatalf("expected relationship intensity to decrease by -1, got %d", intensity)
	}
}

func TestProcessSocialInteractionsTriggersDialogueAtFriendThreshold(t *testing.T) {
	world := NewWorld(12, 12)
	brain := newVisionTestBrain(t, world, 2, 2)
	other := newVisionTestBrain(t, world, 3, 2)

	brain.Owner.Species = Human
	other.Owner.Species = Human
	other.Owner.Personalities = []Personality{Friendly}

	brain.Owner.AddRelationship(other.Owner, "Acquaintance", 9)
	brain.Owner.DialogueAtMs = 0
	other.Owner.DialogueAtMs = 0

	obs := world.GetVision(brain.Owner.Location.X, brain.Owner.Location.Y, brain.Owner.VisionRange)
	brain.ProcessSocialInteractions(obs)

	intensity, _ := brain.Owner.GetRelationshipIntensity(other.Owner.FullName)
	if intensity != 10 {
		t.Fatalf("expected relationship to reach talk threshold 10, got %d", intensity)
	}
	if brain.Owner.LastDialogue == "" {
		t.Fatalf("expected sender to record last dialogue line")
	}
	if other.Owner.LastDialogue == "" {
		t.Fatalf("expected receiver to record last dialogue line")
	}
	if other.Owner.DialogueWith != brain.Owner.FullName {
		t.Fatalf("expected receiver dialogue partner %q, got %q", brain.Owner.FullName, other.Owner.DialogueWith)
	}
}

func TestProcessSocialInteractionsOnlyAffectsSameSpecies(t *testing.T) {
	world := NewWorld(12, 12)
	brain := newVisionTestBrain(t, world, 2, 2)
	other := newVisionTestBrain(t, world, 3, 2)

	brain.Owner.Species = Human
	other.Owner.Species = Wolf
	other.Owner.Personalities = []Personality{Friendly}

	obs := world.GetVision(brain.Owner.Location.X, brain.Owner.Location.Y, brain.Owner.VisionRange)
	brain.ProcessSocialInteractions(obs)

	if _, ok := brain.Owner.GetRelationshipIntensity(other.Owner.FullName); ok {
		t.Fatalf("expected no relationship updates for different species")
	}
}
