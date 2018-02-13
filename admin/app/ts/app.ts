import * as Service from './service';
import * as User from './user';
import * as EventUtil from './utils/event';
import * as TokenUtil from './utils/token';

declare var APP: any;

APP.Service = Service;
APP.User = User;
APP.Event = EventUtil;
APP.Token = TokenUtil;
