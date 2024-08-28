import * as T from "../api/types";
import { useState } from "react";

const Person = ({
  person,
  currentTile,
  grab,
}: {
  person: T.Person;
  currentTile: T.Tile;
  grab: (item: T.Item, person: T.Person) => void;
}) => {
  const [showWidget, setShowWidget] = useState(false);

  const PersonWidget = (person: T.Person) => {
    return (
      <div className="h-36 w-36 absolute bg-slate-400 z-50 cursor-default">
        {person.FullName}
        <div>{person.IsTalking ? "Is talking" : "Is not talking"}</div>
        {currentTile.items?.length && currentTile.items.length > 0 && (
          <div>
            <p>Items:</p>
            {currentTile.items.map((item) => (
              <div className="flex" key={item.Name}>
                <p>{item.Name}</p>
                <button
                  className="border p-2"
                  onClick={() => {
                    grab(item, person);
                  }}
                >
                  Grab
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  };

  const toggleWidget = () => {
    setShowWidget(!showWidget);
  };

  return (
    <div
      key={person.FullName}
      className="bg-orange-500 rounded-3xl w-10/12 h-10/12 flex justify-center items-center cursor-pointer"
      onClick={toggleWidget}
    >
      {person.FullName[0]}
      {showWidget && <PersonWidget {...person} />}
    </div>
  );
};

export default Person;
