import * as T from "../api/types";

const Person = ({ person }: { person: T.Person }) => {
  return (
    <div
      key={person.FullName}
      className="bg-orange-500 rounded-3xl w-10/12 h-10/12 flex justify-center items-center"
    >
      {person.FullName[0]}
    </div>
  );
};

export default Person;
