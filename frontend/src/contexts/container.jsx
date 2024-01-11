import React, { createContext, useState, useEffect } from 'react';
import { EventMessages } from '../api/messages.js'; // Adjust the path as necessary

export const ContainerContext = createContext();

export const ContainerProvider = ({ children }) => {
    const [retrievedContainers, setRetrieveContainers] = useState([]);

    useEffect(() => {
        const addContainer = (newContainer) => {
            let updatedArray = [...retrievedContainers, newContainer];
            updatedArray.sort((a, b) => a.Name.localeCompare(b.Name));
            setRetrieveContainers(updatedArray);
        };

        const removeContainer = (containerToRemove) => {
            const updatedArray = retrievedContainers.filter(container => container.Name !== containerToRemove.Name);
            setRetrieveContainers(updatedArray);
        };

        // Listen for container addition
        window.runtime.EventsOn(EventMessages.ContainerAddUpdate, addContainer);

        // Listen for container removal
        window.runtime.EventsOn(EventMessages.ContainerRemoveUpdate, removeContainer);

        // Cleanup event listeners when context provider is unmounted
        return () => {
            window.runtime.EventsOff(EventMessages.ContainerAddUpdate, addContainer);
            window.runtime.EventsOff(EventMessages.ContainerRemoveUpdate, removeContainer);
        };
    }, [retrievedContainers]);

    return (
        <ContainerContext.Provider value={{ retrievedContainers }}>
            {children}
        </ContainerContext.Provider>
    );
};
