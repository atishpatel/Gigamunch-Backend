
export function IsError(resp: any): boolean {
    if (resp.code && resp.code !== 0 && resp.code !== 200) {
        return true;
    }
    if (resp.error && resp.error.code && resp.error.code !== 0 && resp.error.code !== 200) {
        return true;
    }
    return false;
}