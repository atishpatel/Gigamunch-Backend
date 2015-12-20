// TODO switch to jquery
function validateForm() {
    var email = document.getElementById("email");
    var terp = document.getElementById("terp");
    if(terp.value != ""){
        return false;
    }
    if (!verifyEmail(email.value)) {
        return false;
    }
    return true;
}

function verifyEmail(email) {
    var re = /\S+@\S+\.\S+/;
    return re.test(email);
}