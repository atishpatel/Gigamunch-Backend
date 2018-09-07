import * as utils from './utils';
import * as service from './service';
import * as auth from './auth';

declare var APP: any;

APP.Utils = utils;
APP.Service = service;
// APP.User = user;
APP.Auth = auth;
console.log('app.js loaded');
