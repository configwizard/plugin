// Sidebar.jsx
import React, { useContext } from 'react';
import { PluginContext } from '../../contexts/plugin.jsx';
import { Link } from 'react-router-dom';
import './styles.css';

const Sidebar = () => {
    const { plugins } = useContext(PluginContext);
// Function to handle image load error
    const handleImageError = (e, pluginId) => {
        e.target.onerror = null; // Prevents infinite callback loop
        e.target.src = `/${pluginId}/icon.png`; // Fallback to PNG
    };

    return (
        <div className="sidebar">
            <Link to="/developer">Developer View</Link>
            <Link to="/plugin">Plugins</Link>
            {plugins.map(p => (
                <div key={p.pluginId} className="sidebar-item">
                    <Link to={`/plugin/${p.pluginId}`}>
                        <div className="plugin-info">
                            <img
                                src={`/${p.pluginId}/icon.svg`}
                                alt={`${p.name || p.pluginId} Icon`}
                                className="plugin-icon"
                                onError={(e) => handleImageError(e, p.pluginId)}
                            />
                            <span>{p.name || p.pluginId} {/* Displaying the plugin name or ID */}</span>
                        </div>
                    </Link>
                </div>
            ))}
        </div>
    );
};

export default Sidebar;
