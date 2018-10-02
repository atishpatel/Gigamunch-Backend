import * as Service from './service';
import * as Auth from './auth';
import * as EventUtil from './utils/event';

declare const APP: any;

APP.Auth = Auth;
APP.Service = Service;
APP.Event = EventUtil;
console.log('app.js loaded');
