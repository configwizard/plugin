import React, { useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';

const PLUGIN_BACKEND_EVENT = "plugin_backend_event";
const PLUGIN_FRONTEND_EVENT = "plugin_frontend_event";

const PluginView = () => {
    const { id } = useParams();
    const iframeRef = useRef(null);

    useEffect(() => {
        // Backend to Iframe
        const handleBackendMessage = (eventData) => {
            console.log("event data from backend! ", eventData)
            if (iframeRef.current && id === eventData.Id) {
                console.log("forwarding to ", eventData.Id)
                iframeRef.current.contentWindow.postMessage(eventData, '*');
            }
        };

        // Iframe to Backend
        const handleFrontendMessage = (event) => {
            console.log("event origin ", event.origin)
            // if (event.origin === window.location.origin) { // Replace with actual origin check
                const { pluginID, data } = event.data;
                if (pluginID) {
                    console.log("plugged in", event)
                    window.runtime.EventsEmit(`${PLUGIN_FRONTEND_EVENT}_${pluginID}`, data);
                }
            // }
        };

        window.runtime.EventsOn(PLUGIN_BACKEND_EVENT, handleBackendMessage);
        window.addEventListener("message", handleFrontendMessage);

        return () => {
            window.runtime.EventsOff(PLUGIN_BACKEND_EVENT, handleBackendMessage);
            window.removeEventListener("message", handleFrontendMessage);
        };
    }, [id]);

    if (!id) {
        return <div>Select a plugin from the sidebar.</div>;
    }

    return (
        <div>
            <iframe
                ref={iframeRef}
                sandbox="allow-scripts"
                src={`/${id}/content.html`}
                style={{ width: '100%', height: '500px', border: 'none' }}
            ></iframe>
        </div>
    );
};

export default PluginView;
