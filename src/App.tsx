import { useEffect, useState } from "react";
import * as api from "./api/api";
import * as T from "./api/types";
import Map from "./components/Map";
import InformationBar from "./components/InformationBar";

function App() {
  const [persons, setPersons] = useState<T.PersonCleaned[]>();
  const [world, setWorld] = useState<T.CleanedTile[][]>();

  const continuslyUpdate = () => {
    api.getWorld().then((data) => {
      setWorld(data);

      // Update the persons
      const persons: T.PersonCleaned[] = [];
      data.forEach((row) => {
        row.forEach((tile) => {
          if (tile.Person) {
            persons.push(tile.Person);
          }
        });
      });
      setPersons(persons);
    });
  };

  useEffect(() => {
    // Fetch the initial world state
    api.getWorld().then((data) => {
      setWorld(data);
    });

    // Set up the interval to continuously update the world state
    const intervalId = setInterval(continuslyUpdate, 500);

    // Clear the interval when the component unmounts to prevent memory leaks and infinite loops
    return () => clearInterval(intervalId);
  }, []);

  const grab = (item: T.Item, person: T.PersonCleaned) => {
    console.log(person.FullName + " wants to grab " + item.Name);
    api.grabItem(item, person).then((data) => {
      setWorld(data);
    });
  };

  return (
    <div className="App">
      <div className="w-26">
        <div>
          {
            // Display the persons
            persons?.map((person) => (
              <div className="border w-40" key={person.FullName}>
                <h1 className="text-sm">{person.FullName}</h1>
                <p className="text-xs">Title: {person.Title}</p>
                <p className="text-xs">Age: {person.Age}</p>
                <p className="text-xs">
                  Thought: {person.Thinking.length > 0 ? person.Thinking : ""}
                </p>
                <p className="text-xs">
                  {person.RightArm?.Hand?.Items &&
                  person.RightArm.Hand.Items.length > 0
                    ? "Right Hand: " + person.RightArm.Hand.Items[0].Name
                    : ""}
                </p>

                <p className="text-xs">
                  Current task: {person.CurrentTask.Action}
                </p>
              </div>
            ))
          }
        </div>
      </div>
      <InformationBar persons={persons} />
      <Map world={world} grab={grab} />
    </div>
  );
}

export default App;
