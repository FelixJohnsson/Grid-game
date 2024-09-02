package main

func (w *World) CleanTiles() [][]CleanedTile {
	tiles := w.GetTiles()
	cleanedTiles := make([][]CleanedTile, len(tiles))

	for y, row := range tiles {
		cleanedTiles[y] = make([]CleanedTile, len(row))
		for x, tile := range row {
			var cleanedPersons []PersonCleaned
			cleanedPersons = append(cleanedPersons, PersonCleaned{
				FirstName:  tile.Person.FirstName,
				FamilyName: tile.Person.FamilyName,
				FullName:   tile.Person.FullName,
				Gender:     tile.Person.Gender,
				Age:        tile.Person.Age,
				Title:      tile.Person.Title,
				Location:   tile.Person.Location,
				Thinking:   tile.Person.Thinking,
				Head:       HeadCleaned{tile.Person.Body.Head.LimbStatus},
				Torso:      tile.Person.Body.Torso,
				RightArm:   tile.Person.Body.RightArm,
				LeftArm:    tile.Person.Body.LeftArm,
				RightLeg:   tile.Person.Body.RightLeg,
				LeftLeg:    tile.Person.Body.LeftLeg,

				Strength:         tile.Person.Strength,
				Agility:          tile.Person.Agility,
				Intelligence:     tile.Person.Intelligence,
				Charisma:         tile.Person.Charisma,
				Stamina:          tile.Person.Stamina,
				CombatExperience: tile.Person.CombatExperience,
				CombatSkill:      tile.Person.CombatSkill,
				CombatStyle:      tile.Person.CombatStyle,
				IsIncapacitated:  tile.Person.IsIncapacitated,
				Relationships:    tile.Person.Relationships,
			})
			var cleanedPlant *PlantCleaned
			// Remove the PlantLife from the Plant before sending it to the client
			plant := tile.Plant
			cleanedPlant = &PlantCleaned{
				Name:          plant.Name,
				Age:           plant.Age,
				Health:        plant.Health,
				IsAlive:       plant.IsAlive,
				ProducesFruit: plant.ProducesFruit,
				Fruit:         plant.Fruit,
				PlantStage:    plant.PlantStage,
			}

			cleanedTiles[y][x] = CleanedTile{
				Type:    tile.Type,
				Persons: cleanedPersons,
				Items:   tile.Items,
				Plant:   cleanedPlant,
			}
		}
	}

	return cleanedTiles
}