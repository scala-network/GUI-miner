/*
  shared.js contains functions used by both firstrun and app
 */

let shared = {
  // showError takes an error message and display's it using the
  // bundled modal
  showError: function(message) {
    alert("Testing " + val);
    let errDiv = document.createElement("div");
    errDiv.innerHTML = parsed.data;
    $('.astimodaler-body').addClass('error');
    asticode.modaler.setContent(errDiv);
    asticode.modaler.show();
  },
  // validateWalletAddress checks if the given address is a valid Stellite
  // wallet address
  validateWalletAddress: function(address) {
    /*
      The regular expression to match the address
      ^(Se)\d[0-9a-zA-Z]{94}$/
      was taken from the Bisq pull request
      https://github.com/bisq-network/bisq-desktop/pull/1307/commits/2b2773e666417b179cc07edc19ede4eba4aa4ab6#diff-7e18464877c4444f041e934dc88a6b3bR437
    */
    return /^(Se)\d[0-9a-zA-Z]{94}$/.test(address);
  },
  // bindExternalLinks ensures external links are opened outside of Electron
  bindExternalLinks: function() {
    var shell = require('electron').shell;
    $(document).on('click', 'a[href^="http"]', function(event) {
      event.preventDefault();
      shell.openExternal(this.href);
    });
    // This stops electron from updating the window title when a link
    // is clicked
    $(document).on('click', 'a[href^="#"]', function(event) {
      event.preventDefault();
    });
  },
  isMac: function() {
    return window.navigator.platform.toLowerCase().includes("mac");
  }
}
