

export function GetURLParmas() {
    let vars: any = {};
    let parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, (m, key: string, value: string) => {
        vars[key] = value;
    });
    return vars;
}