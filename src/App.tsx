import { useEffect, useState } from "react";
import * as api from "./api/api";
import * as T from "./api/types";
import Map from "./components/Map";
import InformationBar from "./components/InformationBar";
import MoveControls from "./components/MoveControls";

function App() {
  const [persons, setPersons] = useState<T.Person[]>([]);
  const [buildings, setBuildings] = useState<T.Building[]>();
  const [world, setWorld] = useState<T.World["tiles"]>();

  useEffect(() => {
    setBuildings([]);
    setPersons([]);
    api.getWorld().then((data) => {
      setWorld(data);

      // Loop through the tiles in the world and assign the persons and buildings to the state
      data.forEach((row) => {
        row.forEach((tile) => {
          if (tile.Persons) {
            for (const person of tile.Persons) {
              // Push the person to the persons state
              setPersons((persons) => {
                if (persons) {
                  return [...persons, person];
                } else {
                  return [person];
                }
              });
            }
          }
          if (tile.Building) {
            setBuildings([tile.Building]);
          }
        });
      });
    });
  }, []);

  const move = (direction: string) => {
    api.movePerson(persons[0].FullName, direction).then((data) => {
      setWorld(data);
    });
  };

  const grab = (item: T.Item, person: T.Person) => {
    console.log(person.FullName + " wants to grab " + item.Name);
    api.grabItem(item, person).then((data) => {
      setWorld(data);
    });
  };

  return (
    <div className="App">
      <div className="w-26">
        <MoveControls move={move} />
        <div>
          {
            // Display the persons
            persons.map((person) => (
              <div className="border w-40" key={person.FullName}>
                <h1 className="text-sm">{person.FullName}</h1>
                <p className="text-xs">Title: {person.Title}</p>
                <p className="text-xs">Age: {person.Age}</p>
                <p className="text-xs">
                  Thought: {person.Thinking.length > 0 ? person.Thinking : ""}
                </p>
                <p className="text-xs">
                  {person.RightHand && person.RightHand.length > 0
                    ? "Right Hand: " + person.RightHand[0].Name
                    : ""}
                </p>
              </div>
            ))
          }
        </div>
      </div>
      <InformationBar persons={persons} buildings={buildings} />
      <Map world={world} grab={grab} />
    </div>
  );
}

export default App;
