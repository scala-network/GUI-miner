/*
  This handles the initial user setup
 */
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

          astilectron.sendMessage({name: "pool-list", payload: ""}, function(message) {});
          astilectron.sendMessage({name: "startup", payload: ""}, function(message) {
            $('#username').html(message.payload);
          });

          // On firstrun we'll receive the user's username and start firstrun.html
          /*
          var name = app.getParameterByName('name');
          $('#username').html(name);
           */

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

          $('.button').bind('click', function() {
            var target = $(this).data('target');

            // A bit of validation
            if (target == 'setup-pool') {
              var address = $('#mining_address').val();
              if (address == '') {
                alert("You must enter your address");
                return false;
              }
              // Just make sure they're not using integrated addresses
              if (address.substring(0, 2) != 'Se') {
                alert("Please enter a valid Stellite address starting with 'Se'");
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


          $(document).on('click', '.pool', function(){
            $('.pool').removeClass('selected');
            $(this).addClass('selected');
            $('#start_mining').show();
          });

          app.listen();

          // A couple of steps to get you set up
          $('#intro_anim_a').fadeIn(2500, function(){
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


        })
    },
    listen: function() {
      astilectron.onMessage(function(message) {
          console.log("Got message", message.name);
            switch (message.name) {
                case "pool-list":
                  $('#pool_list').html(message.payload);
                  break;
                case "firstrun":
                  $('#username').html(message.payload);
                  break;
                case "about":

                  //return {payload: "payload"};
                  break;
                case "check.out.menu":
                  //asticode.notifier.info(message.payload);
                  break;
            }
        });
    },
    getParameterByName: function(name, url) {
      if (!url) url = window.location.href;
      name = name.replace(/[\[\]]/g, "\\$&");
      var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
          results = regex.exec(url);
      if (!results) return null;
      if (!results[2]) return '';
      return decodeURIComponent(results[2].replace(/\+/g, " "));
    }

};
