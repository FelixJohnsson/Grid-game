export interface BuildingDetails {
  title: string;
  description: string;
  icon: string;
}

export const buildingData: Record<string, BuildingDetails> = {
  Lumberjack: {
    title: 'Lumberjack',
    description: 'The Lumberjack gathers wood for your village.',
    icon: 'L',
  },
  Mine: {
    title: 'Mine',
    description: 'The Mine extracts valuable minerals.',
    icon: 'M',
  },
  Farm: {
    title: 'Farm',
    description: 'The Farm produces food for your villagers.',
    icon: 'F',
  },
  Barracks: {
    title: 'Barracks',
    description: 'The Barracks trains your military units.',
    icon: 'B',
  },
};
