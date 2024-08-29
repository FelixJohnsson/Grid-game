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
  Farmer = "Farmer",
  Miner = "Miner",
  Lumberjack = "Lumberjack",
  Builder = "Builder",
  Soldier = "Soldier",
  Unemployed = "Unemployed",
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
  House = "bg-blue-500",
  Lumberjack = "bg-orange-900",
  Mine = "bg-gray-500",
  Farm = "bg-yellow-500",
  Barracks = "bg-red-500",
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
  Introvert = "Introvert",
  Extrovert = "Extrovert",
  Ambivert = "Ambivert",
  HardWorker = "Hard Worker",
  Lazy = "Lazy",
  Fighter = "Fighter",
  Peacemaker = "Peacemaker",
  Leader = "Leader",
  Follower = "Follower",
  Diplomat = "Diplomat",
  Aggressive = "Aggressive",
  PeoplePleaser = "People Pleaser",
}

export enum Moods {
  Happy = "Happy",
  Sad = "Sad",
  Angry = "Angry",
  Excited = "Excited",
  Bored = "Bored",
  Tired = "Tired",
  Stressed = "Stressed",
  Relaxed = "Relaxed",
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
  taskList: TaskDetails[];
  onMind: string;
  memories: Memories[];
}

export const Tasks = {
  Idle: { name: "Idle", priority: 0, duration: 0, energyExpense: 0 },
  Work: { name: "Work", priority: 3, duration: 5, energyExpense: 10 },
  Eat: { name: "Eat", priority: 2, duration: 1, energyExpense: -5 },
  Sleep: { name: "Sleep", priority: 1, duration: 8, energyExpense: -15 },
  Move: { name: "Move", priority: 2, duration: 2, energyExpense: 5 },
  Talk: { name: "Talk", priority: 2, duration: 3, energyExpense: 5 },
  Think: { name: "Think", priority: 4, duration: 4, energyExpense: 8 },
} as const;

export type TaskType = keyof typeof Tasks;
export type Task = (typeof Tasks)[TaskType];

export type Relationship = {
  WithPerson: string;
  Relationship: string;
  Intensity: number;
};

export interface Person {
  FullName: string;
  Age: number;
  Title: string;
  Location: Location;
  IsTalking: boolean;
  Thinking: string;
  RightHand?: Item[];
  LeftHand?: Item[];
  Relationhips: Relationship[];
}

// --------------------- Items ---------------------

export interface Material {
  name: string;
  type: string;
  hardness: number;
  weight: number;
  density: number;
  malleability: number;
}

export interface Residue {
  name: string;
  amount: number;
}

export interface Item {
  Name: string;
  Sharpness: number;
  Bluntness: number;
  Weight: number;
  Material: Material[];
  Residues: Residue[];
}

// --------------------- Plants --------------------------
export interface PlantAction {
  Name: string;
  Target?: Tile; // Optional, as it can be null or undefined
  Priority: number;
}

export interface PlantLife {
  active: boolean;
  ctx: any; // Replace `any` with the specific type for context, if known
  cancel: () => void;
  actions: PlantAction[];
}

export interface Nutrients {
  Calories: number;
  Carbs: number;
  Protein: number;
  Fat: number;
  Vitamins: number;
  Minerals: number;
}

export interface Fruit {
  Name: string;
  Taste: string;
  Age: number;
  RipeAge: number;
  IsRipe: boolean;
  Nutrients: Nutrients[];
}

export enum PlantStage {
  Seed = "Seed",
  Sprout = "Sprout",
  Vegetative = "Vegetative",
  Flowering = "Flowering",
  Fruiting = "Fruiting",
}

export interface Plant {
  Name: string;
  Age: number;
  Health: number;
  IsAlive: boolean;
  ProducesFruit: boolean;
  Fruit: Fruit[];
  PlantStage: PlantStage;
  PlantLife?: PlantLife;
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
  Type: TileType;
  Building?: Building;
  Persons?: Person[];
  Items?: Item[];
  Plants?: Plant[];
}

// World represents a 2D array of tiles.
export interface World {
  tiles: Tile[][];
}
