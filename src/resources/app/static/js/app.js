let app = {
    init: function() {
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for the ready signal
        document.addEventListener('astilectron-ready', function() {
          astilectron.sendMessage({
            name: "miner_start",
            payload: ""
          }, function(message) {
            console.log("Send miner_start")
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
                  $('#price').html(stats.price);
                  $('#network_hashrate').html(stats.hashrate);
                  $('#network_difficulty').html(stats.last_block.difficulty);
                  $('#network_height').html(stats.last_block.height);
                  $('#trading_volume').html(stats.volume);
                  $('#trading_tradeogre_volume').html(stats.volume_tradeogre);
                  $('#trading_crex_volume').html(stats.volume_crex);
                  $('#record_volume').html(stats.records.volume);
                  $('#record_price').html(stats.records.price);
                  $('#miner_payout').html(stats.xtl_per_day);
                  break;
                case "miner_stats":
                  var stats = $.parseJSON(message.payload)
                  $('#miner_hashrate').html(stats.hashrate.total[0]);
                  $('#miner_uptime').html(stats.connection.uptime);
                  $('#miner_difficulty').html(stats.results.diff_current);
                  $('#miner_shares').html(stats.results.shares_total);
                  $('#miner_shares_bad').html(stats.results.shares_total - stats.results.shares_good);
                  $('#miner_address').html(stats.address);
                case "miner_start":
                  // TODO: Update miner status
                  break;
                case "miner_end":
                  // TODO: Update miner status
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
    }
};
