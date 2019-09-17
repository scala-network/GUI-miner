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
			// Get the computer username
			astilectron.sendMessage({name: "get-username", payload: ""}, function(message) {
				$('#username').html(message.payload);
			});

			// Get the actual miner path
			astilectron.sendMessage({name: "get-miner-path", payload: ""}, function(message) {
				$('#miner_path').html(message.payload);
			});

			// Get the coins json and cache it locally
			astilectron.sendMessage({
				name: "get-coins-content",
				payload: ""
			}, function(message) {
				var parsed = $.parseJSON(message.payload);
				console.log('[' + new Date().toUTCString() + '] ', "get-coins-content", parsed);
				firstrun.coinsContent = parsed;

				firstrun.bindEvents();
				firstrun.listen();
			});
		});
	},
	listen: function() {
		// Handle error messages
		astilectron.onMessage(function(message) {
			var parsed = $.parseJSON(message.payload);
			switch (message.name) {
				case "fatal_error":
					console.log('[' + new Date().toUTCString() + '] ', "fatal_error", parsed.data);
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
				astilectron.sendMessage({name: "get-pools-list", payload: payloadData}, function(message) {
					// console.log('[' + new Date().toUTCString() + '] ', "get-pools-list", message.payload);

					$("#pool-list").mCustomScrollbar("destroy");
					$('#pool-list').html(message.payload);
					asticode.loader.hide();

					if (firstrun.selected_pool == 0) {
						var fe = $('#pool-list').find('.table-body').first();
						if (fe.length > 0) {
							firstrun.selected_pool = parseInt(fe.data('id'));
						}
					}

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
					address:       $('#miner-address-input').val().trim(),
					pool:          $('#pool-list').find('.table-body.selected').data('id'),
					coin_type:     firstrun.coin_type,
					coin_algo:     firstrun.coin_algo,
					xmrig_algo:    firstrun.xmrig_algo,
					xmrig_variant: firstrun.xmrig_variant.toString(),
					hardware_type: 1 // 1 = CPU mining, 2 = GPU mining
				};
				console.log('[' + new Date().toUTCString() + '] ', "save-configuration", configData);
				astilectron.sendMessage({name: "save-configuration", payload: configData}, function() {
					document.location = 'index.html';
				});
			}
			// render the other coins that can be mined
			if (id == 'other-currencies') {
				var coins = firstrun.coinsContent.coins.filter(function(el) { // remove bloc
					return el.coin_type !== 'bloc';
				});
				coins = coins.map(function(el) { // add name and abbreviation keys
					el.name          = firstrun.coinsContent.names[el.coin_type];
					el.icon          = firstrun.coinsContent.icons[el.coin_type];
					el.abbreviation  = firstrun.coinsContent.abbreviation[el.coin_type];
					el.xmrig_algo    = firstrun.coinsContent.xmrigAlgo[el.coin_type];
					el.xmrig_variant = firstrun.coinsContent.xmrigVariant[el.coin_type];
					return el;
				});
				let html1 = $.fn.tmpl("tmpl-coins-title", coins);
				$('#coins-title').html(html1);
				let html2 = $.fn.tmpl("tmpl-coins", coins);
				$('#currencies-selector').html(html2);

				// Events to mine other currencies buttons
				$('#other-currencies .currency-btn').off('click').on('click', function() {
					let el = $(this);
					firstrun.coin_type =     el.data('coin_type');
					firstrun.coin_algo =     el.data('coin_algo');
					firstrun.xmrig_algo =    el.data('xmrig_algo');
					firstrun.xmrig_variant = el.data('xmrig_variant');
					$('#other-currencies-next-step').trigger('click');
				});
			}
			// change the miner-address content based on the selected coin
			if (id == 'miner-address') {
				// setup
				var mac = $('#miner-address-content');
				let data = {
					coin_name:   firstrun.coinsContent.names[firstrun.coin_type],
					coin_abbr:   firstrun.coinsContent.abbreviation[firstrun.coin_type],
					coin_prefix: firstrun.coinsContent.addressPrefix[firstrun.coin_type]
				};

				// replace text vars
				let html;
				html = $.fn.tmpl("tmpl-miner-address-text", data);
				mac.find('.address-text').html(html);
				// replace input vars
				html = $.fn.tmpl("tmpl-miner-address-input", data);
				mac.find('.address-input input').attr('placeholder', html);
			}
			// change to first screen
			if (id == 'choose-wallet') {
				// reset the selected pool and address
				firstrun.selected_pool = 0;
				$('#miner-address-input').val('');
			}
		});
		// If i've clicked "i already have a wallet" button
		$('#choose-wallet a[data-target="miner-address"]').on('click', function() {
			firstrun.resetToBloc();
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
			} else if (!shared.validateWalletAddress(address, firstrun.coinsContent.addressValidation, firstrun.coin_type)) {
				asticode.notifier.error("Please enter a valid wallet address");
			} else {
				$(this).next().trigger('click');
			}
		});
	},
	resetToBloc: function() {
		const bloc_key = 'bloc';
		var selected_coin = firstrun.coinsContent.coins.filter(function(el) { // remove bloc
			return el.coin_type === bloc_key;
		});
		firstrun.coin_type     = selected_coin[0].coin_type;
		firstrun.coin_algo     = selected_coin[0].coin_algo;
		firstrun.xmrig_algo    = firstrun.coinsContent.xmrigAlgo[bloc_key];
		firstrun.xmrig_variant = firstrun.coinsContent.xmrigVariant[bloc_key];
	},
	coin_type: "",
	coin_algo: "",
	xmrig_algo: "",
	xmrig_variant: "",
	selected_pool: 0,
	coinsContent: {}
};