import { useEffect, useState } from "react";
import { getWorld } from "./api/api";
import * as T from "./api/types";
import Map from "./components/Map";
import InformationBar from "./components/InformationBar";

function App() {
  const [persons, setPersons] = useState<T.Person[]>();
  const [buildings, setBuildings] = useState<T.Building[]>();
  const [world, setWorld] = useState<T.World>();

  useEffect(() => {
    setBuildings([]);
    setPersons([]);
    getWorld().then((data) => {
      console.warn(data);
      setWorld(data);

      // Loop through the tiles in the world and assign the persons and buildings to the state
      data.tiles.forEach((row) => {
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

  return (
    <div className="App">
      <InformationBar persons={persons} buildings={buildings} />
      <Map world={world} />
    </div>
  );
}

export default App;
