/*
  This handles the initial user setup
 */
let firstrun = {
	init: function() {
		asticode.loader.init();
		asticode.modaler.init();
		asticode.notifier.init();

		shared.bindExternalLinks();
		shared.bindTargetLinks();
		shared.extraFunctions();

		// Wait for the ready signal
		document.addEventListener('astilectron-ready', function() {
			// Get the computer name
			astilectron.sendMessage({name: "firstrun", payload: ""}, function(message) {
				$('#username').html(message.payload);
			});

			// Get the actual miner path
			astilectron.sendMessage({name: "get-miner-path", payload: ""}, function(message) {
				$('#miner_path').html(message.payload);
			});

			// Get the coins json an cache it locally
			astilectron.sendMessage({
				name: "coins-content-json",
				payload: ""
			}, function(message) {
				var parsed = $.parseJSON(message.payload);
				console.log('[' + new Date().toUTCString() + '] ', "coins-content-json", parsed);
				firstrun.coinsContent = parsed;

				firstrun.bindEvents();
				firstrun.listen();
			});
			/*
			// overide the firstrun
			var configData = {
				address: 'abLoc...',
				pool: 1,
				coin_type: 'bloc',
				coin_algo: 'cryptonight_haven'
			};
			astilectron.sendMessage({name: "configure", payload: configData}, function(message) {
				document.location = 'index.html';
			});
			*/
		});
	},
	listen: function() {
		// Handle error messages
		astilectron.onMessage(function(message) {
			var parsed = $.parseJSON(message.payload);
			switch (message.name) {
				case "fatal_error":
					shared.showError(parsed.data);
			}
		});
	},
	bindEvents: function() {
		// Functionality based on which page slide is loaded
		$(document).on('click', '[data-target]', function() {
			var id = $(this).data('target');
			if (id == 'select-pool') {
				asticode.loader.show();
				// Return the pool list for the GUI miner
				var payloadData = {
					coin_type: firstrun.coin_type
				};
				astilectron.sendMessage({name: "pool-list", payload: payloadData}, function(message) {
					$("#pool-list").mCustomScrollbar("destroy");
					$('#pool-list').html(message.payload);
					asticode.loader.hide();

					// make the pool list table scrollable
					$("#pool-list").mCustomScrollbar({
						theme:"rounded-dots",
						scrollInertia:400
					});

					// enable the selected ticker
					$('#pool-list').find('.table-body').off('click').on('click', function() {
						$(this).parent().find('.table-body').removeClass('selected');
						$(this).addClass('selected');
						firstrun.selected_pool = parseInt($(this).data('id'));
					});

					// trigger the selected pool
					$('#pool-list').find('.table-body[data-id="' + firstrun.selected_pool + '"]').trigger('click');
				});
			}
			if (id == 'configuring-miner') {
				// Send the configuration to Go backend,
				// then wait for Go's OK to continue
				var configData = {
					address: $('#miner-address-input').val().trim(),
					pool: $('#pool-list').find('.table-body.selected').data('id'),
					coin_type: firstrun.coin_type,
					coin_algo: firstrun.coin_algo
				};
				astilectron.sendMessage({name: "configure", payload: configData}, function(message) {
					document.location = 'index.html';
				});
			}
			// render the other coins that can be mined
			if (id == 'other-currencies') {
				var coins = firstrun.coinsContent.coins.filter(function(el) { // remove bloc
					return el.coin_type !== 'bloc';
				});
				coins = coins.map(function(el) { // add name and abbr keys
					el.name = firstrun.coinsContent.names[el.coin_type];
					el.abbr = firstrun.coinsContent.abbr[el.coin_type];
					return el;
				});
				var html = $.fn.tmpl("tmpl-coins-title", coins);
				$('#coins-title').html(html);
				var html = $.fn.tmpl("tmpl-coins", coins);
				$('#currencies-selector').html(html);

				// Events to mine other currencies buttons
				var cb = $('#other-currencies .currency-btn');
				cb.off('click').on('click', function(event) {
					firstrun.coin_type = $(this).data('coin_type');
					firstrun.coin_algo = $(this).data('coin_algo');
					$('#other-currencies-next-step').trigger('click');
				});
			}
			// change the miner-address content based on the selected coin
			if (id == 'miner-address') {
				// setup
				var mac = $('#miner-address-content');
				let data = {
					coin_name: firstrun.coinsContent.names[firstrun.coin_type],
					coin_abbr: firstrun.coinsContent.abbr[firstrun.coin_type],
					coin_prefix: firstrun.coinsContent.address_prefix[firstrun.coin_type]
				};

				// replace text vars
				var html = $.fn.tmpl("tmpl-miner-address-text", data);
				mac.find('.address-text').html(html);
				// replace input vars
				var html = $.fn.tmpl("tmpl-miner-address-input", data);
				mac.find('.address-input input').attr('placeholder', html);
			}
		});
		// If i've clicked "i already have a wallet" button, reset to BLOC
		$('#choose-wallet a[data-target="miner-address"]').on('click', function(event) {
			firstrun.coin_type = 'bloc';
			firstrun.coin_algo = 'cryptonight_haven';
		});

		// Pool list validation
		$('#select-pool-next-step').on('click', function(event) {
			event.preventDefault();
			asticode.notifier.close();
			if (firstrun.selected_pool == 0) {
				asticode.notifier.error('You must choose one of the pools');
			} else {
				$(this).next().trigger('click');
			}
		});

		// Miner address validation
		$('#miner-address-next-step').on('click', function(event) {
			event.preventDefault();
			asticode.notifier.close();
			var address = $('#miner-address-input').val().trim();
			if (address == '') {
				asticode.notifier.error('You must enter your address');
			} else if (!shared.validateWalletAddress(address, firstrun.coin_type)) {
				asticode.notifier.error("Please enter a valid wallet address");
			} else {
				$(this).next().trigger('click');
			}
		});
	},
	coin_type: 'bloc',
	coin_algo: 'cryptonight_haven',
	selected_pool: 0,
	coinsContent: {}
};