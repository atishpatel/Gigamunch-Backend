const url = new URL(window.location.href);
const authTkn = url.searchParams.get("auth-token");

setTimeout(()=>{
    if(window.ui){
        console.log("set auth-token");
        window.ui.preauthorizeApiKey("auth-token",authTkn)
    }
}, 3000);