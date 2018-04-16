import { GetJWT, GetToken, } from './token';
addEventListener('UserUpdated', UpdateUser);
export function IsLoggedIn() {
    const tkn = GetToken();
    if (tkn === '') {
        return false;
    }
    return true;
}
export let ID = '';
export let Email = '';
export let FirstName = '';
export let LastName = '';
export let PhotoURL = '';
export function UpdateUser() {
    const tkn = GetToken();
    if (!tkn) {
        return;
    }
    const jwt = GetJWT(tkn);
    if (!jwt) {
        return;
    }
    ID = jwt.id;
    Email = jwt.email;
    FirstName = jwt.first_name;
    LastName = jwt.last_name;
    PhotoURL = jwt.photo_url;
}
export function IsAdmin() {
    const jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 2);
}
export function HasCreditCard() {
    const jwt = GetJWT(GetToken());
    if (!jwt) {
        return false;
    }
    return getKthBit(jwt.perm, 0);
}
function getKthBit(x, k) {
    return (((x >> k) & 1) === 1);
}
UpdateUser();
