let app = {
    init: function() {
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // Wait for the ready signal
        document.addEventListener('astilectron-ready', function() {
          // This will send a message to GO
          astilectron.sendMessage({name: "event.name", payload: "hello"}, function(message) {
            console.log("received " + message.payload)
          });
          app.listen();
        })
    },
    listen: function() {
      astilectron.onMessage(function(message) {
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
    }
};
