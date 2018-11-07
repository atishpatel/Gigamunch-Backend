export function IsDev(): boolean {
    if (location.port !== '8080' && location.port !== '8081') {
        return false;
    }
    return true;
}

export function IsStage(): boolean {
    if (location.hostname !== 'gigamunch-omninexus-dev.appspot.com') {
        return false;
    }
    return true;
}

export function IsProd(): boolean {
    if (location.hostname === 'gigamunch-omninexus-dev.appspot.com' || location.hostname !== 'localhost') {
        return false;
    }
    return true;
}
