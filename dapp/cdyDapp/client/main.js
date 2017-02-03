import { Template } from 'meteor/templating';
import { ReactiveVar } from 'meteor/reactive-var';
import { Tasks } from './tasks.js';

//import '../imports/api/tasks.js';

import './main.html';
//import './body.html';

Template.hello.onCreated(function helloOnCreated() {
  // counter starts at 0
  this.counter = new ReactiveVar(0);
  //this.counter = EthBlocks.latest.number;
  //EthBlocks.init();
  //this.counter = EthBlocks.latest.number;
});

Template.hello.helpers({
  counter() {
    return Template.instance().counter.get();
  },
  //currentBlock: function(){
  //  return EthBlocks.latest.number;
  //},
/*  taskss: [
      { text: 'This is task 1' },
      { text: 'This is task 2' },
      { text: 'This is task 3' },
    ],*/

});

Template.body.helpers({
  tasks2: [
    { text: 'This is task 1' },
    { text: 'This is task 2' },
    { text: 'This is task 3' },
  ],
 tasks() {
      return Tasks.find({});//[{ text: 'This is task 1' },];//T
    },
});


Template.body.events({
  'submit .new-task'(event) {
    // Prevent default browser form submit
    event.preventDefault();

    // Get value from form element
    const target = event.target;
    const text = target.text.value;

    // Insert a task into the collection
    Tasks.insert({
      text,
      createdAt: new Date(), // current time
    });

    // Clear form
    target.text.value = '';
  },
});


Template.hello.events({
  'click button'(event, instance) {
    // increment the counter when button is clicked
    instance.counter.set(instance.counter.get() + 1);
  },
});
