package main

// CleanEntity is a function that cleans a Person struct
func (w *World) CleanEntity(Entity *Entity) EntityCleaned {
	return EntityCleaned{
				FirstName:  Entity.FirstName,
				FamilyName: Entity.FamilyName,
				FullName:   Entity.FullName,
				Gender:     Entity.Gender,
				Age:        Entity.Age,
				Title:      Entity.Title,
				Location:   Entity.Location,
				Thinking:   Entity.Thinking,
				Head:       Entity.Body.Head,
				Torso:      Entity.Body.Torso,
				RightArm:   Entity.Body.RightArm,
				LeftArm:    Entity.Body.LeftArm,
				RightLeg:   Entity.Body.RightLeg,
				LeftLeg:    Entity.Body.LeftLeg,

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

				CurrentTask:      Entity.Brain.CurrentTask,
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