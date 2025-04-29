package main

import (
	"fmt"
	"sync"
	"time"
)

// MessageType represents the type of message
type MessageType int

const (
	// Message types
	INFO    MessageType = iota // Regular information
	WARNING                    // Warning message
	ERROR                      // Error message
	SUCCESS                    // Success message
	HEALTH                     // Health-related message
	FOOD                       // Food-related message
	WATER                      // Water-related message
	COMBAT                     // Combat-related message
	DEATH                      // Death-related message
)

// Message represents a structured message for the information panel
type Message struct {
	Text      string      // The message text
	Type      MessageType // The type of message
	Timestamp time.Time   // When the message was created
	EntityID  string      // ID of the entity related to this message (optional)
}

// InfoPanel manages messages for display
type InfoPanel struct {
	messages      []Message
	maxMessages   int
	mutex         sync.Mutex
	messageAdded  chan bool // Signal when a new message is added
}

// Global InfoPanel instance
var GlobalInfoPanel *InfoPanel

// InitInfoPanel initializes the InfoPanel
func InitInfoPanel(maxMessages int) *InfoPanel {
	panel := &InfoPanel{
		messages:     make([]Message, 0, maxMessages),
		maxMessages:  maxMessages,
		messageAdded: make(chan bool, 10),
	}
	
	// Store in global variable for easy access
	GlobalInfoPanel = panel
	
	return panel
}

// AddMessage adds a new message to the panel
func (ip *InfoPanel) AddMessage(text string, msgType MessageType, entityID string) {
	ip.mutex.Lock()
	defer ip.mutex.Unlock()
	
	// Create a new message
	msg := Message{
		Text:      text,
		Type:      msgType,
		Timestamp: time.Now(),
		EntityID:  entityID,
	}
	
	// Add the message to the beginning of the slice so newest are first
	ip.messages = append([]Message{msg}, ip.messages...)
	
	// If we have too many messages, remove the oldest ones
	if len(ip.messages) > ip.maxMessages {
		ip.messages = ip.messages[:ip.maxMessages]
	}
	
	// Notify about the new message
	select {
	case ip.messageAdded <- true:
		// Signal sent
	default:
		// Channel buffer is full, can't send
	}
	
	// Also print to console for debugging
	fmt.Printf("[%s] %s: %s\n", msgTypeToString(msgType), entityID, text)
}

// GetMessages returns all messages
func (ip *InfoPanel) GetMessages() []Message {
	ip.mutex.Lock()
	defer ip.mutex.Unlock()
	
	// Return a copy to avoid race conditions
	result := make([]Message, len(ip.messages))
	copy(result, ip.messages)
	
	return result
}

// GetMessagesForEntity returns messages related to a specific entity
func (ip *InfoPanel) GetMessagesForEntity(entityID string) []Message {
	ip.mutex.Lock()
	defer ip.mutex.Unlock()
	
	var result []Message
	for _, msg := range ip.messages {
		if msg.EntityID == entityID {
			result = append(result, msg)
		}
	}
	
	return result
}

// Helper function to convert MessageType to string
func msgTypeToString(msgType MessageType) string {
	switch msgType {
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case SUCCESS:
		return "SUCCESS"
	case HEALTH:
		return "HEALTH"
	case FOOD:
		return "FOOD"
	case WATER:
		return "WATER"
	case COMBAT:
		return "COMBAT"
	case DEATH:
		return "DEATH"
	default:
		return "UNKNOWN"
	}
}

// Helper global functions to make logging easier

// LogInfo adds an info message to the global info panel
func LogInfo(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, INFO, entityID)
	}
}

// LogWarning adds a warning message to the global info panel
func LogWarning(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, WARNING, entityID)
	}
}

// LogError adds an error message to the global info panel
func LogError(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, ERROR, entityID)
	}
}

// LogSuccess adds a success message to the global info panel
func LogSuccess(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, SUCCESS, entityID)
	}
}

// LogHealth adds a health-related message to the global info panel
func LogHealth(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, HEALTH, entityID)
	}
}

// LogFood adds a food-related message to the global info panel
func LogFood(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, FOOD, entityID)
	}
}

// LogWater adds a water-related message to the global info panel
func LogWater(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, WATER, entityID)
	}
}

// LogCombat adds a combat-related message to the global info panel
func LogCombat(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, COMBAT, entityID)
	}
}

// LogDeath adds a death-related message to the global info panel
func LogDeath(text string, entityID string) {
	if GlobalInfoPanel != nil {
		GlobalInfoPanel.AddMessage(text, DEATH, entityID)
	}
} 