import {
  GetJWT,
  GetToken,
} from './token';

addEventListener('UserUpdated', UpdateUser);

export function IsLoggedIn(): boolean {
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
  const jwt = GetJWT(GetToken());
  ID = jwt.id;
  Email = jwt.email;
  FirstName = jwt.first_name;
  LastName = jwt.last_name;
  PhotoURL = jwt.photo_url;
}

export function HasCreditCard(): boolean {
  const jwt = GetJWT(GetToken());
  return getKthBit(jwt.perm, 0);
}

function getKthBit(x: number, k: number) {
  return (((x >> k) & 1) === 1);
}

UpdateUser();
