import Person from './utils/Person';

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

export interface Item {
  name: string;
  icon: string;
  cost: { wood: number; food: number; stone: number; money: number };
  effect: [ItemEffect, number];
}

export interface BuildingDetails {
  title: string;
  description: string;
  icon: string;
  cost: { wood: number; food: number; stone: number; money: number };
  workers: Person[];
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
  building: BuildingDetails | null;
  inhabitants: Person[];
  isGround: boolean;
  isBuilding: boolean;
  isWater: boolean;
  isRoad: boolean;
}