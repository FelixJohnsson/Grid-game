package main

// CleanPerson is a function that cleans a Person struct
func (w *World) CleanPerson(Person *Person) PersonCleaned {
	return PersonCleaned{
				FirstName:  Person.FirstName,
				FamilyName: Person.FamilyName,
				FullName:   Person.FullName,
				Gender:     Person.Gender,
				Age:        Person.Age,
				Title:      Person.Title,
				Location:   Person.Location,
				Thinking:   Person.Thinking,
				Head:       HeadCleaned{Person.Body.Head.LimbStatus},
				Torso:      Person.Body.Torso,
				RightArm:   Person.Body.RightArm,
				LeftArm:    Person.Body.LeftArm,
				RightLeg:   Person.Body.RightLeg,
				LeftLeg:    Person.Body.LeftLeg,

				Strength:         Person.Strength,
				Agility:          Person.Agility,
				Intelligence:     Person.Intelligence,
				Charisma:         Person.Charisma,
				Stamina:          Person.Stamina,
				CombatExperience: Person.CombatExperience,
				CombatSkill:      Person.CombatSkill,
				CombatStyle:      Person.CombatStyle,
				IsIncapacitated:  Person.IsIncapacitated,
				Relationships:    Person.Relationships,

				CurrentTask:      Person.Body.Head.Brain.CurrentTask,
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
			var cleanedPerson *PersonCleaned = nil
			var cleanedPlant *PlantCleaned = nil
			var shelter *Shelter = nil

			if tile.Person != nil {
				cleanedPersonVal := w.CleanPerson(tile.Person)
				cleanedPerson = &cleanedPersonVal
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
				Person:  cleanedPerson,
				Items:   tile.Items,
				Plant:   cleanedPlant,
				Shelter: shelter,
			}
		}
	}

	return cleanedTiles
}