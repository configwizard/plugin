// PluginView.jsx
import React from 'react';
import { useParams } from 'react-router-dom';

const PluginView = () => {
    const { id } = useParams(); // Get the plugin ID from URL parameters

    // If there's no ID, it means we are on the "/plugin" route. Load the plugins available to install.
    if (!id) {
        return <div>Select a plugin from the sidebar.</div>;
    }

    return (
        <div>
            <iframe
                sandbox
                src={`/${id}/content.html`}
                style={{ width: '100%', height: '500px', border: 'none' }}
            ></iframe>
        </div>
    );
};

export default PluginView;
