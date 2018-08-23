const url = new URL(window.location.href);
const authTkn = url.searchParams.get("auth-token");

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