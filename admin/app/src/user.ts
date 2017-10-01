import { GetJWT, GetToken } from './utils/token';

export function IsLoggedIn(): boolean {
  const tkn = GetToken();
  if (tkn === '') {
    return true;
  }
  return false;
}

export function ID(): string {
  const jwt = GetJWT(GetToken());
  return jwt.id;
}

export function Email(): string {
  const jwt = GetJWT(GetToken());
  return jwt.email;
}

export function FirstName(): string {
  const jwt = GetJWT(GetToken());
  return jwt.first_name;
}

export function LastName(): string {
  const jwt = GetJWT(GetToken());
  return jwt.last_name;
}

export function PhotoURL(): string {
  const jwt = GetJWT(GetToken());
  return jwt.photo_url;
}

export function HasCreditCard(): boolean {
  const jwt = GetJWT(GetToken());
  return getKthBit(jwt.perm, 0);
}

function getKthBit(x: number, k: number) {
  return (((x >> k) & 1) === 1);
}
