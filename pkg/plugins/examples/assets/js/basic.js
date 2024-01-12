// pluginInteraction.js (To be included in the plugin's HTML/JS)

//actions
const RETRIEVE_CONTAINERS = "retrieveContainers";
const RETRIEVE_CONTAINER = "retrieveContainer";

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
            sendMessageToHost({ action: RETRIEVE_CONTAINERS });
            // listenForResponse("containersResponse", callback);
        },
        retrieveContainer: function(id, callback) {
            sendMessageToHost({ action: RETRIEVE_CONTAINER, containerId: id });
            //if we want an immediate but single response.
            // listenForResponse("containerResponse", callback);
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

    window.addEventListener("message", function(event) {
        // Perform security checks here
        console.log("plugin experienced ", event.data)
        if (event.data.action) {
            switch (event.data.action) {
                case "example_action":
                    console.log("containersResponse - example_action", event.pluginID, event.data)
                    // Handle the containers response
                    break;
                // ... handle other actions ...
            }
        }
    });
})(window);
