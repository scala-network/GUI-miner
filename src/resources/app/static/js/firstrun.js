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

			// get the computer name
			astilectron.sendMessage({name: "firstrun", payload: ""}, function(message) {
				$('#username').html(message.payload);
			});

			// get the actual miner path
			astilectron.sendMessage({name: "get-miner-path", payload: ""}, function(message) {
				$('#miner_path').html(message.payload);
			});

			// return the pool list for the GUI miner
			astilectron.sendMessage({name: "pool-list", payload: ""}, function(message) {
				$('#table-inner').html(message.payload);
				// select the first element
				$('div.table-body').first().addClass('selected');
			});
		});
	},
	listen: function() {
		// handle error messages
		astilectron.onMessage(function(message) {
			var parsed = $.parseJSON(message.payload);
			switch (message.name) {
				case "fatal_error":
					shared.showError(parsed.data);
			}
		});
	},
	bindEvents: function() {
		$(document).on('click', '[data-target]', function() {
			var id = $(this).data('target');
			if (id == 'select-pool') {
				// make the pool list table scrollable
				if (!$("#table-inner").hasClass('mCSB_container'))  {
					$("#table-inner").mCustomScrollbar({
						theme:"rounded-dots",
						scrollInertia:400
					});

					// enable the selected ticker
					$('div.table-body').off('click').on('click', function() {
						$(this).parent().find('.table-body').removeClass('selected');
						$(this).addClass('selected');
					});
				}
			}
			if (id == 'configuring-miner') {
				// send the configuration to Go backend,
				// then wait for Go's OK to continue
				var configData = {
					address: $('#miner-address-input').val(),
					pool: $('#select-pool').find('.table-body.selected').data('id')
				};
				astilectron.sendMessage({name: "configure", payload: configData}, function(message) {
					document.location = 'index.html';
				});
			}
		});

		// miner address validation
		$('#miner-address-next-step').on('click', function(event) {
			event.preventDefault();
			asticode.notifier.close();
			var address = $('#miner-address-input').val();
			if (address == '') {
				asticode.notifier.error('You must enter your address');
			} else if (!shared.validateWalletAddress(address)) {
				asticode.notifier.error("Please enter a valid Bloc address starting with 'abLoc'");
			} else {
				$(this).next().trigger('click');
			}
		});
	}
};