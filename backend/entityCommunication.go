package main

import (
	"math/rand"
	"time"
)

const (
	relationshipFriendTalkThreshold = 10
	dialogueCooldownMs              = 7000
	maxRelationshipIntensity        = 100
	minRelationshipIntensity        = -100
)

var simpleDialogueLines = []string{
	"The weather is nice today.",
	"I found a good path nearby.",
	"Stay close to the base.",
	"We should keep gathering resources.",
	"I think this area is safe.",
	"Let's watch out for danger.",
	"I saw useful materials earlier.",
	"Good to see a friendly face.",
	"We are making steady progress.",
	"Keep moving and stay alert.",
}

func (b *Brain) ProcessSocialInteractions(observedTiles []Tile) {
	if b == nil || b.Owner == nil {
		return
	}

	seen := make(map[string]bool)
	for _, tile := range observedTiles {
		other := tile.Entity
		if other == nil || other.FullName == b.Owner.FullName {
			continue
		}
		if seen[other.FullName] {
			continue
		}
		seen[other.FullName] = true

		if other.Species != b.Owner.Species {
			continue
		}

		newIntensity := b.applyPersonalityRelationshipDelta(other)
		if newIntensity >= relationshipFriendTalkThreshold {
			b.tryDialogueWith(other)
		}
	}
}

func (b *Brain) applyPersonalityRelationshipDelta(other *Entity) int {
	if other == nil {
		return 0
	}

	delta := relationshipDeltaFromPersonality(other.Personalities)
	currentIntensity, exists := b.Owner.GetRelationshipIntensity(other.FullName)
	if !exists {
		b.Owner.AddRelationship(other, "Neutral", 0)
		currentIntensity = 0
	}

	nextIntensity := clampRelationshipIntensity(currentIntensity + delta)
	b.Owner.UpdateRelationship(other.FullName, relationshipLabel(nextIntensity), nextIntensity)
	return nextIntensity
}

func relationshipDeltaFromPersonality(personalities []Personality) int {
	if len(personalities) == 0 {
		return 1
	}
	for _, personality := range personalities {
		if personality == Hostile {
			return -1
		}
	}
	return 1
}

func relationshipLabel(intensity int) string {
	switch {
	case intensity >= 10:
		return "Friend"
	case intensity > 0:
		return "Acquaintance"
	case intensity <= -10:
		return "Enemy"
	case intensity < 0:
		return "Distrust"
	default:
		return "Neutral"
	}
}

func clampRelationshipIntensity(value int) int {
	if value > maxRelationshipIntensity {
		return maxRelationshipIntensity
	}
	if value < minRelationshipIntensity {
		return minRelationshipIntensity
	}
	return value
}

func (b *Brain) tryDialogueWith(other *Entity) {
	if other == nil || other.Brain == nil {
		return
	}

	nowMs := time.Now().UnixMilli()
	if nowMs-b.Owner.DialogueAtMs < dialogueCooldownMs {
		return
	}
	if nowMs-other.DialogueAtMs < dialogueCooldownMs {
		return
	}

	line := simpleDialogueLines[rand.Intn(len(simpleDialogueLines))]
	_ = b.SendDialogueMessage(other, line)
}

func (b *Brain) SendDialogueMessage(to *Entity, message string) bool {
	if b == nil || b.Owner == nil || to == nil || to.Brain == nil || message == "" {
		return false
	}

	nowMs := time.Now().UnixMilli()
	b.Owner.DialogueAtMs = nowMs
	b.Owner.DialogueWith = to.FullName
	b.Owner.LastDialogue = message
	b.Owner.DialogueMode = "Said"
	b.Owner.Thinking = `Said: "` + message + `"`
	b.AddMemoryToShortTerm("Talked", to.FullName, to.Location)

	return to.Brain.ReceiveDialogueMessage(b.Owner, message)
}

func (b *Brain) ReceiveDialogueMessage(from *Entity, message string) bool {
	if b == nil || b.Owner == nil || from == nil || message == "" {
		return false
	}

	nowMs := time.Now().UnixMilli()
	b.Owner.DialogueAtMs = nowMs
	b.Owner.DialogueWith = from.FullName
	b.Owner.LastDialogue = message
	b.Owner.DialogueMode = "Heard"
	b.Owner.Thinking = `Heard: "` + message + `"`
	b.AddMemoryToShortTerm("Heard dialogue", from.FullName, from.Location)
	return true
}
