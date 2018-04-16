let app = {
    init: function() {
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        //open links externally by default
        var shell = require('electron').shell;
        $(document).on('click', 'a[href^="http"]', function(event) {
            event.preventDefault();
            shell.openExternal(this.href);
        });

        // Wait for the ready signal
        document.addEventListener('astilectron-ready', function() {
          // Start the miner on start
          astilectron.sendMessage({
            name: "miner_start",
            payload: ""
          }, function(message) {
          });

          $('#start_stop').bind('click', function(e){
            var isStarted = $(this).hasClass('stop');
            if (isStarted) {
              // Stop the miner
              astilectron.sendMessage({
                name: "miner_stop",
                payload: ""
              }, function(message) {
                $('#start_stop').addClass('start');
                $('#start_stop').removeClass('stop');
                $('#start_stop').html('Start mining');

                $('#miner_hashrate').html('0');
                $('#miner_uptime').html('0');
                $('#miner_difficulty').html('0');
                $('#miner_shares').html('0');
                $('#miner_shares_bad').html('0');
                $('#miner_payout').html('0.00 XTL');
              });
            } else {
              // Start the miner
              astilectron.sendMessage({
                name: "miner_start",
                payload: ""
              }, function(message) {
                $('#start_stop').addClass('stop');
                $('#start_stop').removeClass('start');
                $('#start_stop').html('Stop mining');
              });
            }
            e.stopPropogation();
            return false;
          });


          app.listen();
        })
    },
    listen: function() {
      astilectron.onMessage(function(message) {
            switch (message.name) {
                case "stats":
                  var stats = $.parseJSON(message.payload)
                  $('#circulation').html(stats.circulation);
                  $('#market_cap').html(stats.market_cap);
                  $('#price').html(stats.price + ' BTC');
                  $('#network_hashrate').html(stats.hashrate);
                  $('#network_difficulty').html(stats.last_block.difficulty);
                  $('#network_height').html(stats.last_block.height);
                  $('#trading_volume').html(stats.volume + ' BTC');
                  $('#trading_tradeogre_volume').html(stats.volume_tradeogre + ' BTC');
                  $('#trading_crex_volume').html(stats.volume_crex + ' BTC');
                  $('#record_volume').html(stats.records.volume + ' BTC');
                  $('#record_price').html(stats.records.price + ' BTC');
                  $('#miner_payout').html(stats.xtl_per_day + ' XTL');
                  $('#pool').html(stats.pool_html);
                  break;
                case "miner_stats":
                  var stats = $.parseJSON(message.payload)
                  $('#miner_hashrate').html(stats.hashrate.total[0]);
                  $('#miner_uptime').html(app.secondsHumanize(stats.connection.uptime));
                  $('#miner_difficulty').html(stats.results.diff_current);
                  $('#miner_shares').html(stats.results.shares_total);
                  $('#miner_shares_bad').html(stats.results.shares_total - stats.results.shares_good);
                  $('#miner_address').html(stats.address);

                  hashrateChart.data.datasets.forEach((dataset) => {
                      dataset.data.shift();
                      dataset.data.push(stats.hashrate.total[0]);
                  });
                  hashrateChart.update();
                  break;
                case "about":
                  index.about(message.payload);
                  return {payload: "payload"};
                  break;
                case "check.out.menu":
                  asticode.notifier.info(message.payload);
                  break;
            }
        });
    },
    secondsHumanize: function(d) {
        d = Number(d);
        var h = Math.floor(d / 3600);
        var m = Math.floor(d % 3600 / 60);
        var s = Math.floor(d % 3600 % 60);

        var hDisplay = h > 0 ? h + (h == 1 ? " hour, " : " hours, ") : "";
        var mDisplay = m > 0 ? m + (m == 1 ? " minute, " : " minutes, ") : "";
        var sDisplay = s > 0 ? s + (s == 1 ? " second" : " seconds") : "";
        if (m > 0) {
          return hDisplay + mDisplay;
        }
        return hDisplay + mDisplay + sDisplay;
    }
};
