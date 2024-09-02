import * as T from "../api/types";
import { useState } from "react";

const Person = ({
  person,
  currentTile,
  grab,
  onMouseEnter,
  onMouseLeave,
}: {
  person: T.PersonCleaned;
  currentTile: T.CleanedTile;
  grab: (item: T.Item, person: T.PersonCleaned) => void;
  onMouseEnter: (event: React.MouseEvent, text: string) => void;
  onMouseLeave: () => void;
}) => {
  const [showPersonWidget, setShowPersonWidget] = useState(false);

  const togglePersonWidget = () => {
    setShowPersonWidget(!showPersonWidget);
  };

  return (
    <div
      key={person.FullName}
      className="bg-orange-500 rounded-3xl w-10/12 h-10/12 flex justify-center items-center cursor-pointer"
      onClick={togglePersonWidget}
      onMouseEnter={(e) => onMouseEnter(e, person.FullName)}
      onMouseLeave={onMouseLeave}
    >
      {person.FullName[0]}
    </div>
  );
};

export default Person;
