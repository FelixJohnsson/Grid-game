import { faker } from '@faker-js/faker';
import { Person, Jobs } from '../types';

export const createNewPerson = (): Person => {
  const age = faker.number.int({ min: 2, max: 54 });
  const gender = faker.helpers.arrayElement(['male', 'female']);

  const person: Person = {
    name: faker.person.fullName(),
    age: age,
    isChild: age < 18,
    gender: gender,
    description: '',
    icon: 'P',
    occupation: Jobs.Farmer,
    color: '',
    location: {
      x: 0,
      y: 0,
    },
    isMoving: false,
    isTalking: false,
    isSitting: false,
    isHolding: false,
    isEating: false,
    isSleeping: false,
    isWorking: false,
  };

  return person;
};
