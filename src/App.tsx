import './App.css';
import Map from './components/Map';
import LocalWorldState from './utils/worldState';
import Person from './utils/Person';
import { useEffect, useRef, useState } from 'react';

function App() {
  const didRun = useRef(false);

  const [persons, setPersons] = useState<Person[]>([]);

  useEffect(() => {
    if (didRun.current) return;
    didRun.current = true;

    const person = new Person(6, 6);
    LocalWorldState.addPerson(person);

    // Trigger a re-render by updating the state
    setPersons([...LocalWorldState.getPersons()]);
  }, []);

  console.warn('persons', persons);

  return (
    <div className="App">
      {
        persons.map((person) => (
          <div key={person.name}>
            {person.name}
          </div>
        ))
      }
      <Map />

    </div>
  );
}

export default App;
