/*
  shared.js contains functions used by both firstrun and app
 */

let shared = {
  // showError takes an error message and display's it using the
  // bundled modal
  showError: function(message) {
    let errDiv = document.createElement("div");
    errDiv.innerHTML = "<h2>Something went wrong</h2>" +
      "<p>" + message + "</p>";
    $('.astimodaler-body').removeClass('ann');
    $('.astimodaler-body').addClass('error');
    asticode.modaler.setContent(errDiv);
    asticode.modaler.show();
  },
  extraFunctions: function() {
	// rewrite notifier to keep the notification indefinitely
	asticode.notifier.notify = function(type, message) {
		const wrapper = document.createElement("div");
		wrapper.className = "astinotifier-wrapper";
		const item = document.createElement("div");
		item.className = "astinotifier-item " + type;
		const label = document.createElement("div");
		label.className = "astinotifier-label";
		label.innerHTML = message;
		const close = document.createElement("div");
		close.className = "astinotifier-close";
		close.innerHTML = `<i class="fa fa-close"></i>`;
		close.onclick = function() {
			wrapper.remove();
		};
		item.appendChild(label);
		item.appendChild(close);
		wrapper.appendChild(item);
		document.getElementById("astinotifier").prepend(wrapper);
	};
	asticode.notifier.close = function() {
		const el = document.querySelector('.astinotifier-close');
		if (el) el.click();
	};
  },
  // check if the given address is a valid BLOC wallet address
  validateWalletAddress: function(address) {
	var regexp = /^([a-z0-9]{99}|[a-z0-9]{187})$/gi;
    if (address.substring(0, 5) == 'abLoc' && address.match(regexp))
    {
      return true;
    }
    return false;
  },
  // bindExternalLinks ensures external links are opened outside of Electron
  bindExternalLinks: function() {
    var shell = require('electron').shell;
    $(document).on('click', 'a[href^="http"]', function(event) {
      event.preventDefault();
      shell.openExternal(this.href);
    });
  },
  // bindTargetLinks emulates slides with links
  bindTargetLinks: function() {
    $(document).on('click', '[data-target]', function(event) {
      event.preventDefault();
      $(this).closest('.main-section').addClass('hidden');
	  asticode.notifier.close();
      var id = $(this).data('target');
	  $('#' + id).removeClass('hidden');
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
