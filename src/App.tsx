import './App.css';
import cloneDeep from 'lodash/cloneDeep';
import LocalWorldState, { WorldState } from './utils/worldState';
import Person from './utils/Person';
import Building from './utils/Building';
import { useEffect, useRef, useState } from 'react';

function App() {
  const didRun = useRef(false);

  const [persons, setPersons] = useState<Person[]>([]);
  const [buildings, setBuildings] = useState<Building[]>([]);
  const [map, setMap] = useState<WorldState>(LocalWorldState);

  useEffect(() => {
    if (didRun.current) return;
    didRun.current = true;

    const person = new Person(2, 2);
    LocalWorldState.addPerson(person);

    const lumberjack = new Building('lumberjack', { x: 0, y: 2 });
    LocalWorldState.addBuilding(lumberjack);

    setPersons([...LocalWorldState.getPersons()]);
    setBuildings([...LocalWorldState.getBuildings()]);

    LocalWorldState.updateMapState();

    setInterval(() => {
      setMap(cloneDeep(LocalWorldState));
  }, 1000);
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
        {
          map.getMap().map((row, y) => (
            <div key={y} className="flex">
              {
                row.map((gridItem, x) => {
                  return (
                  <div key={x} className={`w-10 h-10 border ${gridItem.isGround ? 'bg-green-500' : 'bg-white'}`}>
                    {
                      gridItem.inhabitants.map((person) => {
                        return (
                        <div key={person.name} className="text-lg text-black h-full bg-slate-100">{person.initials}</div>
                      )})
                    }
                    {
                      gridItem.building ? (
                        <div className={`text-lg text-black h-full ${gridItem.building.color}`}>{gridItem.building.icon}</div>
                      ) : null
                    }
                  </div>
                )})
              }
            </div>
          ))
        }
    </div>
  );
}

export default App;
