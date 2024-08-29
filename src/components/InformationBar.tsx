import * as T from "../api/types";

type Props = {
  persons: T.PersonCleaned[] | undefined;
};

const InformationBar = ({ persons }: Props) => {
  return (
    <div className="flex flex-col items-center">
      <div>
        <h1 className="text-lg underline">Persons</h1>
        <div>
          {persons
            ? persons.map((person, i) => (
                <div key={i}>
                  <p className="bg-orange-500 text-sm">{person.FullName}</p>
                </div>
              ))
            : null}
        </div>
      </div>
    </div>
  );
};

export default InformationBar;
