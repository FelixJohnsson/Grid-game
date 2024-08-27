import { useEffect, useState } from "react";
import * as api from "./api/api";
import * as T from "./api/types";
import Map from "./components/Map";
import InformationBar from "./components/InformationBar";

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
          if (tile.persons) {
            for (const person of tile.persons) {
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
          if (tile.building) {
            setBuildings([tile.building]);
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

  return (
    <div className="App">
      <div>
        <h1>Move</h1>
        <button
          className="border p-4 cursor-pointer"
          onClick={() => move("up")}
        >
          Up
        </button>
        <button
          className="border p-4 cursor-pointer"
          onClick={() => move("down")}
        >
          Down
        </button>
        <button
          className="border p-4 cursor-pointer"
          onClick={() => move("left")}
        >
          Left
        </button>
        <button
          className="border p-4 cursor-pointer"
          onClick={() => move("right")}
        >
          Right
        </button>
      </div>
      <InformationBar persons={persons} buildings={buildings} />
      <Map world={world} />
    </div>
  );
}

export default App;
