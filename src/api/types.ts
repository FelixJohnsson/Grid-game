// ----------------- World ------------------
export enum TileType {
  Grass,
  Water,
  Mountain,
}

type Tile = {
  Type: TileType;
  Persons?: Person[];
  Items?: Item[];
  Plant?: Plant;
  NutritionalValue?: number;
};

export type World = {
  Tiles: Tile[][];
};

// ----------------- People -----------------

enum Jobs {
  Farmer = "Farmer",
  Miner = "Miner",
  Lumberjack = "Lumberjack",
  Builder = "Builder",
  Soldier = "Soldier",
  Unemployed = "Unemployed",
}

type Relationship = {
  WithPerson: string;
  Relationship: string;
  Intensity: number;
};

export type Person = {
  Age: number;
  Title: string;
  FirstName: string;
  FamilyName: string;
  FullName: string;
  Initials: string;
  IsChild: boolean;
  Gender: string;
  Description: string;
  Occupation: Jobs;
  SkinColor: string;
  Personality: string;
  Genes: string[];

  IsMoving: TargetedAction;
  IsTalking: TargetedAction;
  IsSitting: TargetedAction;
  IsEating: TargetedAction;
  IsSleeping: TargetedAction;

  Thinking: string;
  WantsTo: string;
  FeelingSafe: number;
  FeelingScared: number;

  Body?: HumanBody;

  Strength: number;
  Agility: number;
  Intelligence: number;
  Charisma: number;
  Stamina: number;

  CombatExperience: number;
  CombatSkill: number;
  CombatStyle: string;

  Relationships: Relationship[];

  IsIncapacitated: boolean;
  VisionRange: number;
  Location: Location;
  OnTileType: TileType;
};

// ----------------- Brain ------------------

type TargetedAction = {
  Action: string;
  Target: string;
  IsActive: boolean;
  RequiresLimb: string[];
};

type RequestedAction = TargetedAction & {
  From: Person;
};

type Brain = {
  Owner?: Person;
  Active: boolean;
  Ctx: any; // Placeholder for context.Context
  Cancel: any; // Placeholder for context.CancelFunc
  Actions: TargetedAction[];
  IsConscious: boolean;
  IsAlive: boolean;
  BrainDamage: number;
};

type Vision = {
  Buildings: BuildingCleaned[];
  Persons: PersonInVision[];
};

// ----------------- Body -------------------

type LimbStatus = {
  BluntDamage: number;
  SharpDamage: number;
  IsBleeding: boolean;
  IsBroken: boolean;
  Residues: Residue[];
  CoveredWith: Wearable[];
  IsAttached: boolean;
};

enum LimbType {
  RightHand = "RightHand",
  LeftHand = "LeftHand",
  RightFoot = "RightFoot",
  LeftFoot = "LeftFoot",
  RightLeg = "RightLeg",
  LeftLeg = "LeftLeg",
  TheHead = "Head",
  Torso = "Torso",
}

type LimbThatCanHold = LimbStatus & {
  Items: Item[];
  WeightOfItems: number;
};

type Damage = {
  AmountBluntDamage: number;
  AmountSharpDamage: number;
};

type Head = LimbStatus & {
  Brain?: Brain;
};

type LimbThatCanGrab = LimbStatus & {
  Items?: Item[];
  WeightOfItems: number;
};

type LimbThatCantGrab = LimbStatus;

type LimbThatCanMove = LimbStatus;

type Leg = LimbThatCanMove & {
  Foot?: LimbThatCanMove;
};

type Arm = LimbThatCantGrab & {
  Hand?: LimbThatCanGrab;
};

type HumanBody = {
  Head?: Head;
  Torso?: LimbStatus;
  RightArm?: Arm;
  LeftArm?: Arm;
  RightLeg?: Leg;
  LeftLeg?: Leg;
};

// ----------------- Plants -----------------

type PlantAction = {
  Name: string;
  Target?: Tile;
  Priority: number;
};

type PlantLife = {
  active: boolean;
  ctx: any; // Placeholder for context.Context
  cancel: any; // Placeholder for context.CancelFunc
  actions: PlantAction[];
};

type Nutrients = {
  Calories: number;
  Carbs: number;
  Protein: number;
  Fat: number;

  Vitamins: number;
  Minerals: number;
};

type Fruit = {
  Name: string;
  Taste: string;
  Age: number;
  RipeAge: number;
  IsRipe: boolean;
  Nutrients: Nutrients[];
};

enum PlantStage {
  Seed,
  Sprout,
  Vegetative,
  Flowering,
  Fruiting,
}

type Plant = {
  Name: string;
  Age: number;
  Health: number;
  IsAlive: boolean;
  ProducesFruit: boolean;
  Fruit: Fruit[];
  PlantStage: PlantStage;
  PlantLife?: PlantLife;
};

// ----------------- Items ------------------

type Wearable = {
  Name: string;
  Material: string;
  Protection: number;
};

type Material = {
  Name: string;
  Type: string;
  Hardness: number;
  Weight: number;
  Density: number;
  Malleability: number;
};

type Residue = {
  Name: string;
  Amount: number;
};

export type Item = {
  Name: string;
  Sharpness: number;
  Bluntness: number;
  Weight: number;
  Material: Material[];
  Residues: Residue[];
};

// ----------------- Cleaned ------------------

export type CleanedTile = {
  Type: TileType;
  Building?: BuildingCleaned;
  Persons?: PersonCleaned[];
  Items?: Item[];
  Plant?: PlantCleaned;
};

type PlantCleaned = {
  Name: string;
  Age: number;
  Health: number;
  IsAlive: boolean;
  ProducesFruit: boolean;
  Fruit: Fruit[];
  PlantStage: PlantStage;
};

type BuildingCleaned = {
  Name: string;
  Type: string;
  Location: Location;
};

type HeadCleaned = LimbStatus;

export type PersonCleaned = {
  FirstName: string;
  FamilyName: string;
  FullName: string;
  Gender: string;
  Age: number;
  Title: string;

  Location: Location;

  Thinking: string;

  Head: HeadCleaned;
  Torso: LimbStatus;
  RightArm: Arm;
  LeftArm: Arm;
  RightLeg: Leg;
  LeftLeg: Leg;

  Strength: number;
  Agility: number;
  Intelligence: number;
  Charisma: number;
  Stamina: number;

  CombatExperience: number;
  CombatSkill: number;
  CombatStyle: string;

  IsIncapacitated: boolean;

  Relationships: Relationship[];
};

type PersonInVision = {
  FullName: string;
  FirstName: string;
  FamilyName: string;
  Gender: string;
  Age: number;
  Title: string;
  Location: Location;
  Body?: HumanBody;
};

// ----------------- Server response ------------------
export type WorldResponse = {
  message: CleanedTile[][];
  status: number;
};

export type PersonResponse = {
  message: Person[];
  status: number;
};
