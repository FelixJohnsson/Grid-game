export type WorldResponse = {
  message: CleanedTile[][];
  status: number;
};

export interface Location {
  X: number;
  Y: number;
}

export type CognitiveMapEntity = {
  FullName: string;
  SpeciesType: string;
  IsAlive: boolean;
};

export type CognitiveMapPlant = {
  Name: string;
  IsAlive: boolean;
  ProducesFruit: boolean;
  PlantStage: number;
  FruitCount?: number;
  HasRipeFruit?: boolean;
};

export type CognitiveMapKnownTile = {
  Location: Location;
  TileType: TileType;
  Entity?: CognitiveMapEntity;
  Plant?: CognitiveMapPlant;
  Items?: Record<string, number>;
  LastSeenUnixMs?: number;
};

export type CognitiveMapResponse = {
  message: CognitiveMapKnownTile[];
  status: number;
};

export type Material = {
  Name: string;
  Type: string;
  Hardness: number;
  Weight: number;
  Density: number;
  Malleability: number;
};

export type Item = {
  Name: string;
  Sharpness: number;
  Bluntness: number;
  Weight: number;
  Material: Material[];
  Residues: Residue[];
  Location: Location;
};
class Fruit {
  Name: string;
  RipeAge: number;
  IsRipe: boolean;
  NutritionalValue: number;

  constructor(
    name: string,
    ripeAge: number,
    isRipe: boolean,
    nutritionalValue: number
  ) {
    this.Name = name;
    this.RipeAge = ripeAge;
    this.IsRipe = isRipe;
    this.NutritionalValue = nutritionalValue;
  }

  getName(): string {
    return this.Name;
  }

  getNutritionalValue(): number {
    return this.NutritionalValue;
  }
}

export type CleanedTile = {
  Type: TileType;
  // Backend currently returns Entity. Keep Person as fallback for compatibility.
  Entity?: PersonCleaned;
  Person?: PersonCleaned;
  Items?: Item[];
  Plant?: PlantCleaned;
  Shelter?: Shelter;
};

export enum TileType {
  Grass,
  Water,
  Sand,
  Rock,
}

export type LimbThatCantGrab = {
  LimbStatus: LimbStatus;
};

export type LimbThatCanGrab = {
  LimbStatus: LimbStatus;
  Items: Item[];
  WeightOfItems: number;
};

export type LimbThatCanMove = {
  LimbStatus: LimbStatus;
};

export type Arm = {
  LimbThatCantGrab: LimbThatCantGrab;
  Hand: LimbThatCanGrab;
};

export type Leg = {
  LimbThatCanMove: LimbThatCanMove;
  Foot: LimbThatCanMove;
};

type BodyPartType = string;

export type TargetedAction = {
  Action: string;
  Target: string;
  IsActive: boolean;
  RequiresLimb: BodyPartType[];
  Priority: number;
};

export type PersonCleaned = {
  FirstName: string;
  FamilyName: string;
  FullName: string;
  Gender: string;
  Age: number;
  Title: string;

  Location: Location;

  Thinking: string;
  LastDialogue?: string;
  DialogueWith?: string;
  DialogueAtMs?: number;
  DialogueMode?: "Said" | "Heard";

  Head: HeadCleaned;
  Torso?: LimbStatus;
  RightArm?: Arm;
  LeftArm?: Arm;
  RightLeg?: Leg;
  LeftLeg?: Leg;

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
  Personalities?: Personality[];

  CurrentTask: TargetedAction;
  BaseClaimed?: boolean;
  BaseLocation?: Location;
  BaseStockpile?: Record<string, number>;
};

enum PlantStage {
  Seed,
  Sprout,
  Vegetative,
  Flowering,
  Fruiting,
}

export type PlantCleaned = {
  Name: string;
  Age: number;
  Health: number;
  IsAlive: boolean;
  ProducesFruit: boolean;
  Fruit: Fruit[];
  PlantStage: PlantStage;
};

export type Relationship = {
  WithPerson: string;
  Relationship: string;
  Intensity: number;
};

export type Personality = "Friendly" | "Hostile";

export type HeadCleaned = {
  LimbStatus: LimbStatus;
};

export type LimbStatus = {
  BluntDamage: number;
  SharpDamage: number;
  IsBleeding: boolean;
  IsBroken: boolean;
  Residues: Residue[];
  CoveredWith: Wearable[];
  IsAttached: boolean;
};

export type Residue = {
  Name: string;
  Amount: number;
};

export type Wearable = {
  Name: string;
  Material: string;
  Protection: number;
};

export type Shelter = {
  Owner: PersonCleaned;
  Location: Location;
  Inhabitants: PersonCleaned[];
  Inventory: Item[];
};
