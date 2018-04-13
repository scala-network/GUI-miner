let app = {
    init: function() {
        /*asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for the ready signal
        document.addEventListener('astilectron-ready', function() {
          // This will send a message to GO
          astilectron.sendMessage({name: "event.name", payload: "hello"}, function(message) {
            console.log("received " + message.payload)
          });
          app.listen();
        })*/
    },
    listen: function() {
      /*astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
        */
    },
    firstrun: function() {
      console.log("firstrun");

      // TODO: Setup handlers somewhere
      $('.option.wallet-select').bind('click', function() {
        var option = $(this).data('option');
        console.log('clicked', option);
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
          window.setTimeout(function(){
            window.location = '/';
          }, 5000);
        }

        console.log('button', target);
        $(this).closest('.main-section').fadeOut(function(){
          $('#' + target).closest('.main-section').show();
          $('#' + target).fadeIn(1500);
        });
      });

      $('.pool').bind('click', function(){
        $('.pool').removeClass('selected');
        $(this).addClass('selected');
        $('#start_mining').show();
      });

      // TODO: Move to own function
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
      // $('#intro_anim_a').fadeIn(500, function(){
      //   $('#intro_anim_a').fadeOut(500, function(){
      //     $('#intro_anim_b').fadeIn(500, function(){
      //       $('#intro_anim_b').fadeOut(500, function(){
      //         $('#intro_anim_c').fadeIn(500, function(){
      //           $('#intro_anim_c').fadeOut(500, function(){
      //             $('#intro_anim_d').fadeIn(500, function(){
      //               $('#initial-wallet').animateCss('fadeInUp');
      //             });
      //           });
      //         });
      //       });
      //     });
      //   });
      // });

    }
};
