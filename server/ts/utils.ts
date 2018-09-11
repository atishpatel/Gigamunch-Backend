

export function GetURLParmas() {
    let vars: any = {};
    window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, (m: string, key: string, value: string): string => {
        vars[key] = value;
        return value;
    });
    return vars;
}