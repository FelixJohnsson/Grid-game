import Person from './utils/Person';
import Building from './utils/Building';

export enum Jobs {
  Farmer = 'Farmer',
  Miner = 'Miner',
  Lumberjack = 'Lumberjack',
  Builder = 'Builder',
  Soldier = 'Soldier',
  Unemployed = 'Unemployed',
}

enum ItemEffect {
  Strength = 'Strength',
  Intelligence = 'Intelligence',
  Charisma = 'Charisma',
  Speed = 'Speed',
  Health = 'Health',
  Happiness = 'Happiness',
  Hunger = 'Hunger',
  Thirst = 'Thirst',
  Energy = 'Energy',
}

export enum BuildingType {
  House = 'House',
  Lumberjack = 'Lumberjack',
  Mine = 'Mine',
  Farm = 'Farm',
  Barracks = 'Barracks',
}

export enum BuildingState {
  Idle = 'Idle',
  Working = 'Working',
  LackingResources = 'Lacking Resources',
  LackingWorkers = 'Lacking Workers',
}

export enum BuildingDescription {
  House = 'The House provides your inhabitants with shelter.',
  Lumberjack = 'The Lumberjack gathers wood for your village.',
  Mine = 'The Mine extracts valuable minerals.',
  Farm = 'The Farm produces food for your villagers.',
  Barracks = 'The Barracks trains your military units.',
}

export enum BuildingIcon {
  House = 'H',
  Lumberjack = 'L',
  Mine = 'M',
  Farm = 'F',
  Barracks = 'B',
}

export enum BuildingColor {
  House = 'bg-blue-500',
  Lumberjack = 'bg-orange-900',
  Mine = 'bg-gray-500',
  Farm = 'bg-yellow-500',
  Barracks = 'bg-red-500',
}

export interface Item {
  name: string;
  icon: string;
  cost: { wood: number; food: number; stone: number; money: number };
  effect: [ItemEffect, number];
}

export interface Resources {
  inhabitants: Person[];
  wood: number;
  food: number;
  stone: number;
  ores: number;
  money: number;
}

export interface GridItem {
  building: Building | null;
  inhabitants: Person[];
  isGround: boolean;
  isBuilding: boolean;
  isWater: boolean;
  isRoad: boolean;
}