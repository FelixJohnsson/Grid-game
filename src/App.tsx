import './App.css';
import { useEffect, useState } from 'react';
import { getPersons, getBuildings } from './api/api';
import { Person } from './api/types';

function App() {
  const [persons, setPersons] = useState<Person[]>();
  const [buildings, setBuildings] = useState<any[]>();
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    setLoading(true);
    getPersons().then((data) => {
      setPersons(data);
    });
    getBuildings().then((data) => {
      setBuildings(data);
        setLoading(false);
    });
  }, []);

  return (
    <div className="App">
      { loading ? <p>Loading...</p> : 
        <div>
          {persons?.map((person) => {
            console.warn(person);
            return (
              <div className="pt-6" key={person.Name}>
                <p className="text-lg underline">This is your person</p>
                <p className="text-md">Name: {person.Name}</p>
                <p className="text-md">Age: {person.Age}</p>
                <p className="text-md">Is child?: {person.IsChild ? 'Yes' : 'No'}</p>
                <p className="text-md">Location: {'X: ' + person.Location.X + ', Y: ' + person.Location.Y}</p>
              </div>
              )
            })
          }
          {
            buildings?.map((building) => {
              console.warn(building);
              return (
                <div className="pt-6" key={building.Name}>
                  <p className="text-lg underline">This is your building</p>
                  <p className="text-md">Name: {building.Name}</p>
                  <p className="text-md">Type: {building.Type}</p>
                  <p className="text-md">State: {building.State}</p>
                  <p className="text-md">Location: {'X: ' + building.Location.X + ', Y: ' + building.Location.Y}</p>
                  <p className="text-md underline">Workers</p>
                  {
                    building.Workers ? building.Workers?.map((worker: any) => {
                      return (
                        <div className="pt-6" key={worker.Name}>
                          <p className="text-md">Name: {worker.Name}</p>
                          <p className="text-md">Age: {worker.Age}</p>
                          <p className="text-md">Is child?: {worker.IsChild ? 'Yes' : 'No'}</p>
                          <p className="text-md">Location: {'X: ' + worker.Location.X + ', Y: ' + worker.Location.Y}</p>
                        </div>
                      )
                    }) : <p>No workers</p>
                  }
                </div>
              )
            })
          }
        </div>
      }  
    </div>
  );
}

export default App;
