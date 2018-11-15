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
					$('#miner_payout').html(parsed.bloc_per_day + ' BLOC');
					$('#pool_hashrate').html(parsed.pool.hashrate);
					$('#pool_miners').html(parsed.pool.miners);
					$('#pool_last_block').html(parsed.pool.last_block);
					$('#trading_volume').html(parsed.volume + ' BTC');
					$('#record_volume').html(parsed.records.volume + ' BTC');
					$('#record_price').html(parsed.records.price + ' BTC');
					$('#price_stex').html(parsed.price_stex + ' BTC');
					$('#price_tradeogre').html(parsed.price_tradeogre + ' BTC');
					$('#pool-address').html('<a href="' + parsed.pool.url + '">' + parsed.pool.url + '</a>').data('id', parsed.pool.id);
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
				name: "miner_start",
				payload: ""
			}, function(message) {});
		});
		$(document).on('miner-stopped', function() {
			astilectron.sendMessage({
				name: "miner_stop",
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
				name: "miner_stop",
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
		$('#miner_payout').html('0 BLOC');
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
	populatedAddress: false
}