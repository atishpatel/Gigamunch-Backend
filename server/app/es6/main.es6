'use strict';

function verifyEmail(email) {
  const re = /\S+@\S+\.\S+/;
  return re.test(email);
}

// TODO switch to jquery
function validateForm() {
  const email = document.getElementById('email');
  const terp = document.getElementById('terp');
  if (terp.value !== '') {
    return false;
  }
  if (!verifyEmail(email.value)) {
    return false;
  }
  return true;
}
