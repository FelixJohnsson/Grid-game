package main

import (
	"context"
	"time"
)

const (
	plantLifecycleTickInterval = 2 * time.Second
	appleTreeMaxFruit          = 6
	appleTreeSpawnEveryTicks   = 2
	appleFruitRipenAfterTicks  = 1
	appleFruitNutritionalValue = 20
)

// StartPlantLifecycle starts a single world-level lifecycle loop for plants.
func (w *World) StartPlantLifecycle() {
	w.StopPlantLifecycle()

	ctx, cancel := context.WithCancel(context.Background())
	w.mu.Lock()
	w.plantLifecycleCtx = ctx
	w.plantLifecycleCancel = cancel
	w.mu.Unlock()

	go w.runPlantLifecycleLoop(ctx)
}

// StopPlantLifecycle stops the active world-level lifecycle loop, if any.
func (w *World) StopPlantLifecycle() {
	w.mu.Lock()
	cancel := w.plantLifecycleCancel
	w.plantLifecycleCtx = nil
	w.plantLifecycleCancel = nil
	w.mu.Unlock()

	if cancel != nil {
		cancel()
	}
}

func (w *World) runPlantLifecycleLoop(ctx context.Context) {
	ticker := time.NewTicker(plantLifecycleTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.TickPlantLifecycle()
		}
	}
}

// TickPlantLifecycle advances plant age and fruiting state in one simulation tick.
func (w *World) TickPlantLifecycle() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for y := range w.Tiles {
		for x := range w.Tiles[y] {
			plant := w.Tiles[y][x].Plant
			if plant == nil || !plant.IsAlive {
				continue
			}

			plant.Age++
			if !plant.ProducesFruit {
				continue
			}

			// Ripen any fruit that reached its configured ripe age.
			for i := range plant.Fruit {
				if !plant.Fruit[i].IsRipe && plant.Age >= plant.Fruit[i].RipeAge {
					plant.Fruit[i].IsRipe = true
				}
			}

			if plant.Name != AppleTree {
				continue
			}

			if len(plant.Fruit) >= appleTreeMaxFruit {
				continue
			}
			if plant.Age%appleTreeSpawnEveryTicks != 0 {
				continue
			}

			plant.Fruit = append(plant.Fruit, CreateNewFruit(
				"Apple",
				plant.Age+appleFruitRipenAfterTicks,
				false,
				appleFruitNutritionalValue,
			))
		}
	}
}
