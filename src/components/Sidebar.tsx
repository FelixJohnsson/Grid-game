import { buildingData } from '../utils/buildingData';

const BuildingSidebar = ({
  isOpen,
  toggleSidebar,
  onBuildingClick,
}: {
  isOpen: boolean;
  toggleSidebar: () => void;
  onBuildingClick: (building: string) => void;
}) => {
  return (
    <div
      className={`fixed top-0 right-28 h-full bg-gray-800 text-white transition-transform duration-300 ${
        isOpen ? 'translate-x-28' : 'translate-x-full'
      }`}
      style={{ width: '250px' }}
    >
      <button
        className="absolute top-0 left-0 m-2 bg-gray-700 text-white py-1 px-2 rounded"
        onClick={toggleSidebar}
      >
        Buildings
      </button>
      {isOpen && (
        <div className="p-4 mt-10">
          <ul>
            {Object.keys(buildingData).map((building) => (
              <li
                key={building}
                className="mb-2 cursor-pointer"
                onClick={() => onBuildingClick(building)}
              >
                {building}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default BuildingSidebar;
