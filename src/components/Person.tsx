import * as T from "../api/types";
import { useState } from "react";
import PersonWidget from "./PersonWidget";

const Person = ({
  person,
  currentTile,
  grab,
  onMouseEnter,
  onMouseLeave,
}: {
  person: T.Person;
  currentTile: T.Tile;
  grab: (item: T.Item, person: T.Person) => void;
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
      {showPersonWidget && (
        <PersonWidget
          person={person}
          currentTile={currentTile}
          grab={grab}
          togglePersonWidget={togglePersonWidget}
        />
      )}
    </div>
  );
};

export default Person;
