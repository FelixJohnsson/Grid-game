//go:build !display
// +build !display

package main

import "fmt"

// LaunchGame is a no-op in non-display builds.
func (w *World) LaunchGame(player *Entity) {
	fmt.Println("Display support is disabled in this build. Run with `go run -tags display . -display` to enable the Raylib window.")
}
