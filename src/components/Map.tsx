import { Person, PlacedBuilding, BuildingDetails } from '../types';

const Map = ({
  gridSize,
  cellSize,
  placedBuildings,
  onCellClick,
  isPlacing,
}: {
  gridSize: number;
  cellSize: number;
  placedBuildings: PlacedBuilding[];
  onCellClick: (x: number, y: number) => void;
  isPlacing: boolean;
  inhabitants: Person[];
}) => {
  MapState(1, 2);
  const gridCells = Array(gridSize * gridSize).fill(null);

  return (
    <div>
      <div
        className="grid"
        style={{
          gridTemplateColumns: `repeat(${gridSize}, ${cellSize}px)`,
          position: 'relative',
        }}
      >
        {gridCells.map((_, index) => {
          const x = index % gridSize;
          const y = Math.floor(index / gridSize);

          const placedBuilding = placedBuildings.find(
            (building) => building.position.x === x && building.position.y === y
          );

          return (
            <div
              key={index}
              className={`bg-green-600 border border-white flex items-center justify-center cursor-pointer hover:bg-green-500 transition-colors`}
              style={{ width: cellSize, height: cellSize }}
              onClick={() => isPlacing && onCellClick(x, y)}
            >
              {placedBuilding && (
                <p className="text-center z-50">
                  {placedBuilding.building.icon}
                </p>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

// Location = x: 1
// Location = y: 1

// MapState[1][1]

interface GridItem {
  building: BuildingDetails | null;
  inhabitants: Person[];
  isGround: boolean;
  isBuilding: boolean;
  isWater: boolean;
  isRoad: boolean;
}
class WorldMap {
  map: any;
  constructor() {
    this.createMap();
  }

  createMap() {
    const gridItem: GridItem = {
      building: null,
      inhabitants: [],
      isGround: true,
      isBuilding: false,
      isWater: false,
      isRoad: false,
    };

    const MapContainer = [
      [[gridItem], [gridItem], [gridItem]],
      [[gridItem], [gridItem], [gridItem]],
      [[gridItem], [gridItem], [gridItem]],
    ];

    this.map = MapContainer;
  }

  getMap() {
    return this.map;
  }

  updateMap() {}
}

export default Map;
