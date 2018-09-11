import { UserUpdated } from './utils/event';
import { GetJWT, GetToken } from './utils/token';
addEventListener(UserUpdated, UpdateUser);
export function IsLoggedIn() {
    var tkn = GetToken();
    if (tkn === '') {
        return false;
    }
    return true;
}
export var ID = '';
export var Email = '';
export var FirstName = '';
export var LastName = '';
export var PhotoURL = '';
export var Token = '';
export function UpdateUser() {
    var tkn = GetToken();
    if (!tkn) {
        return;
    }
    var jwt = GetJWT(tkn);
    if (!jwt) {
        return;
    }
    ID = jwt.id;
    Email = jwt.email;
    FirstName = jwt.first_name;
    LastName = jwt.last_name;
    PhotoURL = jwt.photo_url;
    Token = tkn;
}
export function IsAdmin() {
    var jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 2);
}
export function HasCreditCard() {
    var jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 0);
}
function getKthBit(x, k) {
    return (((x >> k) & 1) === 1);
}
UpdateUser();
