declare let ui: SwaggerUI;
declare let APP: any;

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
