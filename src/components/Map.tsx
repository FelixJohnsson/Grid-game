import { GridItem } from '../types';

const Map = ({mapState}:{ mapState:  GridItem[][]}) => {
  return (
    <div>
      {
        mapState.map((row, y) => (
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
};

export default Map;
