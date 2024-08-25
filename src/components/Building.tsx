import * as T from "../api/types";

const Building = ({ building }: { building: T.Building }) => {
  return (
    <div
      key={building.Name}
      className={`flex justify-center items-center ${building.Color}`}
    >
      {building.Icon}
    </div>
  );
};

export default Building;
