import { GridItem, BuildingDetails } from "../types";
import Person from './Person'

class Map {
  private map: GridItem[][];

  constructor(gridSize: number) {
    if (this.checkIfJSONMapExists()) {
      const savedMap = this.getLocalStorageMap();
      if (savedMap) {
        this.map = savedMap
      } else {
        this.map = this.createMap(gridSize);
        this.saveToLocalStorage(this.map);
      }
    } else {
      this.map = this.createMap(gridSize);
      this.saveToLocalStorage(this.map);
    }
  }

  private createMap(gridSize: number): GridItem[][] {
    const initialGridItem: GridItem = {
      building: null,
      inhabitants: [],
      isGround: true,
      isBuilding: false,
      isWater: false,
      isRoad: false,
    };

    // Creates a 2D array filled with copies of the initial grid item
    return Array.from({ length: gridSize }, () =>
      Array.from({ length: gridSize }, () => ({ ...initialGridItem }))
    );
  }

  getMap(): GridItem[][] {
    return this.map;
  }

  updateMap(x: number, y: number, newItem: Partial<GridItem>): void {
    this.map[y][x] = { ...this.map[y][x], ...newItem };
  }

  getGridItem(x: number, y: number): GridItem {
    if (!this.doesGridItemExist(x, y)) {
      throw new Error('Invalid grid item at ' + x + ', ' + y);
    }
    return this.map[y][x];
  }

  doesGridItemExist(x: number, y: number): boolean {
    return x >= 0 && y >= 0 && x < this.map.length && y < this.map[0].length;
  }

  doesGridItemHaveBuilding(x: number, y: number): boolean {
    return this.map[y][x].building !== null;
  }

  createNewPerson(): Person {
    return new Person(6, 6);
  }

  addPersonToGrid(x: number, y: number, person: Person): void {
    this.map[y][x].inhabitants.push(person);

    console.warn(`${y} ${x}`, this.map[y][x])
  }

  addBuildingToGrid(x: number, y: number, building: BuildingDetails): void {
    const isOccupied = this.doesGridItemHaveBuilding(x, y);
    if (!isOccupied) {
          this.map[y][x].building = building;
    } else {
      throw new Error('Cannot build on occupied grid item');
    }
  }

  saveToLocalStorage(map: GridItem[][]): void {
    localStorage.setItem('worldMap', JSON.stringify(map));
  }

  checkIfJSONMapExists(): boolean {
    return localStorage.getItem('worldMap') !== null;
  }

  getLocalStorageMap(): GridItem[][] {
    const savedMap = localStorage.getItem('worldMap');
    if (savedMap) {
      return JSON.parse(savedMap);
    } else {
      throw new Error('No saved map found');
    }
  }
}

export default Map;