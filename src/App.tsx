import './App.css';
import Map from './components/Map';
import LocalWorldState from './utils/worldState';
import Person from './utils/Person';
import Building from './utils/Building';
import { useEffect, useRef, useState } from 'react';

function App() {
  const didRun = useRef(false);

  const [persons, setPersons] = useState<Person[]>([]);
  const [buildings, setBuildings] = useState<Building[]>([]);

  useEffect(() => {
    if (didRun.current) return;
    didRun.current = true;

    const person = new Person(6, 6);
    LocalWorldState.addPerson(person);

    const lumberjack = new Building('lumberjack', { x: 8, y: 8 });
    LocalWorldState.addBuilding(lumberjack);

    setPersons([...LocalWorldState.getPersons()]);
    setBuildings([...LocalWorldState.getBuildings()]);
    
  }, []);

  return (
    <div className="App">
        {
          persons.map((person) => (
            <div key={person.name}>
              {person.name + ' at ' + person.location.x + ', ' + person.location.y}
            </div>
          ))        
        }
        {
          buildings.map((building) => (
            <div key={building.title}>
              {building.title + ' at ' + building.location.x + ', ' + building.location.y}
            </div>
          ))
        }
      <Map />

    </div>
  );
}

export default App;
