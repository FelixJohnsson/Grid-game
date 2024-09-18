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
func (w *World) DisplayMap() {
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

func (w *World) LaunchGame() {
	// Initialize the Raylib window
	rl.InitWindow(300, 300, "The game")
	defer rl.CloseWindow()

	// Set the target FPS to 60
	rl.SetTargetFPS(60)

	// Create a ticker that ticks every 100ms
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Game loop
	for !rl.WindowShouldClose() {
		select {
		case <-ticker.C:
			// It's time to render a new frame
			rl.BeginDrawing()
			rl.ClearBackground(rl.RayWhite)
			
			// Call the DisplayMap function to draw the world map
			w.DisplayMap()
			
			rl.EndDrawing()
		default:
			// No tick yet, let's yield to the OS to avoid busy-waiting
			time.Sleep(time.Millisecond)
		}
	}
}
