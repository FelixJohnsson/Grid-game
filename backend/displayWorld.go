//go:build display
// +build display

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
					} else if w.Tiles[x][y].Plant.Name == HighGrass {
						fmt.Print(Green + "H" + Reset)
					} else if w.Tiles[x][y].Plant.Name == Flower {
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

// DisplayCognitiveMap draws a mini-map based on the player's cognitive map
func (w *World) DisplayCognitiveMap(player *Entity, windowWidth, windowHeight int32) {
	// Return if player or cognitive map is nil
	if player == nil || player.Brain == nil {
		return
	}
	knownTiles := player.Brain.GetKnownTilesSnapshot()
	if len(knownTiles) == 0 {
		return
	}

	// Mini-map configuration
	var mapWidth int32 = 200
	var mapHeight int32 = 200
	var mapX int32 = 20
	var mapY int32 = windowHeight - mapHeight - 200
	var tileSize int32 = 4 // Smaller tiles for mini-map

	// Draw mini-map background and border
	rl.DrawRectangle(mapX-5, mapY-5, mapWidth+10, mapHeight+10, rl.DarkGray)
	rl.DrawRectangle(mapX, mapY, mapWidth, mapHeight, rl.Black)

	// Draw title
	rl.DrawText("Cognitive Map", mapX, mapY-20, 16, rl.DarkBlue)

	// Calculate the center of the mini-map (in tile coordinates)
	centerX := player.Location.X
	centerY := player.Location.Y

	// Calculate visible tile range
	tilesVisibleX := int(mapWidth / tileSize)
	tilesVisibleY := int(mapHeight / tileSize)

	minX := centerX - tilesVisibleX/2
	maxX := centerX + tilesVisibleX/2
	minY := centerY - tilesVisibleY/2
	maxY := centerY + tilesVisibleY/2

	// Draw the known tiles from cognitive map
	for loc, tile := range knownTiles {
		// Check if this tile is in the visible range
		if loc.X >= minX && loc.X <= maxX && loc.Y >= minY && loc.Y <= maxY {
			// Calculate position on mini-map with X and Y swapped to fix rotation
			// This flips the map to match main map orientation
			posX := mapX + int32(loc.Y-minY)*tileSize
			posY := mapY + int32(loc.X-minX)*tileSize

			// Draw tile based on cognitive map knowledge
			if tile.Entity.IsAlive {
				// Draw entities
				switch tile.Entity.SpeciesType {
				case Wolf:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Gray)
				case Human:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Red)
				default:
					// Unknown entity
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Purple)
				}
			} else if tile.Plant.IsAlive {
				// Draw plants
				switch tile.Plant.Name {
				case AppleTree:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Beige)
				case OakTree:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Brown)
				case HighGrass:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.DarkGreen)
				case Flower:
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Magenta)
				default:
					// Unknown plant
					rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Green)
				}
			} else if tile.TileType == 0 {
				rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Green)
			} else if tile.TileType == 1 {
				rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Blue)
			} else {
				rl.DrawRectangle(posX, posY, tileSize, tileSize, rl.Gray)
			}
		}
	}

	// Draw player position (center of mini-map)
	var playerPosX int32 = mapX + mapWidth/2
	var playerPosY int32 = mapY + mapHeight/2

	// Draw player marker (white dot with border)
	rl.DrawCircle(playerPosX, playerPosY, 3, rl.Black)
	rl.DrawCircle(playerPosX, playerPosY, 2, rl.White)

	// Add legend
	legendY := mapY + mapHeight + 10
	rl.DrawText("Legend:", mapX, legendY, 14, rl.Black)

	// First row of legend items
	legendY += 20
	rl.DrawRectangle(mapX, legendY, 10, 10, rl.Red)
	rl.DrawText("Human", mapX+15, legendY, 12, rl.Black)

	rl.DrawRectangle(mapX+70, legendY, 10, 10, rl.Gray)
	rl.DrawText("Wolf", mapX+85, legendY, 12, rl.Black)

	// Second row
	legendY += 15
	rl.DrawRectangle(mapX, legendY, 10, 10, rl.Brown)
	rl.DrawText("Tree", mapX+15, legendY, 12, rl.Black)

	rl.DrawRectangle(mapX+70, legendY, 10, 10, rl.Blue)
	rl.DrawText("Water", mapX+85, legendY, 12, rl.Black)
}

// DisplayMap draws the game map to the screen using Raylib
func (w *World) DisplayMap(player *Entity) {
	var tileSize int32 = 8

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
		textColor := rl.Black

		rl.DrawText(fmt.Sprintf("Pain: %d/%d", player.Brain.PainLevel, player.Brain.PainTolerance),
			panelX+10, 105, 16, textColor)

		rl.DrawText(fmt.Sprintf("Hunger: %f", player.Body.Blood.Glucose),
			panelX+10, 125, 16, textColor)

		rl.DrawText(fmt.Sprintf("Thirst: %f", player.Body.Blood.Water),
			panelX+10, 145, 16, textColor)

		rl.DrawText(fmt.Sprintf("Oxygen: %f", player.Body.Blood.Oxygen),
			panelX+10, 165, 16, textColor)

		rl.DrawText(fmt.Sprintf("Toxins: %f", player.Body.Blood.Toxins),
			panelX+10, 185, 16, textColor)

		rl.DrawText(fmt.Sprintf("Adrenaline: %d", player.Body.Blood.Hormones.Adrenaline),
			panelX+10, 205, 16, textColor)

		rl.DrawText(fmt.Sprintf("Cortisol: %d", player.Body.Blood.Hormones.Cortisol),
			panelX+10, 225, 16, textColor)

		rl.DrawText(fmt.Sprintf("Dopamine: %d", player.Body.Blood.Hormones.Dopamine),
			panelX+10, 245, 16, textColor)

		rl.DrawText(fmt.Sprintf("Epinephrine: %d", player.Body.Blood.Hormones.Epinephrine),
			panelX+10, 265, 16, textColor)

		rl.DrawText(fmt.Sprintf("Endorphin: %d", player.Body.Blood.Hormones.Endorphin),
			panelX+10, 285, 16, textColor)

		rl.DrawText(fmt.Sprintf("Serotonin: %d", player.Body.Blood.Hormones.Serotonin),
			panelX+10, 305, 16, textColor)

		// Current task
		if player.Brain.CurrentTask.IsActive {
			rl.DrawText(fmt.Sprintf("Task: %s → %s",
				player.Brain.CurrentTask.Action, player.Brain.CurrentTask.Target),
				panelX+10, 325, 16, rl.DarkPurple)
		} else {
			rl.DrawText("No active task", panelX+10, 325, 16, rl.Gray)
		}

		// Current thought
		if len(player.Thinking) > 0 {
			thought := player.Thinking
			if len(thought) > 35 {
				thought = thought[:32] + "..."
			}
			rl.DrawText(fmt.Sprintf("Thinking: %s", thought), panelX+10, 345, 16, rl.DarkPurple)
		}
	}
	var startY int32 = 370

	// Draw message log section
	rl.DrawText("Message Log", panelX+10, startY, 18, rl.DarkBlue)
	rl.DrawLine(panelX+5, startY, panelX+panelWidth-10, startY, rl.DarkGray)

	messages := GlobalInfoPanel.GetMessages()

	var lineHeight int32 = 20
	maxMsgsVisible := (windowHeight - startY + 50) / lineHeight

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
		yPos := startY + 40 + int32(i)*lineHeight

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
	windowWidth := int32(1200)
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

			// Display the cognitive mini-map
			w.DisplayCognitiveMap(player, windowWidth, windowHeight)

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
