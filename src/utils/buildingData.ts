import { BuildingDetails } from '../types';

export const buildingData: Record<string, BuildingDetails> = {
  House: {
    title: 'House',
    description: 'The House provides your inhabitants with shelter.',
    icon: 'H',
    cost: { wood: 25, food: 10, stone: 10, money: 50 },
    workers: [],
  },
  Lumberjack: {
    title: 'Lumberjack',
    description: 'The Lumberjack gathers wood for your village.',
    icon: 'L',
    cost: { wood: 20, food: 0, stone: 10, money: 100 },
    workers: [],
  },
  Mine: {
    title: 'Mine',
    description: 'The Mine extracts valuable minerals.',
    icon: 'M',
    cost: { wood: 20, food: 0, stone: 5, money: 200 },
    workers: [],
  },
  Farm: {
    title: 'Farm',
    description: 'The Farm produces food for your villagers.',
    icon: 'F',
    cost: { wood: 5, food: 0, stone: 0, money: 50 },
    workers: [],
  },
  Barracks: {
    title: 'Barracks',
    description: 'The Barracks trains your military units.',
    icon: 'B',
    cost: { wood: 50, food: 25, stone: 50, money: 500 },
    workers: [],
  },
};
