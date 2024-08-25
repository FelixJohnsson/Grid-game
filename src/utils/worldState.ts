import Person from './Person';
import Building from './Building';
import MapState from './MapState'

// This world state class will be used for the singular map of one player world state
// This will in the future be just a part of the world state, which will be a collection of all player world states
// This will probably have to be written in C# or some other language that can handle multiple players fast
// But for now, we will use TypeScript to simulate the world state

export class WorldState {
    private gridSize: number;
    public mapState: MapState;
    private persons: Person[];    
    private buildings: Building[];

    constructor(gridSize: number) {
        this.gridSize = gridSize;
        this.mapState = new MapState(this.gridSize);
        this.persons = [];
        this.buildings = [];
    }

    updateMapState(): void {

    }

    getMap(): GridItem[][] {
        return this.mapState.getMap();
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
        if (x < 0 || y < 0 || x >= this.gridSize || y >= this.gridSize) {
            return false;
        }
        return this.mapState.getMap()[y][x].isGround || this.mapState.getMap()[y][x].isRoad || this.mapState.getMap()[y][x].isBuilding;
    }

    movePerson(person: Person, currentX: number, currentY: number, targetX: number, targetY: number): void {
        // We assume that a check that it's possible to move has been made
        console.warn('The brain has decided to move ' + person.name + ' from ' + currentX + ', ' + currentY + ' to ' + targetX + ', ' + targetY);
        // Remove from current location
        this.updateMap(currentX, currentY, { inhabitants: [] }); // We should actually remove the person from the inhabitants array

        // Add to new location
        this.updateMap(targetX, targetY, { inhabitants: [person] });

        // Update person's location
        person.location.x = targetX;
        person.location.y = targetY;

        // Update person's location in the persons array
        const index = this.persons.findIndex(p => p.name === person.name);
        this.persons[index] = person;

        // Get the coordinates of the person
        console.log('Person ' + person.name + ' is now at ' + person.location.x + ', ' + person.location.y);
    }

    getBuildings(): Building[] {
        return this.buildings;
    }

    addBuilding(building: Building): void {
        this.buildings.push(building);
        this.updateMap(building.location.x, building.location.y, { building });
    }

    updateMap(x: number, y: number, newItem: Partial<GridItem>): void {
        this.mapState.getMap()[y][x] = { ...this.mapState.getMap()[y][x], ...newItem };
    }

    doesGridItemExist(x: number, y: number): boolean {
        return x >= 0 && y >= 0 && x < this.mapState.getMap().length && y < this.mapState.getMap()[0].length;
    }

    getGridItem(x: number, y: number): GridItem {
        if (!this.doesGridItemExist(x, y)) {
            throw new Error('Invalid grid item at ' + x + ', ' + y);
        }
        return this.mapState.getMap()[y][x];
    }

    doesGridItemHaveBuilding(x: number, y: number): boolean {
        return this.mapState.getMap()[y][x].building !== null;
    }

    doesGridItemHavePerson(x: number, y: number): boolean {
        return this.mapState.getMap()[y][x].inhabitants.length > 0;
    }
}

const LocalWorldState = new WorldState(20);

export default LocalWorldState;