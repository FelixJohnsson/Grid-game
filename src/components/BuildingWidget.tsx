import { BuildingDetails } from "../utils/buildingData";

const BuildingWidget = ({
  building,
  onClose,
  onPlaceDown,
}: {
  building: BuildingDetails | null;
  onClose: () => void;
  onPlaceDown: () => void;
}) => {
  if (!building) return null;

  return (
    <div
      className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50"
      onClick={onClose}
    >
      <div
        className="bg-white p-6 rounded shadow-lg w-1/3 relative"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-2xl font-bold mb-4">{building.title}</h2>
        <p className="mb-4 text-gray-700">{building.description}</p>
        { building.icon }
        <div className="flex justify-between">
          <button
            className="bg-blue-500 text-white py-1 px-4 rounded"
            onClick={onPlaceDown}
          >
            Place Down
          </button>
          <button
            className="bg-red-500 text-white py-1 px-4 rounded"
            onClick={onClose}
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default BuildingWidget;
