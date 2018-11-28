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
			firstrun.bindEvents();
			firstrun.listen();

			// Get the computer name
			astilectron.sendMessage({name: "firstrun", payload: ""}, function(message) {
				$('#username').html(message.payload);
			});

			// Get the actual miner path
			astilectron.sendMessage({name: "get-miner-path", payload: ""}, function(message) {
				$('#miner_path').html(message.payload);
			});
			/*
			// overide the firstrun
			var configData = {
				address: 'abLoc...',
				pool: 1
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
					console.log(firstrun.selected_pool);
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
			// change the miner-address content based on the selected coin
			if (id == 'miner-address') {
				// setup
				var mac = $('#miner-address-content');
				jQuery.extend({encode: function(text) {return text}}); // remove the escape made by jquery.tmpl.min.js

				// precompile
				$.template("miner_address_text", miner_address_text);
				$.template("miner_address_input", miner_address_input);

				// get coin data based on coin type
				var coin_data = shared.getCoinData(firstrun.coin_type);
				coin_data.extra = coin_data.coin_prefix != '' ? true : null;

				// replace text vars
				var mat = $.tmpl("miner_address_text", coin_data);
				mac.find('.address-text').html(mat);

				// replace input vars
				var mai = $.tmpl("miner_address_input", coin_data);
				mac.find('.address-input input').attr('placeholder', mai.text());
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

		// Events to mine other currencies buttons
		var cb = $('#other-currencies .currency-btn');
		cb.on('click', function(event) {
			firstrun.coin_type = $(this).data('coin_type');
			firstrun.coin_algo = $(this).data('coin_algo');
			$('#other-currencies-next-step').trigger('click');
		});
	},
	coin_type: 'bloc',
	coin_algo: 'cryptonight_haven',
	selected_pool: 0
};