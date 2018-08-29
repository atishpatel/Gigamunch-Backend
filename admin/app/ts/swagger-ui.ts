const authTkn = GetToken();

declare let ui: SwaggerUI;

setTimeout(() => {
    if (ui) {
        console.log("set auth-token");
        ui.preauthorizeApiKey("auth-token", authTkn)
    }
}, 3000);

interface SwaggerUI {
    preauthorizeApiKey(key: string, value: string | null): void;
}

function GetToken(): string {
    const name = 'AUTHTKN=';
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length).replace(/\n/g, '');
        }
    }
    if (location.hostname === 'localhost') {
        const tnk = window.localStorage.getItem('AUTHTKN');
        if (!tnk) {
            return '';
        }
        return tnk;
    }
    return '';
}