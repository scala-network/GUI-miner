/*
  shared.js contains functions used by both firstrun and app
 */

// showError takes an error message and display's it using the
// bundled modal
function showError(message) {
  let errDiv = document.createElement("div");
  errDiv.innerHTML = parsed.data;
  $('.astimodaler-body').addClass('error');
  asticode.modaler.setContent(errDiv);
  asticode.modaler.show();
}

// validateWalletAddress checks if the given address is a valid Stellite
// wallet address
function validateWalletAddress(address) {
  /*
    The regular expression to match the address
    ^(Se)\d[0-9a-zA-Z]{94}$/
    was taken from the Bisq pull request
    https://github.com/bisq-network/bisq-desktop/pull/1307/commits/2b2773e666417b179cc07edc19ede4eba4aa4ab6#diff-7e18464877c4444f041e934dc88a6b3bR437
  */
 return /^(Se)\d[0-9a-zA-Z]{94}$/.test(address);
}
