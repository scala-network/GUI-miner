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

          app.bindEvents();
          app.listen();
        })
    },
    listen: function() {
      var errorCount = 0;
      astilectron.onMessage(function(message) {
        var parsed = $.parseJSON(message.payload)
        switch (message.name) {
          case "fatal_error":
          let errDiv = document.createElement("div");
          errDiv.innerHTML = parsed.message;
          $('.astimodaler-body').addClass('error');
          asticode.modaler.setContent(errDiv);
          asticode.modaler.show();
            break;
          case "network_stats":
            $('#circulation').html(parsed.circulation);
            $('#market_cap').html(parsed.market_cap);
            $('#price').html(parsed.price + ' BTC');
            $('#network_hashrate').html(parsed.hashrate);
            $('#network_difficulty').html(parsed.last_block.difficulty);
            $('#network_height').html(parsed.last_block.height);
            $('#trading_volume').html(parsed.volume + ' BTC');
            $('#trading_tradeogre_volume').html(parsed.volume_tradeogre + ' BTC');
            $('#trading_crex_volume').html(parsed.volume_crex + ' BTC');
            $('#record_volume').html(parsed.records.volume + ' BTC');
            $('#record_price').html(parsed.records.price + ' BTC');
            $('#miner_payout').html(parsed.xtl_per_day + ' XTL');
            $('#pool').html(parsed.pool_html);
            break;
          case "miner_stats":
            $('#miner_hashrate').html(parsed.hashrate.total[0]);
            $('#miner_uptime').html(app.secondsHumanize(parsed.connection.uptime));
            $('#miner_difficulty').html(parsed.results.diff_current);
            $('#miner_shares').html(parsed.results.shares_total);
            $('#miner_shares_bad').html(parsed.results.shares_total - parsed.results.shares_good);
            $('#miner_address').html(parsed.address);
            // Move the graph, we only refresh it once a minute
            if (parsed.update_graph == true) {
              hashrateChart.data.datasets.forEach((dataset) => {
                dataset.data.shift();
                dataset.data.push(parsed.hashrate.total[0]);
              });
              hashrateChart.update();
            }

            if (parsed.connection.error_log.length > 0) {
              let errDiv = document.createElement("div");
              errDiv.innerHTML = parsed.connection.error_log[0].text;
              $('.astimodaler-body').addClass('error');
              asticode.modaler.setContent(errDiv);
              asticode.modaler.show();
              errorCount++;
              $('#miner_errors').html(errorCount);
              window.setTimeout(function(){
                asticode.modaler.hide();
              }, 5000);
            }
            break;
          }
        });
    },
    // Bind to UI events using jQuery
    bindEvents: function() {
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
    },
    // secondsHumanize turns seconds into hours + minutes
    secondsHumanize: function(d) {
        d = Number(d);
        var h = Math.floor(d / 3600);
        var m = Math.floor(d % 3600 / 60);
        var s = Math.floor(d % 3600 % 60);

        var hDisplay = h > 0 ? h + (h == 1 ? " hour, " : " hours, ") : "";
        var mDisplay = m > 0 ? m + (m == 1 ? " minute" : " minutes") : "";
        var sDisplay = s > 0 ? s + (s == 1 ? " second" : " seconds") : "";
        if (m > 0) {
          return hDisplay + mDisplay;
        }
        return hDisplay + mDisplay + sDisplay;
    }
};
