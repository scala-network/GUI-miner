/*
  This handles the actual miner
 */
let app = {
	init: function() {
		asticode.loader.init();
		asticode.modaler.init();
		asticode.notifier.init();

		shared.bindExternalLinks();
		app.setupChart();

		shared.bindTargetLinks();
		shared.extraFunctions();

		// Wait for the ready signal
		document.addEventListener('astilectron-ready', function() {
			app.bindEvents();
			app.listen();
		});
	},
	listen: function() {
		var errorCount = 0;
		astilectron.onMessage(function(message) {
			var parsed = $.parseJSON(message.payload);
			switch (message.name) {
				case "fatal_error":
					shared.showError(parsed.data);
					break;
				case "network_stats":
					console.log('[' + new Date().toUTCString() + '] ', "network_stats", parsed);
					$('#circulation').html(parsed.circulation);
					$('#market_cap').html(parsed.market_cap);
					$('#network_hashrate').html(parsed.hashrate);
					$('#network_difficulty').html(parsed.difficulty);
					$('#network_height').html(parsed.height);
					$('#price').html(parsed.price + ' BTC');
					$('#miner_payout').html(parsed.bloc_per_day);
					$('#pool_hashrate').html(parsed.pool.hashrate);
					$('#pool_miners').html(parsed.pool.miners);
					$('#pool_last_block').html(parsed.pool.last_block);
					$('#record_volume').html(parsed.records.volume + ' BTC');
					$('#record_price').html(parsed.records.price + ' BTC');
					$('#pool-address').html('<a href="' + parsed.pool.url + '">' + parsed.pool.url + '</a>').data('id', parsed.pool.id);

					// Build prices
					let table = '<tbody>';
					table += '<tr>\
						<td>Volume today</td>\
						<td>' + parsed.volume + ' BTC</td>\
					</tr>';
					parsed.prices.forEach(function (item) {
						table += '<tr>\
							<td>' + item.name + '</td>\
							<td>' + item.value + ' BTC</td>\
						</tr>';
					});
					table += '</tbody>';
					$('#exchanges-price').html(table);

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
	// Bind to UI events using jQuery
	bindEvents: function() {
		// Functionality based on which page slide is loaded
		$(document).on('click', '[data-target]', function() {
			var id = $(this).data('target');
			if (id == 'miner-settings') {
				app.loadSettings(function() {
					$("#pool-list").mCustomScrollbar({
						theme:"rounded-dots",
						scrollInertia:400
					});

					$('div.table-body').off('click').on('click', function() {
						$(this).parent().find('.table-body').removeClass('selected');
						$(this).addClass('selected');
					});

					// Select the pool we are currently using
					var curr_pool_id = $('#pool-address').data('id');
					$('#pool-list').find('.table-body').each(function() {
						var id = $(this).data('id');
						if (id == curr_pool_id) {
							$(this).trigger('click');
							return false;
						}
					});
				});
			}
		});

		// Events for miner started/stopped
		$(document).on('miner-started', function() {
			astilectron.sendMessage({
				name: "miner-start",
				payload: ""
			}, function(message) {});
		});
		$(document).on('miner-stopped', function() {
			astilectron.sendMessage({
				name: "miner-stop",
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
		$('#save-miner-settings').on('click', function() {
			asticode.loader.show();

			// Stop the miner first
			astilectron.sendMessage({
				name: "miner-stop",
				payload: ""
			}, function(message) {
				// Save the miner settings
				var configData = {
					address: $('#settings_mining_address').val(),
					pool: $('#pool-list').find('.selected').data('id'),
					threads: parseInt($('#cpu-cores').dropselect('value')),
					max_cpu: parseInt($('#cpu-max').dropselect('value'))
				};
				console.log('[' + new Date().toUTCString() + '] ', "configure", configData);
				astilectron.sendMessage({name: "configure", payload: configData}, function(message){
					document.location = 'index.html';
				});
			});
		});

		// Get the configuration settings
		astilectron.sendMessage({
			name: "config-file",
			payload: ""
		}, function(message) {
			var parsed = $.parseJSON(message.payload);
			console.log('[' + new Date().toUTCString() + '] ', "config-file", parsed);
			if (app.coin_type != parsed.coin_type || app.coin_algo != parsed.coin_algo) {
				app.coin_type = parsed.coin_type;
				app.coin_algo = parsed.coin_algo;
				app.populateMainUI();
			}
		});
	},
	// Change the UI based on the selected coin
	populateMainUI: function() {
		if (typeof app.UI.networkLinks[app.coin_type] !== 'undefined') {
			// Background
			const bg_img = app.UI.mainBg[app.coin_type].image;
			const bg_color = app.UI.mainBg[app.coin_type].color;
			$('body').css('background-image', 'url(' + bg_img + ')').css('background-color', bg_color);

			// Logo
			const img = $('<img/>').attr('src', app.UI.logo[app.coin_type]);
			$('#miner-middle-logo').html(img);

			// Network links
			$('#network-text').html(app.UI.networkLinks[app.coin_type].title);
			let html_nt = '';
			app.UI.networkLinks[app.coin_type].links.forEach(function (item) {
				html_nt += '<li><a href="' + item.link + '">' + item.text + '</a></li>';
			});
			$('#network-links').html(html_nt);

			// Social links
			$('#social-text').html(app.UI.socialLinks[app.coin_type].title);
			let html_sl = '';
			app.UI.socialLinks[app.coin_type].links.forEach(function (item) {
				html_sl += '<li><a href="' + item.link + '"><img src="' + item.img + '"></a></li>';
			});
			$('#social-links').html(html_sl);

			// Various replacements
			$('#miner_coin').text(app.UI.abbr[app.coin_type]);
			$('#download-title').text(app.UI.downloadPage[app.coin_type].title);
			$('#download-link').attr('href', app.UI.downloadPage[app.coin_type].link);
		}
	},
	loadSettings: function(cb) {
		// Get the current miner processing config
		astilectron.sendMessage({name: "get-processing-config", payload: ""}, function(message) {
			var parsed = $.parseJSON(message.payload);
			console.log('[' + new Date().toUTCString() + '] ', "get-processing-config", parsed);

			// Populate the number of cores
			var threadOptions = "";
			for (var i = 1; i <= parsed.max_threads; i++) {
				threadOptions += '<li><a href="#" data-value="' + i + '"><b>' + i + ' CPU</b> CORE</a></li>';
			}
			$('#cpu-cores-values').html(threadOptions);
			$('#cpu-cores').dropselect(); // re-initialize
			$('#cpu-cores').dropselect('select', parsed.threads);

			// Select the appropriate cpu usage
			$('#cpu-max').dropselect(); // re-initialize
			$('#cpu-max').dropselect('select', parsed.max_usage);

			// Return the pool list for the GUI miner
			astilectron.sendMessage({name: "pool-list", payload: ""}, function(message) {
				$("#pool-list").mCustomScrollbar("destroy");
				$('#pool-list').html(message.payload);
				if (typeof cb === 'function') cb();
			});
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
				labels: ["5 minutes ago", "4 minutes ago", "3 minutes ago", "2 minutes ago", "1 minute ago", "Now"],
				datasets: [{
					label: 'H/s',
					data: [0,0,0,0,0,0],
					backgroundColor: ['rgba(13, 17, 45, 1.0)'],
					borderColor: ['rgba(232,212,0,9)'],
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
	coin_type: '',
	coin_algo: '',
	populatedAddress: false,
	networkStatsOnce: false,
	minerStatsOnce: false,
	minerAndNetworkStatsDone: false,
	UI: {
		abbr: {
			"bloc": "BLOC",
			"turtlecoin": "TRTL"
		},
		mainBg: {
			"bloc": {
				"image": "static/img/bg-miner-mining.png", 
				"color": "#001b45"
			},
			"turtlecoin": {
				"image": "static/img/turtlecoin/bg-miner-mining.png", 
				"color": "#053F04"
			}
		},
		logo: {
			"bloc": "static/img/miner/miner-big-logo.png",
			"turtlecoin": "static/img/turtlecoin/logo.png"
		},
		downloadPage: {
			"bloc": {
				'title': 'Bloc Applications', 
				'link': "https://bloc.money/download"
			},
			"turtlecoin": {
				'title': 'Turtlecoin Applications', 
				'link': "https://github.com/turtlecoin/turtlecoin/releases/latest"
			}
		},
		socialLinks: {
			"bloc": {"title": "Social Network", "links": [
				{"link": "https://bloc.money/", "img": "static/img/miner/network_icon_website.png"},
				{"link": "https://discord.gg/5Buudya", "img": "static/img/miner/network_icon_discord.png"},
				{"link": "https://t.me/bloc_money", "img": "static/img/miner/network_icon_telegram.png"},
				{"link": "https://bitcointalk.org/index.php?topic=4108831.0", "img": "static/img/miner/network_icon_bitcointalk.png"},
				{"link": "https://github.com/furiousteam", "img": "static/img/miner/network_icon_github.png"},
				{"link": "https://twitter.com/bloc_money", "img": "static/img/miner/network_icon_twatter.png"},
				{"link": "https://medium.com/@bloc.money", "img": "static/img/miner/network_icon_medium.png"},
				{"link": "https://www.youtube.com/channel/UCdvnEPWhqGtZUEx3EFBrXvA", "img": "static/img/miner/network_icon_yuotube.png"},
				{"link": "https://www.facebook.com/Blocmoney-383098922176113", "img": "static/img/miner/network_icon_facbook.png"},
				{"link": "https://www.instagram.com/bloc.money", "img": "static/img/miner/network_icon_instgram.png"}
			]},
			"turtlecoin": {"title": "Social Network", "links": [
				{"link": "https://twitter.com/_turtlecoin", "img": "static/img/miner/network_icon_twatter.png"},
				{"link": "http://chat.turtlecoin.lol/", "img": "static/img/miner/network_icon_discord.png"},
				{"link": "https://github.com/turtlecoin", "img": "static/img/miner/network_icon_github.png"},
				{"link": "https://www.facebook.com/trtlcoin/", "img": "static/img/miner/network_icon_facbook.png"},
				{"link": "https://www.instagram.com/_turtlecoin/", "img": "static/img/miner/network_icon_instgram.png"},
				{"link": "https://www.reddit.com/r/TRTL/", "img": "static/img/miner/network_icon_website.png"}
			]}
		},
		networkLinks: {
			"bloc": {"title": "The BLOC Network", "links": [
				{"link": "https://bloc.money/", "text": "- Official website"},
				{"link": "https://itunes.apple.com/us/app/bloc-wallet-by-furiousteam-ltd/id1437924269?mt=8", "text": "- iPhone app"},
				{"link": "https://bloc-explorer.com/", "text": "- Bloc Explorer"},
				{"link": "https://bloc-developer.com/", "text": "- Bloc Developper"},
				{"link": "https://bloc-mining.com/", "text": "- Bloc Web Mining"},
				{"link": "https://bloc-mining.eu/", "text": "- Bloc Mining EU"},
				{"link": "https://bloc-mining.us/", "text": "- Bloc Mining US"},
				{"link": "https://bloc-mining.asia/", "text": "- Bloc Mining ASIA"},
				{"link": "https://bloc.cool/", "text": "- Bloc Cool"},
				{"link": "https://t.me/bloc_wallet_bot", "text": "- Telegram wallet"},
				{"link": "https://t.me/bloc_explorer_bot", "text": "- Telegram explorer"},
				{"link": "https://paychange.com/", "text": "- PayChange"},
				{"link": "https://traakx.com/", "text": "- Traakx"}
			]},
			"turtlecoin": {"title": "The TRTL Network", "links": [
				{"link": "https://turtlecoin.lol", "text": "- Official Website"},
				{"link": "https://turtlewallet.lol", "text": "- Turtle Wallet"},
				{"link": "https://trtl.services", "text": "- Turtle Services"},
				{"link": "https://explorer.turtlepay.io", "text": "- Turtle Explorer"},
				{"link": "https://explorer.turtlepay.io/pools.html", "text": "- Turtle Pools"},
				{"link": "https://turtlepay.io", "text": "- TurtlePay"},
				{"link": "http://wiki.turtlecoin.lol", "text": "- Turtle Wiki"},
				{"link": "https://medium.com/@turtlecoin", "text": "- Developer Blog"}
			]}
		}
	}
}