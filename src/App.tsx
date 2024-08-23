import { useState } from 'react';
import './App.css';
import { buildingData } from './utils/buildingData';
import { BuildingDetails, Resources, PlacedBuilding } from './types';
import BuildingSidebar from './components/Sidebar';
import BuildingWidget from './components/BuildingWidget';
import Map from './components/Map';

const InformationTopBar = ({ resources }: { resources: Resources }) => {
  return (
    <div className="bg-gray-800 text-white p-4 flex justify-center">
      <div className="mr-8">Inhabitants: {resources.inhabitants.length}</div>
      <div className="mr-8">Wood: {resources.wood}</div>
      <div className="mr-8">Food: {resources.food}</div>
      <div className="mr-8">Stone: {resources.stone}</div>
      <div className="mr-8">Ores: {resources.ores}</div>
      <div className="mr-8">Money: {resources.money}</div>
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
    inhabitants: [],
    wood: 100,
    food: 100,
    stone: 100,
    ores: 50,
    money: 1000,
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
    setIsWidgetOpen(false);
  };

  const handlePlaceDown = () => {
    console.warn('Place down');
    setIsPlacing(true);
    setSelectedBuilding(selectedBuilding);
    setIsWidgetOpen(false);
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
      setSelectedBuilding(null);
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
        inhabitants={resources.inhabitants}
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
