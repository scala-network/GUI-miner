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
	validateWalletAddress: function(address, coin_type) {
		if (coin_type == 'bloc') {
			// 1. 99 chars for standard address
			// 2. 187 chars for address with integrated payment id
			// ex: abLoc5jeufY8yWkZgjDJnP6DuuhyGE3jb5F6kmKKqqynhbUDgfvvC2FjdP5DjjnoW2R9aecMDETTbdMuFNFzHRWvGNkzHGKHMT9
			var regexp = /^abLoc([a-zA-Z0-9]{94}|[a-zA-Z0-9]{182})$/g;
			if (address.match(regexp)) {
				return true;
			}
		} else if (coin_type == 'xmr') {
			// 1. 95 chars for address
			// 2. address always starts with 4
			// 3. second character can only be a number (0-9), or letters A or B
			// ex: 4581HhZkQHgZrZjKeCfCJxZff9E3xCgHGF25zABZz7oR71TnbbgiS7sK9jveE6Dx6uMs2LwszDuvQJgRZQotdpHt1fTdDhk
			var regexp = /^4([0-9A-B])([a-zA-Z0-9]{93})$/g;
			if (address.substring(0, 1) == '4' && address.match(regexp)) {
				return true;
			}
		} else if (coin_type == 'xhv') {
			// 1. 79 chars for address
			// 2. first 3 chars are hvx
			// ex: hvxyDX9mqBNbQ6ojRrZZYcNPSTGcxtxQ4Ws6mNm6Ag7NTciArFb71HHL8HbACGpMu3iTc42F3YQNj4r
			var regexp = /^hvx([a-zA-Z0-9]{76})$/g;
			if (address.match(regexp)) {
				return true;
			}
		} else if (coin_type == 'trtl') {
			// 1. 99 chars for address
			// 2. first 3 chars are TRTL
			// ex: TRTLv3GvjehhjeYnctiWQx6MRtgWQKURPWfocps8XuMnL9XXgF2GaYgX9vamnUcG35BkQy6VfwUy5CsV9YNomioPGGyVhM8VgAb
			var regexp = /^TRTL([a-zA-Z0-9]{95})$/g;
			if (address.match(regexp)) {
				return true;
			}
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
	},
	getCoinData: function(coin_type) {
		let data = {coin_abbr: 'BLOC', coin_prefix: 'abLoc'};
		if (coin_type == 'xmr') {
			data.coin_abbr = 'XMR';
			data.coin_prefix = '';
		} else if (coin_type == 'xhv') {
			data.coin_abbr = 'XHV';
			data.coin_prefix = 'hvx';
		} else if (coin_type == 'trtl') {
			data.coin_abbr = 'TRTL';
			data.coin_prefix = 'TRTL';
		}
		return data;
	}
}
