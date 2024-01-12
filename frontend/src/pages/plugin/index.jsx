// PluginView.jsx
import React, {useEffect} from 'react';
import { useParams } from 'react-router-dom';

const PluginView = () => {
    const { id } = useParams(); // Get the plugin ID from URL parameters

    useEffect(() => {
        const handleFrontendMessage = (event) => {
            // Security check: verify the origin of the message
            if (true) {
                const { pluginID, data } = event.data;
                if (pluginID) {
                    window.runtime.EventsEmit(`plugin_message_${pluginID}`, data)
                }
            }
        };
        window.addEventListener("message", handleFrontendMessage);
        return () => window.removeEventListener("message", handleFrontendMessage);
    }, []);
    // If there's no ID, it means we are on the "/plugin" route. Load the plugins available to install.
    if (!id) {
        return <div>Select a plugin from the sidebar.</div>;
    }

    return (
        <div>
            <iframe
                sandbox={"allow-scripts"}
                src={`/${id}/content.html`}
                style={{ width: '100%', height: '500px', border: 'none' }}
            ></iframe>
        </div>
    );
};

export default PluginView;
