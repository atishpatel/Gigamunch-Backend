declare let ui: SwaggerUI;

// Set Environment
let APP = APP || {};
APP.IsDev = false;
APP.IsStage = false;
APP.IsProd = false;
switch (location.hostname) {
    case '127.0.0.1':
    case 'localhost':
        APP.IsDev = true;
        break;
    case 'gigamunch-omninexus-dev.appspot.com':
        APP.IsStage = true;
        break;
    default:
        APP.IsProd = true;
}

setTimeout(() => {
    if (ui) {
        console.log("set auth-token");
        APP.Auth.GetToken().then((token: string) => {
            ui.preauthorizeApiKey("auth-token", token)
        })
    }
}, 3000);

interface SwaggerUI {
    preauthorizeApiKey(key: string, value: string | null): void;
}
