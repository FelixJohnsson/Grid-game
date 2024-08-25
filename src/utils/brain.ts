import Person from "./Person";
import { faker } from '@faker-js/faker';
import * as T from '../api/types';

class Brain {
    person: Person;
    personality: T.Personalities;
    hunger: number;
    energy: number;
    mood: [string, number][];
    health: number;
    taskList: T.TaskDetails[]
    onMind: string;
    memories: T.Memories[];

    constructor(person: Person){
        this.person = person;
        this.personality = faker.helpers.arrayElement(Object.values(T.Personalities));
        this.hunger = 0;
        this.energy = 100;
        this.mood = [[T.Moods.Happy, 100], [T.Moods.Sad, 0], [T.Moods.Angry, 0], [T.Moods.Excited, 0], [T.Moods.Bored, 0], [T.Moods.Tired, 0], [T.Moods.Stressed, 0], [T.Moods.Relaxed, 0]];
        this.taskList  = [T.Tasks.Idle];
        this.health = 100;
        this.onMind = '';
        this.memories = [];
    }

    turnOn(){
        console.log('Brain is on for ' + this.person.name);
        setInterval(() => {
            if (this.taskList.length > 0) {
                const currentTask = this.taskList[0];
                console.log('Current task is ' + currentTask.name);
                switch (currentTask.name) {
                    case T.Tasks.Idle.name:
                        this.think();
                        break;
                    case T.Tasks.Work.name:
                        this.work();
                        break;
                    case T.Tasks.Eat.name:
                        this.eat();
                        break;
                    case T.Tasks.Sleep.name:
                        this.sleep();
                        break;
                    case T.Tasks.Move.name:
                        this.move();
                        break;
                    case T.Tasks.Talk.name:
                        this.talk();
                        break;
                    case T.Tasks.Think.name:
                        this.think();
                        break;
                    default:
                        console.log('No task to perform');
                        break;
                }
            }
        }, 3000);
    }

    turnOff(){
        console.log('Brain is off for ' + this.person.name);
    }

    think(){
        console.log('Thinking...');
    }

    talk(){
        console.log('Talking...');
    }

    decideToMove(){
        this.person.isMoving = true;
        const up = { x: this.person.location.x, y: this.person.location.y - 1 };
        const down = { x: this.person.location.x, y: this.person.location.y + 1 };
        const left = { x: this.person.location.x - 1, y: this.person.location.y };
        const right = { x: this.person.location.x + 1, y: this.person.location.y };

        const directions = [up, down, left, right];
        const randomDirection = faker.helpers.arrayElement(directions);

        const isItPossibleToMoveThere = this.person.currentWorldState.isItPossibleToMoveTo(randomDirection.x, randomDirection.y);

        if (isItPossibleToMoveThere) {
            console.log(`${this.person.name} is moving to ${randomDirection.x}, ${randomDirection.y}`);
            this.person.currentWorldState.movePerson(this.person, this.person.location.x, this.person.location.y, randomDirection.x, randomDirection.y);
        }
    }

    addTask = (newTask: T.TaskDetails) => {
        if (this.taskList.length === 0) {
            this.taskList = [newTask];
            console.log('Task list was empty, added new task.');
        } else {
            if (this.taskList.some(t => t.name === newTask.name)) {
                console.log('Task already in the list');
                return;
            } else {
                if (this.taskList[0].name === T.Tasks.Sleep.name) {
                    console.log('Task added to the list, replacing Sleep');
                    this.taskList[0] = newTask; // Replaces Sleep but keeps other tasks
                    this.wakeUp();
                    return;
                }
                if (this.taskList[0].name === T.Tasks.Idle.name) {
                    console.log('Task added to the list, replacing Idle');
                    this.taskList[0] = newTask; // Replaces Idle but keeps other tasks
                    return;
                }
                if (newTask.priority > this.taskList[0].priority) {
                    console.log('Task added to the list with higher priority'); 
                    this.taskList.unshift(newTask);
                } else {
                    console.log('Task added to the list with lower priority');
                    this.taskList.push(newTask);
                }
            }
        }

    }

    move(){
        console.log('Moving...');
    }

    work(){
        if (this.person.isWorking) {
            console.log('Already working');
            return;
        }
        if (this.person.isChild) {
            console.log('Children do not work');
            return;
        }
        if (this.person.isWorkingAt === null) {
            console.log('No workplace');
            return;
        }
        if (this.person.occupation === 'Unemployed') {
            console.log('Unemployed');
            return;
        }
        if (this.energy < 20) {
            console.log('Not enough energy');
            return;
        }

        console.log('Working...');
        if(this.personality === T.Personalities.HardWorker){
            this.onMind = 'Working hard';
        }

        // Replace with a switch statement

    }

    eat(){
        console.log('Eating...');
    }

    sleep(){
        console.log('Sleeping...');
        this.dream();
    }

    dream(){
        const randomDream = faker.helpers.arrayElement(['Flying', 'Falling', 'Being chased', 'Being late', 'Being naked', 'Being lost', 'Being unprepared']);
        this.onMind = 'Dreaming about ' + randomDream;
    }

    wakeUp(){
        console.log('Waking up...');
        this.onMind = '';
    }

    equip(){
        console.log('Equipping...');
    }
}

export default Brain;