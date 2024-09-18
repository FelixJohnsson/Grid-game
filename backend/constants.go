package main

import "time"

var SIZE_OF_MAP = 50
var SIMULATION_RATE time.Duration = 5
var SIMULATION_TIME time.Duration = SIMULATION_RATE * time.Millisecond

// ------------------ Display -------------------
var TILE_SIZE = 10
var SCREEN_SIZE = SIZE_OF_MAP * TILE_SIZE

// ------------------ Plants --------------------

var SUSTAIN_COST = 3
var MINIMUM_NUTRIENT_STORAGE = 200
var MINIMUM_NUTRIENT_STORAGE_FOR_GROWTH = 100
var GROW_UNTIL_STAGE PlantStage  = 5
var FRUITING_COST = 30
var NUTRIENTS_TAKEN = 10

// ------------------ Tile ----------------------

var TILE_NUTRIENT_NEXT_TO_WATER = 15
var TILE_NUTRIENT_NEXT_TO_GRASS = 2
var MAX_TILE_NUTRIENT = 3000
