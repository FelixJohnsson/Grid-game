export enum Jobs {
  Farmer = 'Farmer',
  Miner = 'Miner',
  Lumberjack = 'Lumberjack',
  Builder = 'Builder',
  Soldier = 'Soldier',
}

export interface Person {
  name: string;
  age: number;
  isChild: boolean;
  gender: string;
  description: string;
  icon: string;
  occupation: Jobs;
  color: string;
  location: { x: number; y: number };
  isMoving: boolean;
  isTalking: boolean;
  isSitting: boolean;
  isHolding: boolean;
  isEating: boolean;
  isSleeping: boolean;
  isWorking: boolean;
}

export interface BuildingDetails {
  title: string;
  description: string;
  icon: string;
  cost: { wood: number; food: number; stone: number; money: number };
}

export interface Resources {
  inhabitants: Person[];
  wood: number;
  food: number;
  stone: number;
  ores: number;
  money: number;
}

export interface PlacedBuilding {
  building: BuildingDetails;
  position: { x: number; y: number };
}
