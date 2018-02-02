import * as user from './user';
import * as service from './service';

declare var APP: any;

APP.Service = service;
APP.User = user;
console.log('app.js loaded');
