package main

func cleanCognitiveMap(knownTiles map[Location]CognitiveMapTile) []CognitiveMapKnownTileCleaned {
	if len(knownTiles) == 0 {
		return nil
	}

	cleaned := make([]CognitiveMapKnownTileCleaned, 0, len(knownTiles))
	for location, tile := range knownTiles {
		var entity *CognitiveMapEntityCleaned
		if tile.Entity.FullName != "" || tile.Entity.SpeciesType != "" {
			entity = &CognitiveMapEntityCleaned{
				FullName:    tile.Entity.FullName,
				SpeciesType: tile.Entity.SpeciesType,
				IsAlive:     tile.Entity.IsAlive,
			}
		}

		var plant *CognitiveMapPlantCleaned
		if tile.Plant.Name != "" {
			plant = &CognitiveMapPlantCleaned{
				Name:          tile.Plant.Name,
				IsAlive:       tile.Plant.IsAlive,
				ProducesFruit: tile.Plant.ProducesFruit,
				PlantStage:    tile.Plant.PlantStage,
				FruitCount:    tile.Plant.FruitCount,
				HasRipeFruit:  tile.Plant.HasRipeFruit,
			}
		}

		var items map[string]int
		if len(tile.Items) > 0 {
			items = make(map[string]int, len(tile.Items))
			for itemName, count := range tile.Items {
				items[itemName] = count
			}
		}

		cleaned = append(cleaned, CognitiveMapKnownTileCleaned{
			Location:       location,
			TileType:       tile.TileType,
			Entity:         entity,
			Plant:          plant,
			Items:          items,
			LastSeenUnixMs: tile.LastSeenUnixMs,
		})
	}

	return cleaned
}

// CleanEntity is a function that cleans a Person struct
func (w *World) CleanEntity(Entity *Entity) EntityCleaned {
	var baseLocation *Location
	baseStockpile := map[string]int{}
	if Entity.Brain != nil && Entity.Brain.Base.Claimed {
		location := Entity.Brain.Base.Location
		baseLocation = &location
		for key, value := range Entity.Brain.Base.Stockpile {
			baseStockpile[key] = value
		}
	}
	if len(baseStockpile) == 0 {
		baseStockpile = nil
	}

	return EntityCleaned{
		FirstName:    Entity.FirstName,
		FamilyName:   Entity.FamilyName,
		FullName:     Entity.FullName,
		Gender:       Entity.Gender,
		Age:          Entity.Age,
		Title:        Entity.Title,
		Location:     Entity.Location,
		Thinking:     Entity.Thinking,
		LastDialogue: Entity.LastDialogue,
		DialogueWith: Entity.DialogueWith,
		DialogueAtMs: Entity.DialogueAtMs,
		DialogueMode: Entity.DialogueMode,
		Head:         Entity.Body.Head,
		Torso:        Entity.Body.Torso,
		RightArm:     Entity.Body.RightArm,
		LeftArm:      Entity.Body.LeftArm,
		RightLeg:     Entity.Body.RightLeg,
		LeftLeg:      Entity.Body.LeftLeg,

		Strength:         Entity.Strength,
		Agility:          Entity.Agility,
		Intelligence:     Entity.Intelligence,
		Charisma:         Entity.Charisma,
		Stamina:          Entity.Stamina,
		CombatExperience: Entity.CombatExperience,
		CombatSkill:      Entity.CombatSkill,
		CombatStyle:      Entity.CombatStyle,
		IsIncapacitated:  Entity.IsIncapacitated,
		Relationships:    Entity.Relationships,
		Personalities:    Entity.Personalities,

		CurrentTask:   Entity.Brain.CurrentTask,
		BaseClaimed:   Entity.Brain != nil && Entity.Brain.Base.Claimed,
		BaseLocation:  baseLocation,
		BaseStockpile: baseStockpile,
	}
}

// CleanPlant is a function that cleans a Plant struct
func (w *World) CleanPlant(Plant *Plant) PlantCleaned {
	return PlantCleaned{
		Name:          Plant.Name,
		Age:           Plant.Age,
		Health:        Plant.Health,
		IsAlive:       Plant.IsAlive,
		ProducesFruit: Plant.ProducesFruit,
		Fruit:         Plant.Fruit,
		PlantStage:    Plant.PlantStage,
	}
}

func (w *World) CleanTiles() [][]CleanedTile {
	tiles := w.GetTiles()
	cleanedTiles := make([][]CleanedTile, len(tiles))

	for y, row := range tiles {
		cleanedTiles[y] = make([]CleanedTile, len(row))
		for x, tile := range row {
			var cleanedEntity *EntityCleaned = nil
			var cleanedPlant *PlantCleaned = nil
			var shelter *Shelter = nil

			if tile.Entity != nil {
				cleanedEntityVal := w.CleanEntity(tile.Entity)
				cleanedEntity = &cleanedEntityVal
			}
			if tile.Plant != nil {
				cleanedPlantVal := w.CleanPlant(tile.Plant)
				cleanedPlant = &cleanedPlantVal
			}
			if tile.Shelter != nil {
				shelter = tile.Shelter
			}

			cleanedTiles[y][x] = CleanedTile{
				Type:    tile.Type,
				Entity:  cleanedEntity,
				Items:   tile.Items,
				Plant:   cleanedPlant,
				Shelter: shelter,
			}
		}
	}

	return cleanedTiles
}
