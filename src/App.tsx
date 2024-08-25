import './App.css';
import { useEffect, useState } from 'react';
import { getWorld } from './api/api';
import * as T from './api/types';

function App() {
  const [persons, setPersons] = useState<T.Person[]>();
  const [buildings, setBuildings] = useState<T.Building[]>();
  const [world, setWorld] = useState<T.World>();

  useEffect(() => {
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
      <div className="flex flex-col items-center">
      <div>
        <h1 className="text-lg underline">Persons</h1>
        <div>
          {persons ? (
            persons.map((person, i) => (
              <div key={i}>
                <p className="bg-orange-500 text-sm">{person.Name}</p>
              </div>
            ))
          ) : null}
        </div>
      </div>

      <div>
        <h1 className="text-lg underline pt-6">Buildings</h1>
        <div>
          {buildings ? (
            buildings.map((building, i) => (
              <div key={i}>
                <p className={building.Color}>{building.Name}</p>
              </div>
            ))
          ) : null}
        </div>
      </div>
      </div>
      
      <div className="w-full flex justify-center pt-6">
        {
        world ? (
          <div>
            <div>
              {world.tiles.map((row, y) => (
                <div key={y} style={{ display: 'flex' }}>
                  {row.map((tile, x) => (
                    <div
                      key={x}
                      style={{
                        width: '20px',
                        height: '20px',
                        backgroundColor: tile.type === T.TileType.Grass ? 'green' : tile.type === T.TileType.Water ? 'blue' : 'gray',
                      }}
                    >
                      {tile.building ? (
                        <div className={tile.building.Color}>{tile.building.Icon}</div>
                      ) : null}
                      {tile.persons ? (
                        <div>
                          {tile.persons.map((person, i) => (
                            <div key={i} className="bg-orange-500">{person.Initials}</div>
                          ))}
                        </div>  
                          ) : null}
                    </div>
                  ))}
                </div>
              ))}
            </div>
          </div>
        ) : null
      }
      </div>
    </div>
  );
}

export default App;
