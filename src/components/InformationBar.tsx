import * as T from "../api/types";

type Props = {
  persons: T.Person[] | undefined;
  buildings: T.Building[] | undefined;
};

const InformationBar = ({ persons, buildings }: Props) => {
  return (
    <div className="flex flex-col items-center">
      <div>
        <h1 className="text-lg underline">Persons</h1>
        <div>
          {persons
            ? persons.map((person, i) => (
                <div key={i}>
                  <p className="bg-orange-500 text-sm">{person.Name}</p>
                </div>
              ))
            : null}
        </div>
      </div>

      <div>
        <h1 className="text-lg underline pt-6">Buildings</h1>
        <div>
          {buildings
            ? buildings.map((building, i) => (
                <div key={i}>
                  <p className={building.Color}>{building.Name}</p>
                </div>
              ))
            : null}
        </div>
      </div>
    </div>
  );
};

export default InformationBar;
