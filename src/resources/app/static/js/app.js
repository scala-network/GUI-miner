/*
  This handles the actual miner
 */
let app = {
	init: function() {
		asticode.loader.init();
		asticode.modaler.init();
		asticode.notifier.init();

		shared.bindExternalLinks();

		shared.bindTargetLinks();
		shared.extraFunctions();

		// Wait for the ready signal
		document.addEventListener('astilectron-ready', function() {
			// Get the coins json and cache it locally
			astilectron.sendMessage({
				name: "get-coins-content",
				payload: ""
			}, function(message) {
				var parsed = $.parseJSON(message.payload);
				console.log('[' + new Date().toUTCString() + '] ', "get-coins-content", parsed);
				app.coinsContent = parsed;

				app.bindEvents();
				app.listen();

				// Get the configuration settings
				astilectron.sendMessage({
					name: "get-config-file",
					payload: ""
				}, function(message) {
					var parsed = $.parseJSON(message.payload);
					console.log('[' + new Date().toUTCString() + '] ', "get-config-file", parsed);
					app.coin_type = parsed.coin_type;
					app.prev_coin_type = app.coin_type;
					app.coin_algo = parsed.coin_algo;
					app.xmrig_algo = parsed.xmrig_algo;
					app.xmrig_variant = parsed.xmrig_variant;
					app.populateMainUI();
				});
			});
		});

		// disable escaping of html
		$.fn.tmpl.encReg = /[\x00]/g
		$.fn.tmpl.encMap = {};
	},
	listen: function() {
		var errorCount = 0;
		astilectron.onMessage(function(message) {
			var parsed = $.parseJSON(message.payload);
			switch (message.name) {
				case "fatal_error":
					console.log('[' + new Date().toUTCString() + '] ', "fatal_error", parsed.data);
					shared.showError(parsed.data);
					break;
				case "network_stats":
					console.log('[' + new Date().toUTCString() + '] ', "network_stats", parsed);
					$('#abbreviation').html(parsed.abbreviation);
					$('#supply').html(parsed.maximum_supply);
					$('#circulation').html(parsed.circulation);
					$('#market_cap').html('$' + parsed.market_cap);
					// $('#network_hashrate').html(parsed.hashrate);
					// $('#network_difficulty').html(parsed.difficulty);
					// $('#network_height').html(parsed.height);
					$('#price').html('฿' + parsed.price);
					$('#price_usd').html('$' + parsed.price_usd);
					$('#volume_24h_btc').html('฿' + parsed.volume_usd);
					$('#volume_24h_usd').html('$' + parsed.volume);
					$('#miner_payout').html(parsed.coins_per_day);
					$('#pool_hashrate').html(parsed.pool.hashrate);
					$('#pool_miners').html(parsed.pool.miners);
					$('#pool_last_block').html(parsed.pool.last_block);
					// $('#record_volume').html(parsed.records.volume + ' BTC');
					// $('#record_price').html(parsed.records.price + ' BTC');
					$('#pool-address')
						.html('<a href="' + parsed.pool.url + '" class="text-color">' + parsed.pool.url + '</a>')
						.find('a').css('color', app.coinsContent.textColor[app.coin_type]);

					// Build prices
					// let table = '<tbody>';
					// table += '<tr>\
						// <td>Volume today</td>\
						// <td>' + parsed.volume + ' BTC</td>\
					// </tr>';
					// parsed.prices.forEach(function (item) {
						// table += '<tr>\
							// <td>' + item.name + '</td>\
							// <td>' + item.value + ' BTC</td>\
						// </tr>';
					// });
					// table += '</tbody>';
					// $('#exchanges-price').html(table);

					if (!app.selectedPoolOnce) {
						app.selected_pool = parsed.pool.id;
						app.selectedPoolOnce = true;
					}

					app.networkStatsOnce = true;
					if (!app.minerAndNetworkStatsDone && app.minerStatsOnce) {
						app.minerAndNetworkStatsDone = true;
						app.enableSettingsButton();
					}

					break;
				case "miner_stats":
					console.log('[' + new Date().toUTCString() + '] ', "miner_stats", parsed);
					$('#miner_address').html(parsed.address);
					if (!app.populatedAddress) {
						$('#settings_mining_address').val(parsed.address); // settings wallet address
						app.populatedAddress = true;
					}

					$('#miner_hashrate').html(parsed.hashrate_human);
					$('#miner_uptime').html(parsed.uptime_human);
					$('#miner_difficulty').html(parsed.current_difficulty);
					$('#miner_shares').html(parsed.shares_good + parsed.shares_bad);
					$('#miner_shares_bad').html(parsed.shares_bad);

					if (parsed.errors !== null && parsed.errors.length > 0) {
						shared.showError(parsed.errors[0]);
						errorCount++;
						$('#miner_errors').html(errorCount);
						window.setTimeout(function(){
							asticode.modaler.hide();
						}, 30000);
					}

					// Move the graph, we only refresh it once a minute
					if (parsed.update_graph == true) {
						app.hashrateChart.data.datasets.forEach((dataset) => {
							dataset.data.shift();
							dataset.data.push(parsed.hashrate);
						});
						app.hashrateChart.update();
					}

					app.minerStatsOnce = true;
					if (!app.minerAndNetworkStatsDone && app.networkStatsOnce) {
						app.minerAndNetworkStatsDone = true;
						app.enableSettingsButton();
					}
					break;
			}
		});
	},
	// Change colors based on selected coin
	changeColors: function() {
		const textColor = app.coinsContent.textColor[app.coin_type];
		const boxColor = app.coinsContent.boxColor[app.coin_type];
		const boxBorder = app.coinsContent.boxBorder[app.coin_type];

		// change colors
		$('.text-color').css('color', textColor);
		$('body.miner .miner-box').css('background-color', boxColor + 'cc');
		$('body.miner .miner-box.top-middle .estimated')
			.css('background-color', boxColor)
			.css('border-color', boxBorder);
		$('body.miner .miner-box .change-pool .btn')
			.css('color', textColor)
			.css('background-color', boxColor + 'cc')
			.css('border-color', boxBorder);
		$('body.miner .miner-box .pool-address a').css('color', textColor);

		$('.miner-settings .miner-settings-submit .btn')
			.css('color', textColor)
			.css('background-color', boxColor + 'cc')
			.css('border-color', boxBorder);
		$('.miner-settings .miner-settings-cpus').css('background-color', boxColor + 'ff');
		$('.miner-settings .address-input').css('background-color', boxColor + 'ff');
		$('.miner-settings .dropdown-toggle')
			.css('background-color', boxColor + 'ff')
			.css('border-color', boxBorder);
		$('.miner-settings .miner-settings-back')
			.css('background-color', boxColor + 'cc')
			.css('border-color', boxBorder);
		$('.miner-settings .miner-settings-back a').css('color', textColor);

		$('.miner-settings .miner-settings-coins').css('background-color', boxColor + 'ff');

		$('.whatsnext .whatsnext-box').css('background-color', boxColor + '99');
		$('.whatsnext .whatsnext-box a').css('color', textColor);
		$('.whatsnext .whatsnext-back')
			.css('background-color', boxColor + 'cc')
			.css('border-color', boxBorder);
	},
	// Bind to UI events using jQuery
	bindEvents: function() {
		// Functionality based on which page slide is loaded
		$(document).on('click', '[data-target]', function() {
			var id = $(this).data('target');
			if (id == 'miner-settings') {
				app.loadSettings();
			}
			if (id == 'whatsnext') {
				let html;
				html = $.fn.tmpl("tmpl-help-content", app.coinsContent.helpText[app.coin_type].boxes);
				$('#help-content').html(html);
				app.changeColors();
			}
		});

		// Events for miner started/stopped
		$(document).on('miner-started', function() {
			astilectron.sendMessage({
				name: "start-miner",
				payload: ""
			}, function(message) {});
		});
		$(document).on('miner-stopped', function() {
			astilectron.sendMessage({
				name: "stop-miner",
				payload: ""
			}, function(message) {
				app.resetMinerStats();
			});
		});

		// Miner on/off button
		$('#on-off-switch').on('click', function() {
			var el = $(this), offk = el.find('.off-knob'), onk = el.find('.on-knob');
			if (el.hasClass('off')) {
				var rail_width = el.find('.off-rail').width();
				var own_width = offk.width();
				offk.animate({"left": (rail_width - own_width) + "px"}, "slow", function() {
					el.removeClass('off');
					onk.css('right', '0px');
				});
				$('#nice-bg').fadeIn('slow');
				$(document).trigger('miner-started');
			} else {
				var rail_width = el.find('.on-rail').width();
				var own_width = onk.width();
				onk.animate({"right": (rail_width - own_width) + "px"}, "slow", function() {
					el.addClass('off');
					offk.css('left', '0px');
				});
				$('#nice-bg').fadeOut('slow');
				$(document).trigger('miner-stopped');
			}
		});

		// Start the miner on start
		$('#on-off-switch').trigger('click');

		// Save miner settings button
		$('#save-miner-settings').on('click', function(event) {
			event.preventDefault();

			asticode.notifier.close();
			var address = $('#settings_mining_address').val().trim();
			if (address == '') {
				asticode.notifier.error('You must enter your address');
				return;
			} else if (!shared.validateWalletAddress(address, app.coinsContent.addressValidation, app.coin_type)) {
				asticode.notifier.error("Please enter a valid wallet address");
				return;
			}

			if (app.selected_pool == 0) {
				asticode.notifier.error('You must choose one of the pools');
				return;
			}

			asticode.loader.show();

			// Stop the miner first
			astilectron.sendMessage({
				name: "stop-miner",
				payload: ""
			}, function(message) {
				// Save the miner settings
				var configData = {
					address:       $('#settings_mining_address').val().trim(),
					pool:          app.selected_pool,
					coin_type:     app.coin_type,
					coin_algo:     app.coin_algo,
					xmrig_algo:    app.xmrig_algo,
					xmrig_variant: app.xmrig_variant.toString(),
					// threads:       parseInt($('#cpu-cores').dropselect('value')),
					// threads:       1,
					// max_cpu:       parseInt($('#cpu-max').dropselect('value'))
					// max_cpu:       100,
					hardware_type: parseInt($('#hardware-type').dropselect('value'))
				};
				console.log('[' + new Date().toUTCString() + '] ', "save-configuration", configData);
				astilectron.sendMessage({name: "save-configuration", payload: configData}, function(message){
					document.location = 'index.html';
				});
			});
		});
	},
	// Change the UI based on the selected coin
	populateMainUI: function() {
		app.setupChart();

		let html;

		// Background
		$('body')
			.css('background-image', 'url(' + app.coinsContent.mainBackground[app.coin_type].image + ')')
			.css('background-color', app.coinsContent.mainBackground[app.coin_type].color);

		// Logo
		$('#miner-middle-logo').html(
			$('<img/>').attr('src', app.coinsContent.logo[app.coin_type])
		);

		// Network links
		html = $.fn.tmpl("tmpl-network-links", app.coinsContent.networkLinks[app.coin_type]);
		$('#network-links').html(html);

		// Social links
		html = $.fn.tmpl("tmpl-social-links", app.coinsContent.socialLinks[app.coin_type]);
		$('#social-links').html(html);

		// News box
		if (app.coinsContent.newsBox[app.coin_type].title !== "" && app.coinsContent.newsBox[app.coin_type].image !== "") {
			$('#news-title').text(app.coinsContent.newsBox[app.coin_type].title);
			$('#news-image').attr("src", app.coinsContent.newsBox[app.coin_type].image);
			if (app.coinsContent.newsBox[app.coin_type].link !== "") {
				$('#news-link').attr("href", app.coinsContent.newsBox[app.coin_type].link);
			} else {
				$('#news-link').attr("href", "javascript:;");
			}
			$('#news-box').removeClass('hidden');
		} else {
			$('#news-box').addClass('hidden');
		}

		// Various replacements
		$('#miner_coin').text(app.coinsContent.abbreviation[app.coin_type]);
		$('#cryptunit-widget').html(app.coinsContent.cryptunitWidget[app.coin_type]);
		$('#download-title').text(app.coinsContent.downloadPage[app.coin_type].title);
		$('#download-title-link').attr('href', app.coinsContent.downloadPage[app.coin_type].link);
		$('#download-link').attr('href', app.coinsContent.downloadPage[app.coin_type].link);

		// Powered-by links
		html = $.fn.tmpl("tmpl-powered-by-links", app.coinsContent.poweredByLinks);
		$('#powered-by-links').html(html).show();

		// CoinGecko links
		$('#coingecko-link').attr('href', app.coinsContent.coinGeckoLinks[app.coin_type]);

		app.changeColors();
	},
	loadSettings: function() {
		// Get the current miner processing config
		astilectron.sendMessage({name: "get-processing-config", payload: ""}, function(message) {
			var parsed = $.parseJSON(message.payload);
			console.log('[' + new Date().toUTCString() + '] ', "get-processing-config", parsed);

			let html;

			// Populate the number of cores
			// html = "";
			// for (let i = 1; i <= parsed.max_threads; i++) {
				// html += '<li><a href="#" data-value="' + i + '"><b>' + i + ' CPU</b> CORE</a></li>';
			// }
			// $('#cpu-cores-values').html(html);
			// $('#cpu-cores').dropselect(); // re-initialize
			// $('#cpu-cores').dropselect('select', parsed.threads);

			// Select the appropriate cpu usage
			// $('#cpu-max').dropselect(); // re-initialize
			// $('#cpu-max').dropselect('select', parsed.max_usage);

			// Select the mining hardware
			$('#hardware-type').dropselect(); // re-initialize
			$('#hardware-type').dropselect('select', parsed.hardware_type);

			// Populate the coins
			let coins = {
				'selected': app.coinsContent.names[app.coin_type]
			};
			coins.coins = app.coinsContent.coins.map(function(el) { // add name and abbreviation keys
				el.name = app.coinsContent.names[el.coin_type];
				return el;
			});
			html = $.fn.tmpl("tmpl-miner-settings-coins", coins);
			$('#miner-settings-coins').html(html);

			// Enable the dropselect
			$('#settings-coins').dropselect(); // re-initialize
			$('#settings-coins').dropselect('change', function() {
				asticode.loader.show();

				// Add the xmrig_algo and xmrig_variant to the coins
				let all_coins = app.coinsContent.coins.map(function(el) {
					el.xmrig_algo    = app.coinsContent.xmrigAlgo[el.coin_type];
					el.xmrig_variant = app.coinsContent.xmrigVariant[el.coin_type];
					return el;
				});

				// Change coin_type/coin_algo
				const sel_coin_type = $('#settings-coins').dropselect('value'); // Selected coin_type
				let curr_coin = all_coins.filter(function(el) { // add name and abbreviation keys
					return el.coin_type == sel_coin_type;
				});

				app.prev_coin_type = app.coin_type;
				app.coin_type = sel_coin_type;
				app.coin_algo = curr_coin[0].coin_algo;
				app.xmrig_algo = curr_coin[0].xmrig_algo;
				app.xmrig_variant = String(curr_coin[0].xmrig_variant);

				// The coin has changed, reload the pools
				app._loadPoolsInSettings(function() {
					asticode.loader.hide();
				});
			});

			app._loadPoolsInSettings();
		});
	},
	_loadPoolsInSettings: function(cb) {
		// Populates the pool list in the settings
		var payloadData = {
			coin_type: app.coin_type
		};
		astilectron.sendMessage({name: "get-pools-list", payload: payloadData}, function(message) {
			// console.log('[' + new Date().toUTCString() + '] ', "get-pools-list", message.payload);

			$("#pool-list").mCustomScrollbar("destroy");
			$('#pool-list').html(message.payload);
			$("#pool-list").mCustomScrollbar({
				theme:"rounded-dots",
				scrollInertia:400
			});

			app.selected_pool = 0;

			$('div.table-body').off('click').on('click', function() {
				$(this).parent().find('.table-body').removeClass('selected');
				$(this).addClass('selected');
				app.selected_pool = parseInt($(this).data('id'));
			});

			// Select the pool we are currently using
			$('#pool-list').find('.table-body').each(function() {
				var id = $(this).data('id');
				if (id == app.selected_pool) {
					$(this).trigger('click');
					return false;
				}
			});

			if (typeof cb === 'function') cb();
		});
	},
	resetMinerStats: function() {
		$('#miner_hashrate').html('0.00 H/s');
		$('#miner_uptime').html('0 seconds');
		$('#miner_difficulty').html('0');
		$('#miner_shares').html('0');
		$('#miner_shares_bad').html('0');
		$('#miner_payout').html('0');
	},
	setupChart: function() {
		var chart = $("#hashrate_chart");
		chart.attr('width', $('.chart').width());
		chart.attr('height', $('.chart').height());

		app.hashrateChart = new Chart(chart, {
			type: 'line',
			data: {
				labels: ["5min ago", "4min ago", "3min ago", "2min ago", "1min ago", "Now"],
				datasets: [{
					label: 'H/s',
					data: [0,0,0,0,0,0],
					// backgroundColor: ['rgba(13, 17, 45, 1.0)'],
					backgroundColor: [app.coinsContent.boxColor[app.coin_type]],
					// borderColor: ['rgba(232,212,0,9)'],
					borderColor: [app.coinsContent.textColor[app.coin_type]],
					borderWidth: 1,
				}]
			},
			options: {
				tooltips: {
					mode: 'index',
					intersect: false,
				},
				legend: {
					display: false,
				},
				elements: {
					line: {
						tension: 0, // disables bezier curves
					},
				},
				layout: {
					padding: {
						left: 0,
						right: 0,
						top: 0,
						bottom: 0
					}
				},
				scales:
				{
					yAxes: [{
					//display: false
					}]
				}
			}
		});
	},
	enableSettingsButton: function() {
		$('#loading-settings').addClass('done');
	},
	prev_coin_type: '',
	coin_type: '',
	coin_algo: '',
	xmrig_algo: '',
	xmrig_variant: '',
	selected_pool: 0,
	populatedAddress: false,
	networkStatsOnce: false,
	minerStatsOnce: false,
	selectedPoolOnce: false,
	minerAndNetworkStatsDone: false,
	coinsContent: {}
}