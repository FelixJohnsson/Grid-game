import Person from "./Person";
import { faker } from '@faker-js/faker';

enum Personalities {
    Introvert = 'Introvert',
    Extrovert = 'Extrovert',
    Ambivert = 'Ambivert',
    HardWorker = 'Hard Worker',
    Lazy = 'Lazy',
    Fighter = 'Fighter',
    Peacemaker = 'Peacemaker',
    Leader = 'Leader',
    Follower = 'Follower',
    Diplomat = 'Diplomat',
    Aggressive = 'Aggressive',
    PeoplePleaser = 'People Pleaser'
}

enum Moods {
    Happy = 'Happy',
    Sad = 'Sad',
    Angry = 'Angry',
    Excited = 'Excited',
    Bored = 'Bored',
    Tired = 'Tired',
    Stressed = 'Stressed',
    Relaxed = 'Relaxed'
}

interface Task {
    name: string;
    priority: number;
    duration: number;
    energyExpense: number;
}

class Brain {
    person: Person;
    personality: Personalities;
    hunger: number;
    energy: number;
    mood: [string, number][];
    health: number;
    taskList: Task[]

    constructor(person: Person){
        this.person = person;
        this.personality = faker.helpers.arrayElement(Object.values(Personalities));
        this.hunger = 0;
        this.energy = 100;
        this.mood = [[Moods.Happy, 100], [Moods.Sad, 0], [Moods.Angry, 0], [Moods.Excited, 0], [Moods.Bored, 0], [Moods.Tired, 0], [Moods.Stressed, 0], [Moods.Relaxed, 0]];
        this.taskList  = [{name: 'Nothing', priority: 0, duration: 0, energyExpense: 0}];
        this.health = 100;
    }

    turnOn(){
        console.log('Brain is on for ' + this.person.name);
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

    move(){
        console.log('Moving...');
    }

    work(){
        console.log('Working...');
    }

    eat(){
        console.log('Eating...');
    }

    sleep(){
        console.log('Sleeping...');
    }

    wakeUp(){
        console.log('Waking up...');
    }

    equip(){
        console.log('Equipping...');
    }
}

export default Brain;