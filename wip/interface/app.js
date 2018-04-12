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
      $('.option.wallet').bind('click', function() {
        var option = $(this).data('option');
        console.log('clicked', option);

        $('.welcome').fadeOut(function(){
          $('.setup-wallet').fadeIn(); 
        })
      });

      // A couple of steps to get you set up

      $('.one').animateCss('fadeInDown', function() {
        // Do somthing after animation
        console.log("Done");
        $('.one').animateCss('fadeOutDown', function(){

          $('.one').remove();
          //$('.two').show();
          $('.two').animateCss('fadeInDown', function(){
            $('.three').animateCss('fadeInDown', function() {
              $('.two').animateCss('fadeOutDown', function(){
                $('.two').remove();
              });
              $('.three').animateCss('fadeOutDown', function(){

                $('.three').remove();
                $('.four').animateCss('fadeInDown', function() {
                  $('.selection').animateCss('fadeInUp');
                });
              });
            });
          });
        });
      });
    }
};
