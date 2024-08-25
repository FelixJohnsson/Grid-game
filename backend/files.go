package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ------------------------------------------ People ------------------------------------------

// loadPersonsFromFile loads the array of persons from the file
func loadPersonsFromFile() ([]Person, error) {
	file, err := os.Open("person.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty array if the file does not exist
			return []Person{}, nil
		}
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var persons []Person
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&persons); err != nil {
		return nil, fmt.Errorf("error decoding JSON from file: %v", err)
	}

	return persons, nil
}

// savePersonsToFile saves the array of persons to the file
func savePersonsToFile(persons []Person) error {
	file, err := os.Create("person.json")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON
	if err := encoder.Encode(persons); err != nil {
		return fmt.Errorf("error encoding JSON to file: %v", err)
	}

	return nil
}

// addPerson adds a new person to the array and saves it to the file
func addPerson(p Person) error {
	// Load existing persons
	persons, err := loadPersonsFromFile()
	if err != nil {
		return fmt.Errorf("error loading persons from file: %v", err)
	}

	// Append the new person
	persons = append(persons, p)

	// Save the updated array back to the file
	if err := savePersonsToFile(persons); err != nil {
		return fmt.Errorf("error saving persons to file: %v", err)
	}

	return nil
}

// removePerson removes a person from the array and saves it to the file
func removePerson(p Person) error {
	// Load existing persons
	persons, err := loadPersonsFromFile()
	if err != nil {
		return fmt.Errorf("error loading persons from file: %v", err)
	}

	// Find the person in the array
	index := -1
	for i, person := range persons {
		if person.FirstName == p.FirstName {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("person with name %s not found", p.FirstName)
	}

	// Remove the person from the array
	persons = append(persons[:index], persons[index+1:]...)

	// Save the updated array back to the file
	if err := savePersonsToFile(persons); err != nil {
		return fmt.Errorf("error saving persons to file: %v", err)
	}

	return nil
}

// updatePerson updates a person in the array and saves it to the file
func updatePerson(p Person) error {
	// Load existing persons
	persons, err := loadPersonsFromFile()
	if err != nil {
		return fmt.Errorf("error loading persons from file: %v", err)
	}

	// Find the person in the array
	index := -1
	for i, person := range persons {
		if person.FirstName == p.FirstName {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("person with name %s not found", p.FirstName)
	}

	// Update the person in the array
	persons[index] = p

	// Save the updated array back to the file
	if err := savePersonsToFile(persons); err != nil {
		return fmt.Errorf("error saving persons to file: %v", err)
	}

	return nil
}

// getPersonByName retrieves a person by name from the array
func getPersonByNameFromFile(name string) (*Person, error) {
	// Load existing persons
	persons, err := loadPersonsFromFile()
	if err != nil {
		return nil, fmt.Errorf("error loading persons from file: %v", err)
	}

	// Find the person in the array
	for _, person := range persons {
		if person.FirstName == name {
			return &person, nil
		}
	}

	return nil, fmt.Errorf("person with name %s not found", name)
}

// ------------------------------------------ Buildings ------------------------------------------¨

// loadBuildingsFromFile loads the array of buildings from the file
func loadBuildingsFromFile() ([]Building, error) {
	file, err := os.Open("building.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty array if the file does not exist
			return []Building{}, nil
		}
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var buildings []Building
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&buildings); err != nil {
		return nil, fmt.Errorf("error decoding JSON from file: %v", err)
	}

	return buildings, nil
}

// saveBuildingsToFile saves the array of buildings to the file
func saveBuildingsToFile(buildings []Building) error {
	file, err := os.Create("building.json")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON
	if err := encoder.Encode(buildings); err != nil {
		return fmt.Errorf("error encoding JSON to file: %v", err)
	}

	return nil
}

// addBuilding adds a new building to the array and saves it to the file
func addBuilding(b Building) error {
	// Load existing buildings
	buildings, err := loadBuildingsFromFile()
	if err != nil {
		return fmt.Errorf("error loading buildings from file: %v", err)
	}

	// Append the new building
	buildings = append(buildings, b)

	// Save the updated array back to the file
	if err := saveBuildingsToFile(buildings); err != nil {
		return fmt.Errorf("error saving buildings to file: %v", err)
	}

	return nil
}

// removeBuilding removes a building from the array and saves it to the file
func removeBuilding(b Building) error {
	// Load existing buildings
	buildings, err := loadBuildingsFromFile()
	if err != nil {
		return fmt.Errorf("error loading buildings from file: %v", err)
	}

	// Find the building in the array
	index := -1
	for i, building := range buildings {
		if building.Type == b.Type {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("building with title %s not found", b.Type)
	}

	// Remove the building from the array
	buildings = append(buildings[:index], buildings[index+1:]...)

	// Save the updated array back to the file
	if err := saveBuildingsToFile(buildings); err != nil {
		return fmt.Errorf("error saving buildings to file: %v", err)
	}

	return nil
}

// updateBuilding updates a building in the array and saves it to the file
func updateBuilding(b Building) error {
	// Load existing buildings
	buildings, err := loadBuildingsFromFile()
	if err != nil {
		return fmt.Errorf("error loading buildings from file: %v", err)
	}

	// Find the building in the array
	index := -1
	for i, building := range buildings {
		if building.Type == b.Type {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("building with title %s not found", b.Type)
	}

	// Update the building in the array
	buildings[index] = b

	// Save the updated array back to the file
	if err := saveBuildingsToFile(buildings); err != nil {
		return fmt.Errorf("error saving buildings to file: %v", err)
	}

	return nil
}

// ------------------------------------------ World ------------------------------------------

// loadWorldFromFile loads the world from the file
func loadWorldFromFile() (World, error) {
	file, err := os.Open("world.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty world if the file does not exist
			return World{}, nil
		}
		return World{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var world World
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&world); err != nil {
		return World{}, fmt.Errorf("error decoding JSON from file: %v", err)
	}

	return world, nil
}

// saveWorldToFile saves the world to the file
func saveWorldToFile(world World) error {
	file, err := os.Create("world.json")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON
	if err := encoder.Encode(world); err != nil {
		return fmt.Errorf("error encoding JSON to file: %v", err)
	}

	return nil
}

// updateWorld updates the world and saves it to the file
func updateWorld(world World) error {
	// Save the updated world back to the file
	if err := saveWorldToFile(world); err != nil {
		return fmt.Errorf("error saving world to file: %v", err)
	}

	return nil
}
