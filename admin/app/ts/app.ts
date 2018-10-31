import * as Service from './service';
import * as Auth from './auth';
import * as EventUtil from './utils/event';

declare const APP: any;

APP.Auth = Auth;
APP.Service = Service;
APP.Event = EventUtil;
console.log('app.js loaded');

function GetURLParmas() {
    let vars: any = {};
    window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function (m, key, value) {
        vars[key] = value;
        return value;
    });
    return vars;
}

GetURLParmas()
