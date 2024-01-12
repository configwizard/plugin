// pluginInteraction.js (To be included in the plugin's HTML/JS)
(function(window) {
    function sendMessageToHost(data) {
        const message = {
            pluginID: window.pluginID,
            data: data
        };
        window.parent.postMessage(message, '*'); // Adjust the target origin as needed
    }

    window.pluginHostInteractions = {
        retrieveContainers: function(callback) {
            sendMessageToHost({ action: "retrieveContainers" });
            listenForResponse("containersResponse", callback);
        },
        retrieveContainer: function(id, callback) {
            sendMessageToHost({ action: "retrieveContainer", containerId: id });
            listenForResponse("containerResponse", callback);
        },
        // ... other functions ...
    };

    function listenForResponse(action, callback) {
        const handler = (event) => {
            if (event.data.action === action) {
                callback(event.data.data);
                window.removeEventListener("message", handler);
            }
        };
        window.addEventListener("message", handler);
    }
})(window);
