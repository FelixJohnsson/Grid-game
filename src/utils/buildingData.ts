import { BuildingDetails, BuildingType } from '../types';
import Person from './Person';

export const buildingData: Record<string, BuildingDetails> = {
  house: {
    title: BuildingType.House,
    description: 'The House provides your inhabitants with shelter.',
    icon: 'H',
    cost: { wood: 25, food: 10, stone: 10, money: 50 },
    workers: [],
  },
  lumberjack: {
    title: BuildingType.Lumberjack,
    description: 'The Lumberjack gathers wood for your village.',
    icon: 'L',
    cost: { wood: 20, food: 0, stone: 10, money: 100 },
    workers: [],
  },
  mine: {
    title: BuildingType.Mine,
    description: 'The Mine extracts valuable minerals.',
    icon: 'M',
    cost: { wood: 20, food: 0, stone: 5, money: 200 },
    workers: [],
  },
  farm: {
    title: BuildingType.Farm,
    description: 'The Farm produces food for your villagers.',
    icon: 'F',
    cost: { wood: 5, food: 0, stone: 0, money: 50 },
    workers: [],
  },
  barracks: {
    title: BuildingType.Barracks,
    description: 'The Barracks trains your military units.',
    icon: 'B',
    cost: { wood: 50, food: 25, stone: 50, money: 500 },
    workers: [],
  },
};

class Building {
  title: string;
  description: string;
  icon: string;
  cost: { wood: number; food: number; stone: number; money: number };
  workers: Person[];
  state = 'idle';

  constructor(title: string) {
    const building = buildingData[title];
    this.title = building.title;
    this.description = building.description;
    this.icon = building.icon;
    this.cost = building.cost;
    this.workers = building.workers;
  }
}

export default Building;