package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Display the map of the world in the console, super simple for now
func (w *World) DisplayMapInTerminal() {
	if true {
		fmt.Println()
		fmt.Println()
		for y := 0; y < w.Height; y++ {
			for x := 0; x < w.Width; x++ {
				if w.Tiles[x][y].Entity != nil {
					if w.Tiles[x][y].Entity != nil {
						fmt.Print(Red + string(w.Tiles[x][y].Entity.Species[0]) + Reset)
					} else {
						fmt.Print("X")
					}
				} else if w.Tiles[x][y].Plant != nil {
					if w.Tiles[x][y].Plant.Name == AppleTree {
						fmt.Print("A")
					} else if w.Tiles[x][y].Plant.Name == OakTree {
						fmt.Print(Yellow + "T" + Reset)
					} else if (w.Tiles[x][y].Plant.Name == HighGrass) {
						fmt.Print(Green + "H" + Reset)
					} else if (w.Tiles[x][y].Plant.Name == Flower) {
						fmt.Print(Red + "F" + Reset)
					}

				} else if w.Tiles[x][y].Type == 1 {
					fmt.Print(Blue + "W" + Reset)
				} else if w.Tiles[x][y].Shelter != nil {
					fmt.Print(Yellow + "S" + Reset)
				} else {
					fmt.Print(Green + "G" + Reset)
				}

			}
			fmt.Println()
		}
	}
}

// DisplayMap draws the game map to the screen using Raylib
func (w *World) DisplayMap(player *Entity) {
	var tileSize int32 = 10
	
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			tile := w.Tiles[x][y]
			posX := int32(int32(x) * tileSize)
			posY := int32(int32(y) * tileSize)

			if tile.Entity != nil {
				switch tile.Entity.Species {
				case Wolf:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Gray)
					rl.DrawText(string(tile.Entity.Species[0]), posX+5, posY+5, 10, rl.Black)

				case Human:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Red)
					rl.DrawText(string(tile.Entity.Species[0]), posX+5, posY+5, 10, rl.Black)
				}

			} else if tile.Plant != nil {
				// Draw different plants
				switch tile.Plant.Name {
				case AppleTree:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Beige)
					rl.DrawText("A", posX+5, posY+5, 10, rl.Black)
				case OakTree:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Brown)
					rl.DrawText("T", posX+5, posY+5, 10, rl.Black)
				case HighGrass:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.DarkGreen)
					rl.DrawText("H", posX+5, posY+5, 10, rl.Black)
				case Flower:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Magenta)
					rl.DrawText("F", posX+5, posY+5, 10, rl.Black)
				}
			} else if tile.Type == 1 {
				// Water tile
				rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Blue)
			} else if tile.Shelter != nil {
				// Draw shelter as yellow
				rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Gold)
				rl.DrawText("S", posX+5, posY+5, 10, rl.Black)
			} else {
				// Grass tile
				rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Green)
			}
		}
	}
}

// DisplayInfoPanel draws the information panel on the right side
func (w *World) DisplayInfoPanel(player *Entity, windowWidth, windowHeight int32) {
	if GlobalInfoPanel == nil {
		return
	}

	// Panel dimensions
	var panelWidth int32 = 400
	var panelX int32 = windowWidth - panelWidth
	
	// Draw panel background
	rl.DrawRectangle(panelX, 0, panelWidth, windowHeight, rl.LightGray)
	rl.DrawLine(panelX, 0, panelX, windowHeight, rl.DarkGray)
	
	// Draw panel title
	rl.DrawText("Information Panel", panelX+10, 10, 20, rl.Black)
	
	// Draw entity status section
	if player != nil {
		rl.DrawText("Entity Status", panelX+10, 40, 18, rl.DarkBlue)
		rl.DrawText(fmt.Sprintf("Name: %s", player.FullName), panelX+10, 65, 16, rl.Black)
		rl.DrawText(fmt.Sprintf("Position: (%d, %d)", player.Location.X, player.Location.Y), panelX+10, 85, 16, rl.Black)
		
		// Health status (with color)
		healthColor := rl.Green
		if player.Brain.PainLevel > player.Brain.PainTolerance/2 {
			healthColor = rl.Orange
		}
		if player.Brain.PainLevel > player.Brain.PainTolerance*3/4 {
			healthColor = rl.Red
		}
		rl.DrawText(fmt.Sprintf("Pain: %d/%d", player.Brain.PainLevel, player.Brain.PainTolerance), 
			panelX+10, 105, 16, healthColor)
		
		// Hunger and thirst
		hungerColor := rl.Green
		if player.Brain.PhysiologicalNeeds.Hunger > 30 {
			hungerColor = rl.Orange
		}
		if player.Brain.PhysiologicalNeeds.Hunger > 70 {
			hungerColor = rl.Red
		}
		rl.DrawText(fmt.Sprintf("Hunger: %d", player.Brain.PhysiologicalNeeds.Hunger), 
			panelX+10, 125, 16, hungerColor)
		
		thirstColor := rl.Green
		if player.Brain.PhysiologicalNeeds.Thirst > 30 {
			thirstColor = rl.Orange
		}
		if player.Brain.PhysiologicalNeeds.Thirst > 70 {
			thirstColor = rl.Red
		}
		rl.DrawText(fmt.Sprintf("Thirst: %d", player.Brain.PhysiologicalNeeds.Thirst), 
			panelX+10, 145, 16, thirstColor)
		
		// Current task
		if player.Brain.CurrentTask.IsActive {
			rl.DrawText(fmt.Sprintf("Task: %s → %s", 
				player.Brain.CurrentTask.Action, player.Brain.CurrentTask.Target), 
				panelX+10, 165, 16, rl.DarkPurple)
		} else {
			rl.DrawText("No active task", panelX+10, 165, 16, rl.Gray)
		}
		
		// Current thought
		if len(player.Thinking) > 0 {
			thought := player.Thinking
			if len(thought) > 35 {
				thought = thought[:32] + "..."
			}
			rl.DrawText(fmt.Sprintf("Thinking: %s", thought), panelX+10, 185, 16, rl.DarkPurple)
		}
	}
	
	// Draw message log section
	rl.DrawText("Message Log", panelX+10, 220, 18, rl.DarkBlue)
	rl.DrawLine(panelX+5, 245, panelX+panelWidth-10, 245, rl.DarkGray)
	
	messages := GlobalInfoPanel.GetMessages()
	
	var startY int32 = 250
	var lineHeight int32 = 20
	maxMsgsVisible := (windowHeight - startY - 10) / lineHeight
	
	// Display most recent messages
	numToShow := int(maxMsgsVisible)
	if numToShow > len(messages) {
		numToShow = len(messages)
	}
	
	for i := 0; i < numToShow; i++ {
		msg := messages[i]
		
		// Choose color based on message type
		color := rl.Black
		switch msg.Type {
		case INFO:
			color = rl.DarkBlue
		case WARNING:
			color = rl.Orange
		case ERROR:
			color = rl.Red
		case SUCCESS:
			color = rl.Green
		case HEALTH:
			color = rl.Purple
		case FOOD:
			color = rl.Brown
		case WATER:
			color = rl.Blue
		case COMBAT:
			color = rl.Maroon
		case DEATH:
			color = rl.Black
		}
		
		// Format message text
		text := msg.Text
		if len(text) > 35 {
			text = text[:32] + "..."
		}
		
		// Calculate display position
		yPos := startY + int32(i)*lineHeight
		
		// Display message
		timeStr := msg.Timestamp.Format("15:04:05")
		rl.DrawText(timeStr, panelX+10, yPos, 12, rl.DarkGray)
		rl.DrawText(text, panelX+90, yPos, 16, color)
	}
}

func (w *World) LaunchGame(player *Entity) {
	// Initialize the InfoPanel
	InitInfoPanel(100)

	// Log initial messages
	LogInfo("Game started", "System")
	LogInfo("Welcome to the simulation", "System")
	LogInfo(fmt.Sprintf("Player character: %s", player.FullName), "System")
	
	// Window dimensions
	windowWidth := int32(1300)
	windowHeight := int32(1000)
	
	// Initialize the Raylib window
	rl.InitWindow(windowWidth, windowHeight, "Grid Game Simulation")
	defer rl.CloseWindow()

	// Set the target FPS to 60
	rl.SetTargetFPS(60)

	// Create a ticker that ticks every 100ms
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	lastUpdate := time.Now()

	// Game loop
	for !rl.WindowShouldClose() {
		select {
		case <-ticker.C:
			// It's time to render a new frame
			rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)
			
			// Call the DisplayMap function to draw the world map
			w.DisplayMap(player)
			
			// Display the information panel
			w.DisplayInfoPanel(player, windowWidth, windowHeight)
			
			// Display time since last update
			timeSinceUpdate := time.Since(lastUpdate)
			rl.DrawText(timeSinceUpdate.String(), 10, 10, 12, rl.Black)
			
			lastUpdate = time.Now()
			
			rl.EndDrawing()
		default:
			// No tick yet, let's yield to the OS to avoid busy-waiting
			time.Sleep(time.Millisecond)
		}
	}
}
