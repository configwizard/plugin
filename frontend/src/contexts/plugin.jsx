import React, { createContext, useState, useEffect } from 'react';
import { requestPlugins } from '../api/plugins.js'; // Adjust the path as necessary
export const PluginContext = createContext();

export const PluginProvider = ({ children }) => {
    const [plugins, setPlugins] = useState([]); // Changed from pluginIds to plugins

    useEffect(() => {
        const fetchPlugins = async () => {
            const pluginData = await requestPlugins();
            setPlugins(pluginData); // Update state with the full plugin data
        };

        fetchPlugins();
    }, []);

    return (
        <PluginContext.Provider value={{ plugins }}>
            {children}
        </PluginContext.Provider>
    );
};
