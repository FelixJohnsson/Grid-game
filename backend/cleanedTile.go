package main

func (w *World) CleanTiles() [][]CleanedTile {
	tiles := w.GetTiles()
	cleanedTiles := make([][]CleanedTile, len(tiles))

	for y, row := range tiles {
		cleanedTiles[y] = make([]CleanedTile, len(row))
		for x, tile := range row {
			var cleanedBuilding *BuildingCleaned
			if tile.Building != nil {
				cleanedBuilding = &BuildingCleaned{
					Name:     tile.Building.Name,
					Type:     string(tile.Building.Type),
					Location: tile.Building.Location,
				}
			}

			var cleanedPersons []PersonCleaned
			for _, person := range tile.Persons {
				cleanedPersons = append(cleanedPersons, PersonCleaned{
					FirstName:        person.FirstName,
					FamilyName:       person.FamilyName,
					FullName:         person.FullName,
					Gender:           person.Gender,
					Age:              person.Age,
					Title:            person.Title,
					Location:         person.Location,
					Thinking:         person.Thinking,
					Body:             person.Body,
					Strength:         person.Strength,
					Agility:          person.Agility,
					Intelligence:     person.Intelligence,
					Charisma:         person.Charisma,
					Stamina:          person.Stamina,
					CombatExperience: person.CombatExperience,
					CombatSkill:      person.CombatSkill,
					CombatStyle:      person.CombatStyle,
					IsIncapacitated:  person.IsIncapacitated,
					Relationships:    person.Relationships,
				})
			}
			var cleanedPlants []*PlantCleaned
			for _, plant := range tile.Plants {
				// Remove the PlantLife from the Plant before sending it to the client
				cleanedPlants = append(cleanedPlants, &PlantCleaned{
					Name:          plant.Name,
					Age:           plant.Age,
					Health:        plant.Health,
					IsAlive:       plant.IsAlive,
					ProducesFruit: plant.ProducesFruit,
					Fruit:         plant.Fruit,
					PlantStage:    plant.PlantStage,
				})
			}

			cleanedTiles[y][x] = CleanedTile{
				Type:     tile.Type,
				Building: cleanedBuilding,
				Persons:  cleanedPersons,
				Items:    tile.Items,
				Plants:   cleanedPlants,
			}
		}
	}

	return cleanedTiles
}