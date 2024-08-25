// ----------------------------- Buildings -----------------------------

export interface Building {
  Type: BuildingType;
  Name: string;
  State: BuildingState;
  Description: BuildingDescription;
  Icon: BuildingIcon;
  Color: BuildingColor;
}

export enum Jobs {
  Farmer = 'Farmer',
  Miner = 'Miner',
  Lumberjack = 'Lumberjack',
  Builder = 'Builder',
  Soldier = 'Soldier',
  Unemployed = 'Unemployed',
}

export enum BuildingType {
  House = "House",
  WoodCabin = "Lumberjack",
  Mine = "Mine",
  Farm = "Farm",
  Barracks = "Barracks",
}

export enum BuildingState {
  Idle = "Idle",
  Working = "Working",
  LackingResources = "Lacking Resources",
  LackingWorkers = "Lacking Workers",
}

export enum BuildingDescription {
  House = "The House provides your inhabitants with shelter.",
  Lumberjack = "The Lumberjack gathers wood for your village.",
  Mine = "The Mine extracts valuable minerals.",
  Farm = "The Farm produces food for your villagers.",
  Barracks = "The Barracks trains your military units.",
}

export enum BuildingIcon {
  House = "H",
  Lumberjack = "L",
  Mine = "M",
  Farm = "F",
  Barracks = "B",
}

export enum BuildingColor {
  House = 'bg-blue-500',
  Lumberjack = 'bg-orange-900',
  Mine = 'bg-gray-500',
  Farm = 'bg-yellow-500',
  Barracks = 'bg-red-500',
}

export interface Resources {
  inhabitants: Person[];
  wood: number;
  food: number;
  stone: number;
  ores: number;
  money: number;
}
export interface Location {
    X: number;
    Y: number;
}

// --------------------- Brain and Person ---------------------
export enum Personalities {
    Introvert = 'Introvert',
    Extrovert = 'Extrovert',
    Ambivert = 'Ambivert',
    HardWorker = 'Hard Worker',
    Lazy = 'Lazy',
    Fighter = 'Fighter',
    Peacemaker = 'Peacemaker',
    Leader = 'Leader',
    Follower = 'Follower',
    Diplomat = 'Diplomat',
    Aggressive = 'Aggressive',
    PeoplePleaser = 'People Pleaser'
}

export enum Moods {
    Happy = 'Happy',
    Sad = 'Sad',
    Angry = 'Angry',
    Excited = 'Excited',
    Bored = 'Bored',
    Tired = 'Tired',
    Stressed = 'Stressed',
    Relaxed = 'Relaxed'
}

export interface Memories {
    name: string;
    importance: number;
}

export interface TaskDetails {
    name: string;
    priority: number;
    duration: number;
    energyExpense: number;
}

export interface Brain {
    person: Person;
    personality: Personalities;
    hunger: number;
    energy: number;
    mood: Moods[];
    health: number;
    taskList: TaskDetails[]
    onMind: string;
    memories: Memories[];
}

export const Tasks = {
    Idle: { name: 'Idle', priority: 0, duration: 0, energyExpense: 0 },
    Work: { name: 'Work', priority: 3, duration: 5, energyExpense: 10 },
    Eat: { name: 'Eat', priority: 2, duration: 1, energyExpense: -5 },
    Sleep: { name: 'Sleep', priority: 1, duration: 8, energyExpense: -15 },
    Move: { name: 'Move', priority: 2, duration: 2, energyExpense: 5 },
    Talk: { name: 'Talk', priority: 2, duration: 3, energyExpense: 5 },
    Think: { name: 'Think', priority: 4, duration: 4, energyExpense: 8 },
} as const;

export type TaskType = keyof typeof Tasks;
export type Task = typeof Tasks[TaskType];


export interface Person {
    Age: number;
    Name: string;
    Initials: string;
    IsChild: boolean;
    Gender: string;
    Description: string;
    Icon: string;
    Occupation: Jobs;
    IsWorkingAt: Building | null;
    Color: string;
    Location: Location;
    IsMoving: boolean;
    IsTalking: boolean;
    IsSitting: boolean;
    IsHolding: boolean;
    IsEating: boolean;
    IsSleeping: boolean;
    IsWorking: boolean;
    Thinking: string;
    WantsTo: string;
    Inventory: Item[];
    Genes: string[];
    Brain: Brain;
}

// --------------------- Items ---------------------

export enum ItemEffect {
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


// --------------------- World State ---------------------

// Enum representing different types of terrain.
export enum TileType {
  Grass,
  Water,
  Mountain,
}

// Tile represents a single tile in the world.
export interface Tile {
  type: TileType;
  building?: Building;  // Optional Building
  persons?: Person[];   // Optional array of Persons
}

// World represents a 2D array of tiles.
export interface World {
  tiles: Tile[][];
}