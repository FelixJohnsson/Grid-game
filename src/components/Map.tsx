import { useEffect, useState } from 'react';
import LocalWorldState from '../utils/worldState'; 
import { GridItem } from '../types';

const Map = () => {
  const [localMap, setLocalMap] = useState<GridItem[][]>();
  
  useEffect(() => {
    setInterval(() => {
    setLocalMap(LocalWorldState.getMap());
    }, 1000);
  }, []);

  return (
    <div>
      {
        localMap?.map((row, y) => (
          <div key={y} className="flex">
            {
              row.map((gridItem, x) => {
                return (
                <div key={x} className={`w-10 h-10 border ${gridItem.isGround ? 'bg-green-500' : 'bg-white'}`}>
                  {
                    gridItem.inhabitants.map((person) => {
                      return (
                      <div key={person.name} className="text-lg text-black bg-slate-100">{person.initials}</div>
                    )})
                  }
                </div>
              )})
            }
          </div>
        ))
      }
    </div>
  );
};

export default Map;
