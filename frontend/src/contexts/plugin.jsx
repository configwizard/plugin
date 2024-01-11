import React, { createContext, useState, useEffect } from 'react';
import { requestPlugins } from '../api/plugins.js'; // Adjust the path as necessary
export const PluginContext = createContext();

export const PluginProvider = ({ children }) => {
    const [pluginIds, setPluginIds] = useState([]);

    useEffect(() => {
        const fetchPlugins = async () => {
            const ids = await requestPlugins();
            setPluginIds(ids); // Update state with plugin IDs
        };

        fetchPlugins();
    }, []);

    return (
        <PluginContext.Provider value={{ pluginIds }}>
            {children}
        </PluginContext.Provider>
    );
};
