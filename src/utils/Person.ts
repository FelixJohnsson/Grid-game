import { faker } from '@faker-js/faker';
import { Jobs, Item } from '../types';
import LocalWorldState, { WorldState } from './worldState';
import Brain from './brain';
import Building from './Building';

class Person {
  age: number;
  name: string;
  initials: string;
  isChild: boolean;
  gender: string;
  description: string;
  icon: string;
  occupation: Jobs;
  isWorkingAt: Building | null;
  color: string;
  location: { x: number; y: number };
  isMoving: boolean;
  isTalking: boolean;
  isSitting: boolean;
  isHolding: boolean;
  isEating: boolean;
  isSleeping: boolean;
  isWorking: boolean;
  thinking: string;
  wantsTo: string;
  inventory: Item[];
  currentWorldState: WorldState;
  genes = [];
  brain: Brain;

  constructor(x: number, y: number) {
    const age = faker.number.int({ min: 2, max: 54 });
    this.name = faker.person.fullName();
    this.initials = this.name.split(' ')[0][0] + this.name.split(' ')[1][0];
    this.age = age;
    this.isChild = age < 18;
    this.gender = faker.helpers.arrayElement(['male', 'female']);
    this.description = '';
    this.icon = 'P';
    this.occupation = Jobs.Unemployed;
    this.isWorkingAt = null;
    this.color = '';
    this.location = {
      x: x,
      y: y,
    };
    this.isMoving = false;
    this.isTalking = false;
    this.isSitting = false;
    this.isHolding = false;
    this.isEating = false;
    this.isSleeping = false;
    this.isWorking = false;
    this.thinking = '';
    this.wantsTo = '';
    this.inventory = [];
    this.currentWorldState = LocalWorldState;

    this.genes = [];
    this.brain = new Brain(this);

    console.log(`${this.name} has been created`);

    this.turnOnBrain();
  }

  turnOnBrain() {
    this.brain.turnOn();
  }

  decideToMove() {
    console.log(`${this.name} is deciding to move`);
    // Simple logic to decide if the person should move
    this.isMoving = true;
    const up = { x: this.location.x, y: this.location.y - 1 };
    const down = { x: this.location.x, y: this.location.y + 1 };
    const left = { x: this.location.x - 1, y: this.location.y };
    const right = { x: this.location.x + 1, y: this.location.y };

    const directions = [up, down, left, right];
    const randomDirection = faker.helpers.arrayElement(directions);

    const isItPossibleToMoveThere = this.currentWorldState.isItPossibleToMoveTo(randomDirection.x, randomDirection.y);

    if (isItPossibleToMoveThere) {
      this.currentWorldState.movePerson(this, this.location.x, this.location.y, randomDirection.x, randomDirection.y);
    }
  }

  startWorking() {
    if (!this.isWorking) {
      this.isWorking = true;
    }
  }

  stopWorking() {
    if (this.isWorking) {
      this.isWorking = false;
      console.log(`${this.name} has stopped working`);
    }
  }

  eat() {
    if (!this.isEating) {
      this.isEating = true;
      console.log(`${this.name} is eating`);
    }
  }

  sleep() {
    if (!this.isSleeping) {
      this.isSleeping = true;
      console.log(`${this.name} is sleeping`);
    }
  }

  wakeUp() {
    if (this.isSleeping) {
      this.isSleeping = false;
      console.log(`${this.name} has woken up`);
    }
  }

  talk() {
    if (!this.isTalking) {
      this.isTalking = true;
      console.log(`${this.name} is talking`);
    }
  }

  stopTalking() {
    if (this.isTalking) {
      this.isTalking = false;
      console.log(`${this.name} has stopped talking`);
    }
  }
}

export default Person;
