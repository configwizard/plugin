import React, {useEffect, useContext, useState} from 'react';
import ReactJson from 'react-json-view';
import { ContainerContext } from "../../contexts/container.jsx"
import { requestContainers } from '../../api/containers.js'; // Adjust the path as necessary
import { requestPlugins } from '../../api/plugins.js'; // Adjust the path as necessary
import './styles.css'; // Importing the CSS file
const DeveloperView = () => {
    const { retrievedContainers } = useContext(ContainerContext);
    const [plugins, setPlugins] = useState([]); // State to store plugins

    useEffect(() => {
        const fetchPlugins = async () => {
            const pluginIds = await requestPlugins();
            setPlugins(pluginIds); // Update the state
        };

        fetchPlugins();
    }, []); // Empty dependency array to run only once on component mount

    return (
        <div className="three-column-layout">
            <div className="column">
                <h2>Containers </h2>
                <div className="button-grid">
                    <button onClick={requestContainers}>List Containers</button>
                    <button>Select Container</button>
                    <button>Delete Containerd</button>
                    <button>Button 4</button>
                    <button>Button 5</button>
                    <button>Button 6</button>
                </div>
                <div className="json-viewer">
                    <ReactJson
                        src={retrievedContainers}
                        theme="rjv-default"
                    />
                </div>
            </div>
            <div className="column">
                <h2>Objects</h2>
                <div className="button-grid">
                    <button>Button 1</button>
                    <button>Button 2</button>
                    <button>Button 3</button>
                    <button>Button 4</button>
                    <button>Button 5</button>
                    <button>Button 6</button>
                </div>
                <div className="json-viewer">
                    <ReactJson
                        src={plugins}
                        theme="rjv-default"
                    />
                </div>
            </div>
            <div className="column">
                <h2>Selected Object</h2>
                <div className="button-grid">
                    <button>Button 1</button>
                    <button>Button 2</button>
                    <button>Button 3</button>
                    <button>Button 4</button>
                    <button>Button 5</button>
                    <button>Button 6</button>
                </div>
                <div className="json-viewer">
                    <ReactJson
                        src={{}}
                        theme="rjv-default"
                    />
                </div>
            </div>
        </div>
    );
};

export default DeveloperView;
