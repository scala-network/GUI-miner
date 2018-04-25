/*
  This handles the miner
 */
const remote = require('electron').remote;
let app = {
    init: function() {
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        shared.bindExternalLinks();
        app.setupChart();

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
            shared.showError(parsed.Data);
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
              app.hashrateChart.data.datasets.forEach((dataset) => {
                dataset.data.shift();
                dataset.data.push(parsed.hashrate);
              });
              app.hashrateChart.update();
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
        e.stopPropagation();
        return false;
      });

      $('.header-button.settings').bind('click', function(){
        app.loadSettings();
      });
      $('.header-button.help').bind('click', function(){
        $('#help').toggleClass('dn');
      });
      $('.header-button.minimize').bind('click', function(){
        remote.getCurrentWindow().minimize();
      });
      $('.header-button.exit').bind('click', function(){
        remote.getCurrentWindow().close();
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
          pool: $('#pool_list').find('.selected').data('id'),
          threads: parseInt($('#threads option:selected').attr('value')),
          max_cpu: parseInt($('#max_cpu option:selected').attr('value'))
        };
        if (configData.address == '') {
          alert("You must enter your address");
          return false;
        }
        // Just make sure they're not using integrated addresses or
        // invalid ones
        if (shared.validateWalletAddress(configData.address) == false)
        {
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

      // get-processing-config get the current miner processing config
      astilectron.sendMessage({name: "get-processing-config", payload: ""}, function(message) {
        var parsed = $.parseJSON(message.payload)
        $('#max_threads').html(parsed.max_threads);
        if (parsed.max_threads <= 1) {
          $('#max_threads_multiple').hide();
        } else $('#max_threads_multiple').show();

        if (parsed.type == 'xmrig') {
          $('.xmrig-extra').show();
        } else $('.xmrig-extra').hide();

        // For xmrig's GPU only setup we don't show the CPU tuning options
        if (parsed.type != "xmrig-gpu") {
          // TODO: Do this in a better way, i.e - not as text
          var threadOptions = "<select>";
          var startThreadCount = 1;
          if (parsed.type == 'xmr-stak') {
            startThreadCount = 0;
          }
          for (var i = startThreadCount; i <= parsed.max_threads; i++) {
            if (i == parsed.threads) {
              threadOptions += '<option value="' + i + '" selected>' + i + '</option>';
            } else {
              threadOptions += '<option value="' + i + '">' + i +'</option>';
            }
          }
          threadOptions += "</select>";
          $('#threads').html(threadOptions);
          $('#threads select').niceSelect();
          // Not set means 100%
          if (parsed.max_usage == 0) {
            $('#max_cpu select').find('option[value=100]').attr('selected','selected');
          }
          else $('#max_cpu select').find('option[value=' + parsed.max_usage + ']').attr('selected','selected');
          $('#max_cpu select').niceSelect();
        } else {
          // GPU only
          $('.cpu-tuning').hide();
          $('.gpu-tuning').show();
        }
      });

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
            backgroundColor: [
              'rgba(13, 17, 45, 1.0)'
            ],
            borderColor: [
              'rgba(232,212,0,9)'
            ],
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
};
