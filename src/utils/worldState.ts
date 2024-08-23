import { GridItem, BuildingDetails } from '../types';
import Person from './Person';
import Map from './MapState'

// This world state class will be used for the singular map of one player world state
// This will in the future be just a part of the world state, which will be a collection of all player world states
// This will probably have to be written in C# or some other language that can handle multiple players fast
// But for now, we will use TypeScript to simulate the world state

export class WorldState {
    private gridSize: number;
    public map: GridItem[][];
    private persons: Person[];    
    private buildings: BuildingDetails[];

  constructor(gridSize: number) {
    this.gridSize = gridSize;
    this.map = new Map(this.gridSize).getMap();
    this.persons = [];
    this.buildings = [];
  }

    getMap(): GridItem[][] {
        return this.map;
    }

    createNewPerson(x: number, y: number): Person {
        const person = new Person(x, y);
        this.addPerson(person);
        return person;
    }

    getPersons(): Person[] {
        return this.persons;
    }
    
    addPerson(person: Person): void {
        this.persons.push(person);
        this.updateMap(person.location.x, person.location.y, { inhabitants: [person] });
    }

    isItPossibleToMoveTo(x: number, y: number): boolean {
        console.log('Checking if it is possible to move to ' + x + ', ' + y);
        if (x < 0 || y < 0 || x >= this.gridSize || y >= this.gridSize) {
            console.log('Cannot move to ' + x + ', ' + y + ' because it is out of bounds');
            return false;
        }
        console.log('Checking if it is possible to move to ' + x + ', ' + y + ' and the result is ' + this.map[y][x].isGround || this.map[y][x].isRoad || this.map[y][x].isBuilding);
        return this.map[y][x].isGround || this.map[y][x].isRoad || this.map[y][x].isBuilding;
    }

    movePerson(person: Person, currentX: number, currentY: number, targetX: number, targetY: number): void {
        const possible = this.isItPossibleToMoveTo(targetX, targetY);

        if (!possible) {
            throw new Error('Cannot move to ' + targetX + ', ' + targetY);
        } else {
            this.map[currentY][currentX].inhabitants = this.map[currentY][currentX].inhabitants.filter((p) => p !== person);
            this.map[targetY][targetX].inhabitants.push(person);
            person.location = { x: targetX, y: targetY };
        }
    }

    getBuildings(): BuildingDetails[] {
        return this.buildings;
    }


    addBuilding(building: BuildingDetails): void {
        this.buildings.push(building);
    }

    updateMap(x: number, y: number, newItem: Partial<GridItem>): void {
        this.map[y][x] = { ...this.map[y][x], ...newItem };
    }

    doesGridItemExist(x: number, y: number): boolean {
        return x >= 0 && y >= 0 && x < this.map.length && y < this.map[0].length;
    }

    getGridItem(x: number, y: number): GridItem {
        if (!this.doesGridItemExist(x, y)) {
            throw new Error('Invalid grid item at ' + x + ', ' + y);
        }
        return this.map[y][x];
    }

    doesGridItemHaveBuilding(x: number, y: number): boolean {
        return this.map[y][x].building !== null;
    }

    doesGridItemHavePerson(x: number, y: number): boolean {
        return this.map[y][x].inhabitants.length > 0;
    }
}

const LocalWorldState = new WorldState(50);

export default LocalWorldState;