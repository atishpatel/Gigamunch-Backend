import * as Service from './service';
import * as User from './user';
import * as EventUtil from './utils/event';

declare var APP: any;

APP.Service = Service;
APP.User = User;
APP.Event = EventUtil;
