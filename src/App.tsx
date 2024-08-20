import { useState } from 'react';
import './App.css';
import { buildingData, BuildingDetails } from './utils/buildingData';
import BuildingSidebar from './components/Sidebar';
import BuildingWidget from './components/BuildingWidget';

interface PlacedBuilding {
  building: BuildingDetails;
  position: { x: number; y: number };
}

interface Resources {
  wood: number;
  food: number;
  stone: number;
  ores: number;
}

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
}) => {
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

const InformationTopBar = ({ resources }: { resources: Resources }) => {
  return (
    <div className="bg-gray-800 text-white p-4 flex justify-center">
      <div className="mr-8">Wood: {resources.wood}</div>
      <div className="mr-8">Food: {resources.food}</div>
      <div className="mr-8">Stone: {resources.stone}</div>
      <div className="mr-8">Ores: {resources.ores}</div>
    </div>
  );
};

function App() {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const [selectedBuilding, setSelectedBuilding] =
    useState<BuildingDetails | null>(null);
  const [placedBuildings, setPlacedBuildings] = useState<PlacedBuilding[]>([]);
  const [isPlacing, setIsPlacing] = useState(false);
  const [isWidgetOpen, setIsWidgetOpen] = useState(false);

  const [resources, setResources] = useState<Resources>({
    wood: 100,
    food: 100,
    stone: 100,
    ores: 50,
  });

  const cellSize = 50;
  const gridSize = Math.floor(window.innerWidth / cellSize);

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  const handleBuildingClick = (building: string) => {
    setSelectedBuilding(buildingData[building]);
    setIsWidgetOpen(true);
  };

  const closeWidget = () => {
    setIsPlacing(false); // Ensure placing mode is turned off
  };

  const handlePlaceDown = () => {
    console.warn('Place down');
    setIsPlacing(true);
    setSelectedBuilding(selectedBuilding); // Keep the selected building for placing
    setIsWidgetOpen(false); // Close the widget after placing
  };

  const handleCellClick = (x: number, y: number) => {
    console.warn('Cell clicked', x, y);
    console.warn('Is placing', isPlacing);
    console.warn('Selected building', selectedBuilding);
    if (isPlacing && selectedBuilding) {
      console.warn('Placing building', selectedBuilding, x, y);
      setPlacedBuildings((prev) => [
        ...prev,
        { building: selectedBuilding, position: { x, y } },
      ]);
      console.warn('Placed building', selectedBuilding, x, y);
      setIsPlacing(false);
      setSelectedBuilding(null); // Clear the selected building after placement
    }
  };

  return (
    <div className="App">
      <InformationTopBar resources={resources} />
      <Map
        gridSize={gridSize}
        cellSize={cellSize}
        placedBuildings={placedBuildings}
        onCellClick={handleCellClick}
        isPlacing={isPlacing}
      />
      <BuildingSidebar
        isOpen={isSidebarOpen}
        toggleSidebar={toggleSidebar}
        onBuildingClick={handleBuildingClick}
      />
      {isWidgetOpen && (
        <BuildingWidget
          building={selectedBuilding}
          onClose={closeWidget}
          onPlaceDown={handlePlaceDown}
        />
      )}
    </div>
  );
}

export default App;
