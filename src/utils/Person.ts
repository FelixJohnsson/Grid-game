import { faker } from '@faker-js/faker';
import { Jobs, Item, Task } from '../types';
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

  addTask = (task: Task) => {
    this.brain.addTask(task);
  }

  addEmployer = (building: Building) => {
    if (this.isChild) return;
    this.isWorkingAt = building;
    building.workers.push(this);
    this.occupation = building.title === 'Lumberjack' ? Jobs.Lumberjack : building.title === 'Mine' ? Jobs.Miner : building.title === 'Farm' ? Jobs.Farmer : Jobs.Unemployed;
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
