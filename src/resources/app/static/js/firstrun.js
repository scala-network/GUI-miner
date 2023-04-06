/*
  This handles the initial user setup
 */
const remote = require('electron').remote;
let firstrun = {
  init: function() {
      asticode.loader.init();
      asticode.modaler.init();
      asticode.notifier.init();

      shared.bindExternalLinks();

      // Wait for the ready signal
      document.addEventListener('astilectron-ready', function() {

        astilectron.sendMessage({name: "firstrun", payload: ""}, function(message) {
          $('#username').html(message.payload);
        });

        astilectron.sendMessage({name: "get-miner-path", payload: ""}, function(message) {
          $('#miner_path').html(message.payload);
        });

        // The pool-list command returns the pool list for the GUI miner
        astilectron.sendMessage({name: "pool-list", payload: ""}, function(message) {
          $('#pool_list').html(message.payload);
        });

        firstrun.bindEvents();
        firstrun.listen();

        // Just wait a second for the window to show and the user to focus
        window.setTimeout(function(){
          firstrun.animateIntro();
        }, 2000);
      })
  },
  listen: function() {
    astilectron.onMessage(function(message) {
      var parsed = $.parseJSON(message.payload);
      switch (message.name) {
        case "fatal_error":
          shared.showError(parsed.data);
      }
    });
  },
  // Bind to UI events using jQuery
  bindEvents: function() {
    $('.option.wallet-select').bind('click', function() {
      var option = $(this).data('option');
      if (option == 'no-wallet') {
        $('.welcome').fadeOut(function(){
          $('.setup-wallet').fadeIn();
        });
      } else {
        $('.welcome').fadeOut(function(){
          $('.setup-mining').fadeIn();
        });
      }
    });

    // TODO: Part of the show more pools hack
    $(document).on('click', '#show_pool_list', function(){
      $(this).hide();
      $('#pool_list_bottom').slideDown();
    });

    $('.button').bind('click', function() {
      var target = $(this).data('target');

      // A bit of validation
      if (target == 'setup-pool') {
        var address = $('#mining_address').val();
        if (address == '') {
          alert("You must enter your address");
          return false;
        }
        // Just make sure they're not using integrated addresses or
        // invalid ones
        if (shared.validateWalletAddress(address) == false)
        {
          alert("Please enter a valid Scala address starting with 'Sv'");
          return false;
        }
      }
      if (target == 'configure_miner') {
        // Send config to Go backend
        // then wait for Go's OK to continue
        var configData = {
          address: $('#mining_address').val(),
          pool: $('#pool_list').find('.selected').data('id')
        };
        astilectron.sendMessage({name: "configure", payload: configData}, function(message){
          document.location = 'index.html';
        });
      }

      $(this).closest('.main-section').fadeOut(function(){
        $('#' + target).closest('.main-section').show();
        $('#' + target).fadeIn(1500);
      });
    });

    $('#exit').bind('click', function(){
      remote.getCurrentWindow().close();
    });


    $(document).on('click', '.pool', function(){
      $('.pool').removeClass('selected');
      $(this).addClass('selected');
      $('#start_mining').show();
    });
  },
  animateIntro: function() {
    // A couple of steps to get you set up
    $('#intro_anim_logo').fadeOut(1000, function() {
      $('#intro_anim_a').fadeIn(2500, function(){
        if (!shared.isMac()) {
          $('#exit').fadeIn(1000);
        }
        $('#intro_anim_a').fadeOut(1000, function(){
          $('#intro_anim_b').fadeIn(2000, function(){
            $('#intro_anim_b').fadeOut(1000, function(){
              $('#intro_anim_c').fadeIn(1500, function(){
                $('#initial-wallet').animateCss('fadeInUp');
              });
            });
          });
        });
      });
    });
  }
};
