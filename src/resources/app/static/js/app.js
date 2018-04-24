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
        // This stops electron from updating the window title when a link
        // is clicked
        $(document).on('click', 'a[href^="#"]', function(event) {
            event.preventDefault();
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
            $('#network_difficulty').html(parsed.difficulty);
            $('#network_height').html(parsed.height);
            $('#trading_volume').html(parsed.volume + ' BTC');
            $('#trading_tradeogre_volume').html(parsed.volume_tradeogre + ' BTC');
            $('#trading_crex_volume').html(parsed.volume_crex + ' BTC');
            $('#record_volume').html(parsed.records.volume + ' BTC');
            $('#record_price').html(parsed.records.price + ' BTC');
            $('#miner_payout').html(parsed.xtl_per_day + ' XTL');
            $('#pool').html(parsed.pool_html);
            break;
          case "miner_stats":
            $('#miner_hashrate').html(parsed.hashrate_human);
            $('#miner_uptime').html(parsed.uptime_human);
            $('#miner_difficulty').html(parsed.current_difficulty);
            $('#miner_shares').html(parsed.shares_good + parsed.shares_bad);
            $('#miner_shares_bad').html(parsed.shares_bad);
            $('#miner_address').html(parsed.address);
            // Move the graph, we only refresh it once a minute
            if (parsed.update_graph == true) {
              hashrateChart.data.datasets.forEach((dataset) => {
                dataset.data.shift();
                dataset.data.push(parsed.hashrate);
              });
              hashrateChart.update();
            }

            if (parsed.errors !== null && parsed.errors.length > 0) {
              let errDiv = document.createElement("div");
              errDiv.innerHTML = parsed.errors[0];
              $('.astimodaler-body').addClass('error');
              asticode.modaler.setContent(errDiv);
              asticode.modaler.show();
              errorCount++;
              $('#miner_errors').html(errorCount);
              window.setTimeout(function(){
                asticode.modaler.hide();
              }, 4000);
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

            app.resetMinerStats();
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

      $('.settings-button').bind('click', function(){
        app.loadSettings();
      });
      $('.help-button').bind('click', function(){
        $('#help').toggleClass('dn');
      });

      $(document).on('click', '#change_pool', function(){
        app.loadSettings();
      });

      $('.close-settings').bind('click', function(){
        $('#settings').toggleClass('dn');
      });

      $('.close-help').bind('click', function(){
        $('#help').toggleClass('dn');
      });

      $(document).on('click', '.pool', function(){
        $('.pool').removeClass('selected');
        $(this).addClass('selected');
      });

      $('#update').bind('click', function(){
        var configData = {
          address: $('#settings_mining_address').val(),
          pool: $('#pool_list').find('.selected').data('id')
        };
        if (configData.address == '') {
          alert("You must enter your address");
          return false;
        }
        // Just make sure they're not using integrated addresses
        if (configData.address.substring(0, 2) != 'Se') {
          alert("Please enter a valid Stellite address starting with 'Se'");
          return false;
        }
        $('#update').html('Updating...');
        astilectron.sendMessage({name: "reconfigure", payload: configData}, function(message){
          $('.current .pool h3').html('Updating...');
          $('#settings').toggleClass('dn');
          $('#update').html('Update');
          $('#miner_address').html("Updating")
          app.resetMinerStats();
          asticode.notifier.info('Miner reconfigured');
        });

      });
    },
    loadSettings: function() {
      $('#settings_mining_address').val($('#miner_address').html());
      // The pool-list command returns the pool list for the GUI miner
      astilectron.sendMessage({name: "pool-list", payload: ""}, function(message) {
        $('#pool_list').html(message.payload);
        var currentPool = $('.current .pool').data('id');
        $('.pool[data-id="' + currentPool + '"]').addClass('selected');
        $('#settings').toggleClass('dn');
      });
    },
    resetMinerStats: function() {
      $('#miner_hashrate').html('0 H/s');
      $('#miner_uptime').html('0');
      $('#miner_difficulty').html('0');
      $('#miner_shares').html('0');
      $('#miner_shares_bad').html('0');
      $('#miner_payout').html('0.00 XTL');
    },
};
