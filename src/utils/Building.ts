import { BuildingType, BuildingState, BuildingDescription, BuildingIcon, BuildingColor } from '../types';
import Person from './Person';

export const buildingData: Record<string, Building> = {
  house: {
    name: '',
    title: BuildingType.House,
    description: BuildingDescription.House,
    icon: BuildingIcon.House,
    cost: { wood: 25, food: 10, stone: 10, money: 50 },
    workers: [],
    location: { x: 0, y: 0 },
    state: BuildingState.Idle,
    color: BuildingColor.House,
  },
  lumberjack: {
    name: '',
    title: BuildingType.Lumberjack,
    description: BuildingDescription.Lumberjack,
    icon: BuildingIcon.Lumberjack,
    cost: { wood: 20, food: 0, stone: 10, money: 100 },
    workers: [],
    location: { x: 0, y: 0 },
    state: BuildingState.Idle,
    color: BuildingColor.Lumberjack,
  },
  mine: {
    name: '',
    title: BuildingType.Mine,
    description: BuildingDescription.Mine,
    icon: BuildingIcon.Mine,
    cost: { wood: 20, food: 0, stone: 5, money: 200 },
    workers: [],
    location: { x: 0, y: 0 },
    state: BuildingState.Idle,
    color: BuildingColor.Mine,
  },
  farm: {
    name: '',
    title: BuildingType.Farm,
    description: BuildingDescription.Farm,
    icon: BuildingIcon.Farm,
    cost: { wood: 5, food: 0, stone: 0, money: 50 },
    workers: [],
    location: { x: 0, y: 0 },
    state: BuildingState.Idle,
    color: BuildingColor.Farm,
  },
  barracks: {
    name: '',
    title: BuildingType.Barracks,
    description: BuildingDescription.Barracks,
    icon: BuildingIcon.Barracks,
    cost: { wood: 50, food: 25, stone: 50, money: 500 },
    workers: [],
    location: { x: 0, y: 0 },
    state: BuildingState.Idle,
    color: BuildingColor.Barracks,
  },
};

class Building {
  name: string;
  title: BuildingType;
  description: BuildingDescription;
  icon: BuildingIcon;
  cost: { wood: number; food: number; stone: number; money: number };
  workers: Person[];
  state = BuildingState.Idle;
  location: { x: number, y: number };
  color: string;

  constructor(title: string, name: string, location: { x: number, y: number }) {
    this.name = name;
    const building = buildingData[title];
    this.title = building.title;
    this.description = building.description;
    this.icon = building.icon;
    this.cost = building.cost;
    this.workers = building.workers;
    this.location = location;
    this.color = building.color;
  }
}

export default Building;