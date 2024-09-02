import * as T from "../api/types";

type Props = {
  person: T.PersonCleaned;
  togglePersonWidget: () => void;
  currentTile: T.CleanedTile;
  grab: (item: T.Item, person: T.PersonCleaned) => void;
};

const PersonWidget = ({
  person,
  togglePersonWidget,
  currentTile,
  grab,
}: Props) => {
  return (
    <div
      className="h-36 w-36 fixed bg-blue-200 z-50 cursor-default top-100 left-100"
      onClick={togglePersonWidget}
    >
      <div>
        <button className="border p-2" onClick={togglePersonWidget}>
          X
        </button>
      </div>
      <div>
        <p>{person.FullName}</p>
        <p>Age: {person.Age}</p>
      </div>
    </div>
  );
};

export default PersonWidget;
